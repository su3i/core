package airbyte

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/darksuei/suei-intelligence/internal/config"
	domain "github.com/darksuei/suei-intelligence/internal/domain/etl"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache"
)

type AirbyteContext struct {
	cfg    *config.AirbyteConfig
	ctx    context.Context
}

func retrieveAccessToken(cfg *config.AirbyteConfig) (string, error) {
	cacheKey := "airbyte__access__token"

	accessToken, err := cache.GetCache().Get(cacheKey)

	if err == nil {
		return accessToken, nil
	}

	payload := map[string]string{
		"client_id":     cfg.AirbyteClientId,
		"client_secret": cfg.AirbyteClientSecret,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	resp, err := http.Post(
		cfg.AirbyteEndpoint+"/api/v1/applications/token",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to request access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("received empty access token from Airbyte")
	}

	cache.GetCache().Set(cacheKey, result.AccessToken, 10 * time.Minute)

	log.Printf("Successfully retrieved new access token - %s", result.AccessToken)

	return result.AccessToken, nil
}

func Initialize(c *config.AirbyteConfig) domain.ETL {
	return &AirbyteContext {
		cfg: c,
		ctx: context.Background(),
	}
}

func (c *AirbyteContext) CreateSourceConnection(name string, configuration interface{}) error {
	token, err := retrieveAccessToken(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve access token: %w", err)
	}

	payload := map[string]interface{}{
		"name":     name,
		"workspaceId": c.cfg.AirbyteWorkspaceId,
		"configuration": configuration,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.cfg.AirbyteEndpoint+"/api/public/v1/sources", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create source connection: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create source connection (status %d): %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Successfully created new source - %s", respBody)

	return nil
}
