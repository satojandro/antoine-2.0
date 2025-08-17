package mcp

import (
	"antoine-cli/internal/models"
	"context"
	"time"
)

type ExaClient struct {
	*BaseMCPClient
}

func NewExaClient() *ExaClient {
	return &ExaClient{
		BaseMCPClient: NewBaseMCPClient(30 * time.Second),
	}
}

func (e *ExaClient) Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	// Implementar llamada específica a Exa
	// Por ahora, mock response
	response := &MCPResponse{
		ID:     generateID(),
		Result: params,
	}
	return response, nil
}

func (e *ExaClient) SearchHackathons(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Hackathon, error) {
	params := map[string]interface{}{
		"query":   query,
		"filters": filters,
		"type":    "hackathons",
	}

	response, err := e.Call(ctx, "search", params)
	if err != nil {
		return nil, err
	}

	var hackathons []*models.Hackathon
	if err := mapToStruct(response.Result, &hackathons); err != nil {
		return nil, err
	}

	return hackathons, nil
}

func (e *ExaClient) SearchProjects(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Project, error) {
	params := map[string]interface{}{
		"query":   query,
		"filters": filters,
		"type":    "projects",
	}

	response, err := e.Call(ctx, "search", params)
	if err != nil {
		return nil, err
	}

	var projects []*models.Project
	if err := mapToStruct(response.Result, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (e *ExaClient) SearchTrends(ctx context.Context, tech []string, timeframe string) (interface{}, error) {
	params := map[string]interface{}{
		"technologies": tech,
		"timeframe":    timeframe,
		"type":         "trends",
	}

	response, err := e.Call(ctx, "search", params)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}
