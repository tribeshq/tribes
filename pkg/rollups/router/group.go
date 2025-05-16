// The Group package provides functionality for grouping routes with common prefixes and middleware.
// Groups are created with the Group method that receives a prefix for its configuration.
// The group information is then stored in the Group struct.
//
// The recommended way to implement a new group is to:
//   - Create a new router with NewRouter()
//   - Create a group with router.Group("prefix")
//   - Add middleware with group.Use(middleware)
//   - Add handlers with group.HandleAdvance() or group.HandleInspect()
//
// Example shows the creation of a user group with authentication:
//
//	package main
//
//	import (
//		"encoding/json"
//		"github.com/rollmelette/rollmelette"
//		"github.com/tribeshq/tribes/pkg/router"
//	)
//
//	func main() {
//		router := router.NewRouter()
//
//		userGroup := router.Group("users")
//		{
//			userGroup.HandleAdvance("", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//				var user struct {
//					Name     string `json:"name"`
//					Email    string `json:"email"`
//					Password string `json:"password"`
//				}
//				if err := json.Unmarshal(payload, &user); err != nil {
//					return err
//				}
//				return nil
//			})
//
//			userGroup.HandleAdvance("login", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//				var login struct {
//					Email    string `json:"email"`
//					Password string `json:"password"`
//				}
//				if err := json.Unmarshal(payload, &login); err != nil {
//					return err
//				}
//				return nil
//			})
//
//			userGroup.HandleInspect("", func(env rollmelette.EnvInspector, payload []byte) error {
//				return nil
//			})
//
//			userGroup.HandleInspect(":id", func(env rollmelette.EnvInspector, payload []byte) error {
//				return nil
//			})
//
//			userGroup.HandleAdvance(":id", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//				var update struct {
//					Name  string `json:"name"`
//					Email string `json:"email"`
//				}
//				if err := json.Unmarshal(payload, &update); err != nil {
//					return err
//				}
//				return nil
//			})
//
//			userGroup.HandleAdvance("delete/:id", func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
//				return nil
//			})
//		}
//	}

package router

import (
	"strings"

	"github.com/rollmelette/rollmelette"
)

type Group struct {
	router     *Router
	prefix     string
	middleware []Middleware
}

func (g *Group) Use(middleware ...Middleware) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) Group(prefix string) *Group {
	if prefix == "" {
		prefix = "/"
	}
	prefix = strings.TrimPrefix(prefix, "/")
	fullPrefix := g.prefix + "/" + prefix
	fullPrefix = strings.TrimPrefix(fullPrefix, "/")

	return &Group{
		router: g.router,
		prefix: fullPrefix,
	}
}

func (g *Group) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	if path == "" {
		path = "/"
	}
	path = strings.TrimPrefix(path, "/")
	fullPath := g.prefix + "/" + path
	fullPath = strings.TrimPrefix(fullPath, "/")

	g.router.HandleAdvance(fullPath, func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
		// Apply group middleware in reverse order
		for i := len(g.middleware) - 1; i >= 0; i-- {
			handler = g.middleware[i](handler).(AdvanceHandlerFunc)
		}
		return handler(env, metadata, deposit, payload)
	})
}

func (g *Group) HandleInspect(path string, handler InspectHandlerFunc) {
	if path == "" {
		path = "/"
	}
	path = strings.TrimPrefix(path, "/")
	fullPath := g.prefix + "/" + path
	fullPath = strings.TrimPrefix(fullPath, "/")

	g.router.HandleInspect(fullPath, func(env rollmelette.EnvInspector, payload []byte) error {
		// Apply group middleware in reverse order
		for i := len(g.middleware) - 1; i >= 0; i-- {
			handler = g.middleware[i](handler).(InspectHandlerFunc)
		}
		return handler(env, payload)
	})
}
