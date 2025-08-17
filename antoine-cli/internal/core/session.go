package core

import (
	"fmt"
	"sync"
	"time"
)

type Session struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Context    map[string]interface{} `json:"context"`
	History    []Command              `json:"history"`
	StartTime  time.Time              `json:"start_time"`
	LastActive time.Time              `json:"last_active"`
	Active     bool                   `json:"active"`
}

type Command struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Input     string                 `json:"input"`
	Output    interface{}            `json:"output"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(userID string) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:         generateSessionID(),
		UserID:     userID,
		Context:    make(map[string]interface{}),
		History:    []Command{},
		StartTime:  time.Now(),
		LastActive: time.Now(),
		Active:     true,
	}

	sm.sessions[session.ID] = session
	return session
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}

func (sm *SessionManager) UpdateSession(sessionID string, updates map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	for key, value := range updates {
		session.Context[key] = value
	}

	session.LastActive = time.Now()
	return nil
}

func (sm *SessionManager) AddCommand(sessionID string, cmd Command) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.History = append(session.History, cmd)
	session.LastActive = time.Now()
	return nil
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
