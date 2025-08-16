package mcp

import (
	"context"
	"time"
)

type E2BClient struct {
	*BaseMCPClient
}

type CodeExecutionRequest struct {
	Code        string            `json:"code"`
	Language    string            `json:"language"`
	Environment string            `json:"environment"`
	Timeout     int               `json:"timeout"`
	Args        map[string]string `json:"args"`
}

type CodeExecutionResult struct {
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
	Duration int    `json:"duration_ms"`
}

func NewE2BClient() *E2BClient {
	return &E2BClient{
		BaseMCPClient: NewBaseMCPClient(60 * time.Second),
	}
}

func (e *E2BClient) Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	response := &MCPResponse{
		ID:     generateID(),
		Result: params,
	}
	return response, nil
}

func (e *E2BClient) ExecuteCode(ctx context.Context, req *CodeExecutionRequest) (*CodeExecutionResult, error) {
	response, err := e.Call(ctx, "execute_code", req)
	if err != nil {
		return nil, err
	}

	var result CodeExecutionResult
	if err := mapToStruct(response.Result, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *E2BClient) RunAnalysis(ctx context.Context, script string, data interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"script": script,
		"data":   data,
	}

	response, err := e.Call(ctx, "run_analysis", params)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}
