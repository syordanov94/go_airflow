package users

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"go-airflow/airflow"

	"github.com/gin-gonic/gin"
	"github.com/life4/genesis/slices"
)

type AirflowClient interface {
	PostVariables(ctx context.Context, key, value string) error
	GetVariable(ctx context.Context, key string) (string, error)
	TriggerDag(ctx context.Context, dagID string) error
}

type UsersHandler struct {
	airflowCli AirflowClient
}

func NewUsersHandler(airflowCli AirflowClient) *UsersHandler {
	return &UsersHandler{
		airflowCli: airflowCli,
	}
}

type UsersReq struct {
	Users []User `json:"users"`
}

type User struct {
	Name string `json:"name"`
}

func (h UsersHandler) PostUsersV1(c *gin.Context) {
	var req UsersReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a list of all the usernames and join them with a comma
	users := slices.Map(req.Users, func(user User) string {
		return user.Name
	})
	val := strings.Join(users, ",")

	// post the users to the airflow API
	err := h.airflowCli.PostVariables(c.Request.Context(), airflow.UserVariableKey, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, req)
}

type UsersScoreResp struct {
	Scores []UserScore `json:"scores"`
}

type UserScore struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func (h UsersHandler) GetUserScoreV1(c *gin.Context) {
	// get the requested username from the request
	userName, found := c.Params.Get("userName")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user name not found"})
		return
	}

	// get the user scores from airflow
	scores, err := h.airflowCli.GetVariable(c.Request.Context(), airflow.PlayerScoresVariableKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// split the scores by comma
	userScores := strings.Split(scores, "|")

	for _, userScoreStr := range userScores {
		// split the user score by colon and check if the username matches the requested one
		userScore := strings.Split(userScoreStr, "#")
		if userScore[0] == userName {
			score, err := strconv.ParseFloat(userScore[1], 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusOK, UsersScoreResp{
				Scores: []UserScore{
					{
						Name:  userScore[0],
						Score: score,
					},
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func (h *UsersHandler) TriggerPlayerScoreUpdateManuallyV1(c *gin.Context) {
	// trigger the DAG to update the player scores
	err := h.airflowCli.TriggerDag(c.Request.Context(), airflow.PlayerScoresCalcDagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player scores update triggered"})
}
