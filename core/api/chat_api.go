package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/genai"
)

type ChatApi struct {
	chatService *services.ChatService
	authApi     *AuthApi
}

func NewChatApi(chatService *services.ChatService, authApi *AuthApi) *ChatApi {
	return &ChatApi{
		chatService: chatService,
		authApi:     authApi,
	}
}

func (a *ChatApi) Register(r chi.Router) {
	r.Route("/api/chat", func(r chi.Router) {
		r.Use(a.authApi.JWTMiddleware)
		r.Use(a.authApi.RBACMiddleware("read", "chat"))
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
	
	// Get Root Agent
	rootAgent, err := a.chatService.Registry.GetRoot(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("failed to get root agent: %v", err)})
		return
	}

	// Create Runner
	rnr, err := runner.New(runner.Config{
		AppName:        "AiGenCMS",
		Agent:          rootAgent,
		SessionService: a.chatService.SessionService,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("failed to create runner: %v", err)})
		return
	}

	userContent := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: req.Message},
		},
	}
	
	// Use a fixed session ID for now, or generate one
	sessionID := "default-session"
	userID := "default-user"

	var finalResponse string
	for evt, err := range rnr.Run(ctx, userID, sessionID, userContent, agent.RunConfig{}) {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("agent error: %v", err)})
			return
		}

		if evt.Content != nil {
			for _, part := range evt.Content.Parts {
				if part.Text != "" {
					finalResponse += part.Text
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{Response: finalResponse})
}
