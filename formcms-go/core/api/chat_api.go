package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type ChatApi struct {
	chatService *services.ChatService
}

func NewChatApi(chatService *services.ChatService) *ChatApi {
	return &ChatApi{
		chatService: chatService,
	}
}

func (a *ChatApi) Register(r chi.Router) {
	r.Route("/api/chat", func(r chi.Router) {
		r.Post("/message", a.Message)
	})
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (a *ChatApi) Message(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	// Create a new session
	sess := session.New()
	
	// Get Root Agent
	rootAgent, err := a.chatService.Registry.GetRoot(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("failed to get root agent: %v", err)})
		return
	}

	// Create user prompt event
	userEvent := session.NewUserPromptEvent(genai.Text(req.Message))
	
	// Run the agent
	ic := agent.NewInvocationContext(ctx, rootAgent, nil, nil, sess, nil, nil, userEvent)
	
	var finalResponse string
	for evt, err := range rootAgent.Run(ic) {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("agent error: %v", err)})
			return
		}

		if evt.State == session.StateCompleted && evt.Content != nil {
			for _, part := range evt.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					finalResponse += string(txt)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{Response: finalResponse})
}
