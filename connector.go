// Package ccconnect provides a unified interface for connecting to various
// chat and AI platform APIs. It is a fork of chenhg5/cc-connect with
// additional platform support and improved error handling.
package ccconnect

import (
	"context"
	"errors"
	"io"
)

// Role represents the role of a message sender in a conversation.
type Role string

const (
	// RoleUser represents a message from the user.
	RoleUser Role = "user"
	// RoleAssistant represents a message from the AI assistant.
	RoleAssistant Role = "assistant"
	// RoleSystem represents a system-level message or instruction.
	RoleSystem Role = "system"
)

// Message represents a single message in a conversation.
type Message struct {
	// Role is the sender's role in the conversation.
	Role Role `json:"role"`
	// Content is the text content of the message.
	Content string `json:"content"`
}

// Request encapsulates the parameters for a chat completion request.
type Request struct {
	// Messages is the conversation history to send to the model.
	Messages []Message `json:"messages"`
	// Model is the identifier of the model to use (platform-specific).
	Model string `json:"model,omitempty"`
	// MaxTokens limits the number of tokens in the response.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Temperature controls randomness in generation (0.0–2.0).
	Temperature float64 `json:"temperature,omitempty"`
	// Stream indicates whether to stream the response token by token.
	Stream bool `json:"stream,omitempty"`
}

// Response holds the result of a chat completion request.
type Response struct {
	// Content is the generated text from the model.
	Content string `json:"content"`
	// Model is the model that generated the response.
	Model string `json:"model,omitempty"`
	// Usage contains token usage statistics for the request.
	Usage *Usage `json:"usage,omitempty"`
}

// Usage tracks token consumption for a single request.
type Usage struct {
	// PromptTokens is the number of tokens in the input.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens is the number of tokens in the output.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens is the sum of prompt and completion tokens.
	TotalTokens int `json:"total_tokens"`
}

// StreamChunk represents a single chunk received during streaming.
type StreamChunk struct {
	// Delta is the incremental text content for this chunk.
	Delta string
	// Done indicates whether this is the final chunk in the stream.
	Done bool
}

// Connector defines the interface that all platform connectors must implement.
type Connector interface {
	// Chat sends a request and returns a complete response.
	Chat(ctx context.Context, req *Request) (*Response, error)
	// ChatStream sends a request and streams the response via the returned reader.
	// The caller is responsible for closing the returned ReadCloser.
	ChatStream(ctx context.Context, req *Request) (io.ReadCloser, error)
	// Name returns the human-readable name of the platform connector.
	Name() string
}

// ErrEmptyMessages is returned when a request contains no messages.
var ErrEmptyMessages = errors.New("ccconnect: request must contain at least one message")

// ErrInvalidRole is returned when a message contains an unrecognised role.
var ErrInvalidRole = errors.New("ccconnect: invalid message role")

// ErrStreamingNotSupported is returned when a connector does not support streaming.
var ErrStreamingNotSupported = errors.New("ccconnect: streaming is not supported by this connector")

// ValidateRequest performs basic validation on a Request before it is sent.
func ValidateRequest(req *Request) error {
	if req == nil {
		return errors.New("ccconnect: request must not be nil")
	}
	if len(req.Messages) == 0 {
		return ErrEmptyMessages
	}
	for _, m := range req.Messages {
		switch m.Role {
		case RoleUser, RoleAssistant, RoleSystem:
			// valid
		default:
			return ErrInvalidRole
		}
	}
	return nil
}
