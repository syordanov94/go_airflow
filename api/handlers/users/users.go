package users

import (
	"context"
	"go-airflow/airflow"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/life4/genesis/slices"
)

type UsersReq struct {
	Users []User `json:"users"`
}

type User struct {
	Name string `json:"name"`
}

type UsersHandler struct {
	airflowCli AirflowClient
}

type AirflowClient interface {
	PostVariables(ctx context.Context, key, value string) error
}

func NewUsersHandler(airflowCli AirflowClient) *UsersHandler {
	return &UsersHandler{
		airflowCli: airflowCli,
	}
}

func (h UsersHandler) PostUsersV1(c *gin.Context) {
	var req UsersReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a list of all the user names and join them with a comma
	users := slices.Map(req.Users, func(user User) string {
		return user.Name
	})

	val := strings.Join(users, ",")
	err := h.airflowCli.PostVariables(c.Request.Context(), airflow.UserVariableKey, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, req)
	return
}
