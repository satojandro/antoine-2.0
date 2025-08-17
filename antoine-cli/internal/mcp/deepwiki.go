package mcp

import (
	"context"
	"fmt"
	"time"
)

type DeepWikiClient struct {
	*BaseMCPClient
}

func NewDeepWikiClient() *DeepWikiClient {
	return &DeepWikiClient{
		BaseMCPClient: NewBaseMCPClient(45 * time.Second),
	}
}

func (d *DeepWikiClient) Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	response := &MCPResponse{
		ID:     generateID(),
		Result: params,
	}
	return response, nil
}

func (d *DeepWikiClient) GenerateOverview(ctx context.Context, repoURL string) (string, error) {
	params := map[string]interface{}{
		"repository": repoURL,
		"type":       "overview",
	}

	response, err := d.Call(ctx, "generate_overview", params)
	if err != nil {
		return "", err
	}

	overview, ok := response.Result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response type")
	}

	return overview, nil
}

func (d *DeepWikiClient) GenerateDocumentation(ctx context.Context, repoURL string, sections []string) (map[string]string, error) {
	params := map[string]interface{}{
		"repository": repoURL,
		"sections":   sections,
	}

	response, err := d.Call(ctx, "generate_documentation", params)
	if err != nil {
		return nil, err
	}

	var docs map[string]string
	if err := mapToStruct(response.Result, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}
