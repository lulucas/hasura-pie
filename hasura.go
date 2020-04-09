package pie

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"time"
)

type Action struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input            json.RawMessage `json:"input"`
	SessionVariables struct {
		XHasuraUserId string `json:"x-hasura-user-id"`
		XHasuraRole   string `json:"x-hasura-role"`
	} `json:"session_variables"`
}

type Event struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Trigger   struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
	Event struct {
		SessionVariables struct {
			XHasuraRole         string `json:"x-hasura-role"`
			XHasuraAllowedRoles string `json:"x-hasura-allowed-roles"`
			XHasuraUserID       string `json:"x-hasura-user-id"`
		} `json:"session_variables"`
		Op   string `json:"op"`
		Data struct {
			Old json.RawMessage `json:"old"`
			New json.RawMessage `json:"new"`
		} `json:"data"`
	} `json:"event"`
}

func hasuraErrorResponse(err error) interface{} {
	return echo.Map{
		"message": err.Error(),
		"code":    "400",
	}
}
