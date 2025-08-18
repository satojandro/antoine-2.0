package mcp

import (
	"context"
	"fmt"
	"time"
)

// MCPClient defines the interface for MCP (Model Context Protocol) clients
type MCPClient interface {
	Connect(endpoint string) error
	Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error)
	Subscribe(event string, handler EventHandler) error
	Disconnect() error
	IsConnected() bool
	Health() error
}

// MCPResponse represents a response from an MCP server
type MCPResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error response from an MCP server
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// EventHandler is a function type for handling MCP events
type EventHandler func(event *MCPEvent) error

// MCPEvent represents an event from an MCP server
type MCPEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// BaseMCPClient provides a base implementation of MCPClient
type BaseMCPClient struct {
	endpoint   string
	connected  bool
	timeout    time.Duration
	retryCount int
	handlers   map[string][]EventHandler
}

// NewBaseMCPClient creates a new base MCP client
func NewBaseMCPClient(timeout time.Duration) *BaseMCPClient {
	return &BaseMCPClient{
		timeout:    timeout,
		retryCount: 3,
		handlers:   make(map[string][]EventHandler),
	}
}

// Connect establishes a connection to the MCP server
func (c *BaseMCPClient) Connect(endpoint string) error {
	c.endpoint = endpoint
	// TODO: Implementar lógica de conexión específica del protocolo
	c.connected = true
	return nil
}

// Call makes a method call to the MCP server
func (c *BaseMCPClient) Call(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	if !c.connected {
		return nil, fmt.Errorf("client not connected to MCP server")
	}

	// TODO: Implementar lógica real de llamada MCP
	// Por ahora, simulamos una respuesta exitosa
	response := &MCPResponse{
		ID:     generateRequestID(),
		Result: map[string]interface{}{"status": "ok", "method": method},
		Error:  nil,
	}

	// Simular latencia de red
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
		// Continuar
	}

	return response, nil
}

// IsConnected returns whether the client is connected
func (c *BaseMCPClient) IsConnected() bool {
	return c.connected
}

// Disconnect closes the connection to the MCP server
func (c *BaseMCPClient) Disconnect() error {
	c.connected = false
	return nil
}

// Subscribe registers an event handler for a specific event type
func (c *BaseMCPClient) Subscribe(event string, handler EventHandler) error {
	if c.handlers[event] == nil {
		c.handlers[event] = []EventHandler{}
	}
	c.handlers[event] = append(c.handlers[event], handler)
	return nil
}

// Health checks the health of the MCP connection
func (c *BaseMCPClient) Health() error {
	if !c.connected {
		return fmt.Errorf("client not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := c.Call(ctx, "health", nil)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if response.Error != nil {
		return fmt.Errorf("health check error: %s", response.Error.Message)
	}

	return nil
}

// EmitEvent emits an event to all registered handlers
func (c *BaseMCPClient) EmitEvent(eventType string, data interface{}) error {
	handlers, exists := c.handlers[eventType]
	if !exists {
		return nil // No handlers registered for this event type
	}

	event := &MCPEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			// Log error but continue with other handlers
			fmt.Printf("Event handler error for %s: %v\n", eventType, err)
		}
	}

	return nil
}

// SetTimeout sets the timeout for MCP operations
func (c *BaseMCPClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetRetryCount sets the retry count for failed operations
func (c *BaseMCPClient) SetRetryCount(count int) {
	c.retryCount = count
}

// GetEndpoint returns the current endpoint
func (c *BaseMCPClient) GetEndpoint() string {
	return c.endpoint
}

// Helper function to generate request IDs
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// CallWithRetry makes a method call with retry logic
func (c *BaseMCPClient) CallWithRetry(ctx context.Context, method string, params interface{}) (*MCPResponse, error) {
	var lastErr error

	for attempt := 0; attempt < c.retryCount; attempt++ {
		response, err := c.Call(ctx, method, params)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			break
		}

		// Wait before retry (exponential backoff)
		if attempt < c.retryCount-1 {
			waitTime := time.Duration(attempt+1) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
				// Continue to next attempt
			}
		}
	}

	return nil, fmt.Errorf("call failed after %d attempts: %w", c.retryCount, lastErr)
}

// Validate checks if the client configuration is valid
func (c *BaseMCPClient) Validate() error {
	if c.endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	if c.timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.retryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}

	return nil
}
