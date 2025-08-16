package mcp

import (
	"context"
	"fmt"
	"time"
)

type MCPClient interface {
	Connect(endpoint string) error
	Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error)
	Subscribe(event string, handler EventHandler) error
	Disconnect() error
	IsConnected() bool
	Health() error
}

type MCPResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EventHandler func(event *MCPEvent) error

type MCPEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type BaseMCPClient struct {
	endpoint   string
	connected  bool
	timeout    time.Duration
	retryCount int
	handlers   map[string][]EventHandler
}

func NewBaseMCPClient(timeout time.Duration) *BaseMCPClient {
	return &BaseMCPClient{
		timeout:    timeout,
		retryCount: 3,
		handlers:   make(map[string][]EventHandler),
	}
}

func (c *BaseMCPClient) Connect(endpoint string) error {
	c.endpoint = endpoint
	// Implementar lógica de conexión específica del protocolo
	c.connected = true
	return nil
}

func (c *BaseMCPClient) IsConnected() bool {
	return c.connected
}

func (c *BaseMCPClient) Disconnect() error {
	c.connected = false
	return nil
}

func (c *BaseMCPClient) Subscribe(event string, handler EventHandler) error {
	if c.handlers[event] == nil {
		c.handlers[event] = []EventHandler{}
	}
	c.handlers[event] = append(c.handlers[event], handler)
	return nil
}

func (c *BaseMCPClient) Health() error {
	if !c.connected {
		return fmt.Errorf("client not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Call(ctx, "health", nil)
	return err
}
