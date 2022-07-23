package server

import (
	"context"
)

// Option is a function that can be used to set options on the Server.
type Option func(*Server)

// WithContext sets the context on the Server.
func WithContext(ctx context.Context) Option {
	return func(s *Server) {
		s.ctx = ctx
	}
}

// WithName sets the name on the Server.
func WithName(name string) Option {
	return func(s *Server) {
		s.name = name
	}
}

// WithVersion sets the version on the Server.
func WithVersion(version string) Option {
	return func(s *Server) {
		s.version = version
	}
}

// WithPort sets the port on the Server.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}
