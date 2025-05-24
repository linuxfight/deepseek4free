package api

import (
	"hashTest/solver"
	"net/http"
)

const chatCreateBody = `{"agent":"chat"}`
const authUrl = "https://chat.deepseek.com/api/v0/users/login"
const chatCreateUrl = "https://chat.deepseek.com/api/v0/chat_session/create"
const chatEditUrl = "https://chat.deepseek.com/api/v0/chat_session/update_title"
const chatListUrl = "https://chat.deepseek.com/api/v0/chat_session/fetch_page"
const completionUrl = "https://chat.deepseek.com/api/v0/chat/completion"
const historyBaseUrl = "https://chat.deepseek.com/api/v0/chat/history_messages?chat_session_id="
const logoutUrl = "https://chat.deepseek.com/api/v0/users/logout"
const powUrl = "https://chat.deepseek.com/api/v0/chat/create_pow_challenge"
const profileUrl = "https://chat.deepseek.com/api/v0/users/current"
const quotaUrl = "https://chat.deepseek.com/api/v0/users/feature_quota"

type Client struct {
	ApiKey     string
	rangersId  string
	powSolver  *solver.Solver
	httpClient *http.Client
}

// New is a method to create a new DeepSeek mobile API client. If ApiKey is "", you would need to log in.
func New(powSolver *solver.Solver, rangersId string, apiKey string) *Client {
	if rangersId == "" {
		panic("no rangersId")
	}

	return &Client{
		ApiKey:     apiKey,
		rangersId:  rangersId,
		powSolver:  powSolver,
		httpClient: &http.Client{},
	}
}
