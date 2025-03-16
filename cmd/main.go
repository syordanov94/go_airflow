package main

import (
	"go-airflow/airflow"
	"go-airflow/api/handlers/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	handlers := initHandlers()
	for _, handler := range handlers {
		switch handler.OperationType {
		case http.MethodPost:
			router.POST(handler.Path, handler.Func)
		case http.MethodGet:
			router.GET(handler.Path, handler.Func)
		case http.MethodDelete:
			router.DELETE(handler.Path, handler.Func)
		case http.MethodPatch:
			router.PATCH(handler.Path, handler.Func)
		case http.MethodPut:
			router.PUT(handler.Path, handler.Func)
		case http.MethodOptions:
			router.OPTIONS(handler.Path, handler.Func)
		case http.MethodHead:
			router.HEAD(handler.Path, handler.Func)
		default:
			panic("Invalid operation type")
		}
	}

	router.Run("localhost:8080")
}

type HandlerCfg struct {
	Func          gin.HandlerFunc
	OperationType string
	Path          string
}

func initHandlers() []HandlerCfg {

	// TODO: Don't init the client with hardcoded values
	aifrlowCli := airflow.NewAirflowCli("localhost:8080", "http", "airflow", "airflow")

	ret := []HandlerCfg{}

	usersHandler := users.NewUsersHandler(aifrlowCli)
	ret = append(ret, HandlerCfg{
		Func:          usersHandler.PostUsersV1,
		OperationType: http.MethodPost,
		Path:          "/users",
	})

	return ret
}

// func main() {
// 	conf := airflow.NewConfiguration()
// 	conf.Host = "localhost:8080"
// 	conf.Scheme = "http"
// 	cli := airflow.NewAPIClient(conf)

// 	cred := airflow.BasicAuth{
// 		UserName: "test",
// 		Password: "password",
// 	}
// 	ctx := context.WithValue(context.Background(), airflow.ContextBasicAuth, cred)

// 	dag, resp, err := cli.DAGApi.GetDag(ctx, "example_dag").Execute()
// 	if resp != nil {
// 		body := []byte{}
// 		resp.Body.Read(body)
// 		fmt.Println(resp.Status + string(body))
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Print(dag.DagId)

// 	run := airflow.NewDAGRunWithDefaults()
// 	_, resp, err = cli.DAGRunApi.PostDagRun(ctx, "example_dag").DAGRun(*run).Execute()
// 	if resp != nil {
// 		body := []byte{}
// 		resp.Body.Read(body)
// 		fmt.Println(resp.Status + string(body))
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// variable
// 	var (
// 		key string
// 		val string
// 	)

// 	key = "users"
// 	val = strings.Join([]string{"user1", "user2"}, ",")
// 	_, resp, err = cli.VariableApi.PostVariables(ctx).Variable(airflow.Variable{
// 		Key:   &key,
// 		Value: &val,
// 	}).Execute()
// 	if resp != nil {
// 		body := []byte{}
// 		resp.Body.Read(body)
// 		fmt.Println(resp.Status + string(body))
// 	}
// }
