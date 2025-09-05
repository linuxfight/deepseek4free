package api

import (
	"fmt"
	"net/http"

	"github.com/linuxfight/deepseek4free/pkg/api/models"
)

// Login is a method to get an API token with credentials
func (c *Client) Login(email, password, deviceId string) (string, error) {
	body := fmt.Sprintf(`{"email":"%s","password":"%s","device_id":"%s","os":"android"}`, email, password, deviceId)

	var data models.AuthResponse
	if err := c.execute(authUrl, body, http.MethodPost, &data); err != nil {
		return "", err
	}

	return data.Data.BizData.User.Token, nil
}

func (c *Client) Logout() error {
	var data models.NullResponse
	if err := c.execute(logoutUrl, "", http.MethodPost, &data); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetProfile() (models.Profile, error) {
	var data models.ProfileResponse
	if err := c.execute(profileUrl, "", http.MethodGet, &data); err != nil {
		return models.Profile{}, err
	}

	return data.Data.BizData, nil
}

func (c *Client) GetQuota() (models.ThinkingQuota, error) {
	var data models.QuotaResponse
	if err := c.execute(quotaUrl, "", http.MethodGet, &data); err != nil {
		return models.ThinkingQuota{}, err
	}
	return data.Data.BizData.Thinking, nil
}

func (c *Client) getPow(endpoint string) (models.PowChallenge, error) {
	body := fmt.Sprintf(`{"target_path":"%s"}`, endpoint)

	var data models.PowResponse
	if err := c.execute(powUrl, body, http.MethodPost, &data); err != nil {
		return models.PowChallenge{}, err
	}

	return data.Data.BizData.Challenge, nil
}
