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

	ret = append(ret, HandlerCfg{
		Func:          usersHandler.GetUserScoreV1,
		OperationType: http.MethodGet,
		Path:          "/users/:userName/score",
	})

	ret = append(ret, HandlerCfg{
		Func:          usersHandler.TriggerPlayerScoreUpdateManuallyV1,
		OperationType: http.MethodPost,
		Path:          "/users/score",
	})

	return ret
}
