package agents

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/innomon/agentic/pkg/registry"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type RouterAgentConfig struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	SubAgents   []string `yaml:"sub_agents"`
}

func RegisterRouterAgent() {
	registry.RegisterAgentType("router", func(ctx context.Context, name string, cfg *RouterAgentConfig, models registry.ModelRegistry, tools registry.ToolRegistry, sub []agent.Agent) (agent.Agent, error) {
		return agent.New(agent.Config{
			Name:        name,
			Description: cfg.Description,
			SubAgents:   sub,
			Run: func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
				return func(yield func(*session.Event, error) bool) {
					content := ic.UserContent()
					if content == nil || len(content.Parts) == 0 {
						yield(&session.Event{State: session.StateFailed}, fmt.Errorf("no input provided to router"))
						return
					}

					text := ""
					if txtPart, ok := content.Parts[0].(genai.Text); ok {
						text = strings.ToLower(string(txtPart))
					} else {
						yield(&session.Event{State: session.StateFailed}, fmt.Errorf("non-text input to router not supported yet"))
						return
					}

					var targetAgent agent.Agent

					// Simple Keyword-Based Intent Classification
					if strings.Contains(text, "data") || strings.Contains(text, "schema") || strings.Contains(text, "entity") || strings.Contains(text, "record") {
						for _, a := range sub {
							if strings.Contains(strings.ToLower(a.Name()), "data") || strings.Contains(strings.ToLower(a.Name()), "cms") {
								targetAgent = a
								break
							}
						}
					} else if strings.Contains(text, "ui") || strings.Contains(text, "component") || strings.Contains(text, "dashboard") || strings.Contains(text, "view") {
						for _, a := range sub {
							if strings.Contains(strings.ToLower(a.Name()), "ui") || strings.Contains(strings.ToLower(a.Name()), "a2ui") {
								targetAgent = a
								break
							}
						}
					}

					// Default fallback to first sub-agent
					if targetAgent == nil && len(sub) > 0 {
						targetAgent = sub[0]
					}

					if targetAgent != nil {
						for evt, err := range targetAgent.Run(ic) {
							if !yield(evt, err) {
								return
							}
						}
					} else {
						yield(&session.Event{
							State: session.StateCompleted,
							Content: &genai.Content{
								Role: "model",
								Parts: []genai.Part{
									genai.Text("I'm sorry, I don't have the right sub-agents configured to handle this request."),
								},
							},
						}, nil)
					}
				}
			},
		})
	})
}
