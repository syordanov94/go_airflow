package airflow

import (
	"context"
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
)

const (
	UserVariableKey         = "users"
	PlayerScoresVariableKey = "player_scores"
	PlayerScoresCalcDagID   = "player_score_calculation_dag"
)

type AirflowClient struct {
	Cli   *airflow.APIClient
	Creds airflow.BasicAuth
}

func NewAirflowCli(host, scheme, username, password string) *AirflowClient {
	conf := airflow.NewConfiguration()
	conf.Host = host
	conf.Scheme = scheme
	cli := airflow.NewAPIClient(conf)

	creds := airflow.BasicAuth{
		UserName: username,
		Password: password,
	}

	return &AirflowClient{
		Cli:   cli,
		Creds: creds,
	}
}

func (c AirflowClient) PostVariables(ctx context.Context, key, value string) error {
	ctx = context.WithValue(ctx, airflow.ContextBasicAuth, c.Creds)
	_, resp, err := c.Cli.VariableApi.PostVariables(ctx).Variable(airflow.Variable{
		Key:   &key,
		Value: &value,
	}).Execute()
	if resp != nil {
		if resp.StatusCode != 200 {
			return fmt.Errorf("error: %s", resp.Status)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (c AirflowClient) GetVariable(ctx context.Context, key string) (string, error) {
	ctx = context.WithValue(ctx, airflow.ContextBasicAuth, c.Creds)
	variable, resp, err := c.Cli.VariableApi.GetVariable(ctx, key).Execute()
	if resp != nil {
		if resp.StatusCode != 200 {
			return "", fmt.Errorf("error: %s", resp.Status)
		}
	}

	if err != nil {
		return "", err
	}

	if variable.Value == nil {
		return "", fmt.Errorf("variable %s not found", key)
	}

	return *variable.Value, nil
}

func (c AirflowClient) TriggerDag(ctx context.Context, dagID string) error {
	ctx = context.WithValue(ctx, airflow.ContextBasicAuth, c.Creds)
	_, resp, err := c.Cli.DAGRunApi.PostDagRun(ctx, dagID).DAGRun(*airflow.NewDAGRunWithDefaults()).Execute()
	if resp != nil {
		if resp.StatusCode != 200 {
			return fmt.Errorf("error: %s", resp.Status)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
