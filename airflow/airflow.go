package airflow

import (
	"context"
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
)

const (
	UserVariableKey = "users"
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
