package mcp

import (
	"antoine-cli/internal/models"
	"context"
	"fmt"
	"time"
)

type GitHubClient struct {
	*BaseMCPClient
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		BaseMCPClient: NewBaseMCPClient(30 * time.Second),
	}
}

func (g *GitHubClient) Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	// Implementar llamada específica a GitHub MCP
	response := &MCPResponse{
		ID:     generateID(),
		Result: params,
	}
	return response, nil
}

func (g *GitHubClient) AnalyzeRepository(ctx context.Context, repoURL string, options *models.AnalysisOptions) (*models.AnalysisResult, error) {
	params := map[string]interface{}{
		"repository": repoURL,
		"options":    options,
	}

	response, err := g.Call(ctx, "analyze_repository", params)
	if err != nil {
		return nil, err
	}

	var result models.AnalysisResult
	if err := mapToStruct(response.Result, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (g *GitHubClient) GetRepositoryInfo(ctx context.Context, repoURL string) (*models.Repository, error) {
	params := map[string]interface{}{
		"repository": repoURL,
	}

	response, err := g.Call(ctx, "get_repository", params)
	if err != nil {
		return nil, err
	}

	var repo models.Repository
	if err := mapToStruct(response.Result, &repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (g *GitHubClient) ListFiles(ctx context.Context, repoURL string, path string) ([]string, error) {
	params := map[string]interface{}{
		"repository": repoURL,
		"path":       path,
	}

	response, err := g.Call(ctx, "list_files", params)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := mapToStruct(response.Result, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (g *GitHubClient) ReadFile(ctx context.Context, repoURL, filePath string) (string, error) {
	params := map[string]interface{}{
		"repository": repoURL,
		"file_path":  filePath,
	}

	response, err := g.Call(ctx, "read_file", params)
	if err != nil {
		return "", err
	}

	content, ok := response.Result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response type")
	}

	return content, nil
}
