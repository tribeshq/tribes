// The recommended way to implement a new router is to:
//   - Create a new router with NewRouter()
//   - Add middleware with router.Use()
//   - Add handlers with router.HandleAdvance() or router.HandleInspect()
//   - Create groups with router.Group()
//
// Exemplo de uso do router:
//
//	package main
//
//	import (
//		"encoding/json"
//		"fmt"
//		"github.com/rollmelette/rollmelette"
//		"github.com/tribeshq/tribes/pkg/rollups_router"
//	)
//
//	func main() {
//		router := rollups_router.NewRouter()
//
//		router.Use(rollups_router.LoggingMiddleware)
//
//		router.HandleAdvance("/create", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//			var req struct {
//				Name string `json:"name"`
//			}
//			if err := json.Unmarshal(payload, &req); err != nil {
//				return fmt.Errorf("invalid request: %w", err)
//			}
//			return nil
//		})
//
//		router.HandleInspect("/status", func(env rollmelette.EnvInspector, payload []byte) error {
//			return nil
//		})
//
//		// Grupo de usuários
//		userGroup := router.Group("users")
//		userGroup.HandleAdvance("/create", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//			// ... lógica de criação de usuário
//			return nil
//		})
//	}
package rollups_router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rollmelette/rollmelette"
)

// AdvanceHandler handles advance requests
type AdvanceHandlerFunc func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error

// InspectHandler handles inspect requests
type InspectHandlerFunc func(env rollmelette.EnvInspector, payload []byte) error

// Router handles rollmelette requests
type Router struct {
	advanceHandlers map[string]AdvanceHandlerFunc
	inspectHandlers map[string]InspectHandlerFunc
	middlewares     []Middleware
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		advanceHandlers: make(map[string]AdvanceHandlerFunc),
		inspectHandlers: make(map[string]InspectHandlerFunc),
		middlewares:     make([]Middleware, 0),
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

// Group creates a new route group
func (r *Router) Group(prefix string) *Group {
	return &Group{
		router: r,
		prefix: prefix,
	}
}

// HandleAdvance registers a new advance handler
func (r *Router) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	if path == "" {
		path = "/"
	}
	path = strings.TrimPrefix(path, "/")

	// Apply middlewares in reverse order
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).(AdvanceHandlerFunc)
	}

	r.advanceHandlers[path] = handler
}

// HandleInspect registers a new inspect handler
func (r *Router) HandleInspect(path string, handler InspectHandlerFunc) {
	if path == "" {
		path = "/"
	}
	path = strings.TrimPrefix(path, "/")

	// Apply middlewares in reverse order
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).(InspectHandlerFunc)
	}

	r.inspectHandlers[path] = handler
}

// AdvanceRequest represents the structure of an advance request
type AdvanceRequest struct {
	Path    string          `json:"path"`
	Payload json.RawMessage `json:"payload"`
}

// parseAdvanceRawPayload parses the raw payload into an AdvanceRequest
func parseAdvanceRawPayload(payload []byte) (*AdvanceRequest, error) {
	var req AdvanceRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	if req.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	return &req, nil
}

// Advance handles advance requests
func (r *Router) Advance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	req, err := parseAdvanceRawPayload(payload)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := strings.TrimPrefix(req.Path, "/")
	handler, exists := r.advanceHandlers[path]
	if !exists {
		return fmt.Errorf("no handler found for path: %s", path)
	}

	return handler(env, metadata, deposit, req.Payload)
}

// Inspect handles inspect requests
func (r *Router) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	var req struct {
		Path    string          `json:"path"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if req.Path == "" {
		return fmt.Errorf("path is required")
	}

	path := strings.TrimPrefix(req.Path, "/")
	handler, exists := r.inspectHandlers[path]
	if !exists {
		return fmt.Errorf("no handler found for path: %s", path)
	}

	return handler(env, req.Payload)
}
