// The recommended way to implement a new middleware is to:
//   - Create a function that receives a handler and returns a handler
//   - Use type switch to handle different handler types
//   - Apply the middleware to the router with router.Use()
//
// Example shows the creation and usage of middlewares:
//
// package main
//
// import (
// 	"encoding/json"
// 	"log"
// 	"github.com/rollmelette/rollmelette"
// 	"github.com/tribeshq/tribes/pkg/router"
// )
//
// func main() {
// 	router := router.NewRouter()
//
// 	router.Use(router.ErrorHandlingMiddleware)
// 	router.Use(router.ValidationMiddleware)
// 	router.Use(router.LoggingMiddleware)
//
// 	router.HandleAdvance("/example", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
// 		var req struct {
// 			Name string `json:"name"`
// 		}
// 		if err := json.Unmarshal(payload, &req); err != nil {
// 			return err
// 		}
// 		log.Println("Received payload:", req.Name)
// 		return nil
// 	})
// }

package router

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rollmelette/rollmelette"
)

// Middleware represents a function that wraps a handler
type Middleware func(interface{}) interface{}

// LoggingMiddleware logs request details
func LoggingMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			log.Printf("Advance request - Sender: %s, Payload: %s", metadata.MsgSender.String(), string(payload))
			return h(env, metadata, deposit, payload)
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			log.Printf("Inspect request")
			return h(env, payload)
		})
	default:
		return handler
	}
}

// ValidationMiddleware validates request payload
func ValidationMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			var req Request
			if err := json.Unmarshal(payload, &req); err != nil {
				return fmt.Errorf("invalid request format: %v", err)
			}

			if len(payload) == 0 {
				return fmt.Errorf("empty data")
			}

			return h(env, metadata, deposit, payload)
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			return h(env, payload)
		})
	default:
		return handler
	}
}

// ErrorHandlingMiddleware provides consistent error handling
func ErrorHandlingMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			err := h(env, metadata, deposit, payload)
			if err != nil {
				env.Report([]byte(fmt.Sprintf("Error: %v", err)))
			}
			return err
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			err := h(env, payload)
			if err != nil {
				env.Report([]byte(fmt.Sprintf("Error: %v", err)))
			}
			return err
		})
	default:
		return handler
	}
}
