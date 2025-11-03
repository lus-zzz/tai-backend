package flowy

import (
	"flowy-sdk/pkg/client"
	"flowy-sdk/pkg/config"
	"flowy-sdk/services/agent"
	"flowy-sdk/services/knowledge"
	"flowy-sdk/services/model"
)

// SDK Flowy SDK main entry point
type SDK struct {
	config *config.Config
	client client.HTTPClient

	// Services
	Model     model.Service
	Knowledge knowledge.Service
	Agent     agent.Service
}

// New creates a new Flowy SDK instance
func New(cfg *config.Config) *SDK {
	if cfg == nil {
		cfg = config.DefaultConfig().LoadFromEnv()
	} else {
		cfg.LoadFromEnv()
	}

	// Create HTTP client
	httpClient := client.New(cfg)

	// Create service instances
	modelService := model.NewService(httpClient)
	knowledgeService := knowledge.NewService(httpClient)
	agentService := agent.NewService(httpClient)

	return &SDK{
		config:    cfg,
		client:    httpClient,
		Model:     modelService,
		Knowledge: knowledgeService,
		Agent:     agentService,
	}
}
