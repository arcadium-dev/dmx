//  Copyright 2026 arcadium.dev <info@arcadium.dev>
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package telnet // import "arcadium.dev/dmx/telnet"

import (
	"fmt"
	"log/slog"
	"net"

	telnet "github.com/globalcyberalliance/telnet-go"
)

const (
	// DefaultAddr defines the default address of the telnet server when one is
	// not provided via the options. This default to the telnet port on the local
	// host.
	DefaultAddr = ":23"
)

type (
	// Session is an alias for telnet-go.Session. This allows for this package to
	// be used without importing telnet-go as well.
	Session = telnet.Session

	// A Handler responds to telnet activity.
	Handler interface {
		ServeTELNET(*Session)
	}

	// HandlerFunc is an adapter to allow the use of an ordinary function as
	// a telnet handler.
	HandlerFunc func(*telnet.Session)
)

// ServeTELNET calls f(s).
func (f HandlerFunc) ServeTELNET(s *Session) {
	f(s)
}

type (
	// MiddlewareFunc is a function which takes a Handler and returns a Handler.
	// Typically the returned handler is a closure which does something with the
	// Session passed to it, and then calls the handler passed as the parameter
	// to the MiddlewareFunc.
	MiddlewareFunc = func(Handler) Handler

	// Server represent a telnet server.
	Server struct {
		addr   string
		logger *slog.Logger

		middleware []MiddlewareFunc
		handler    Handler

		listener net.Listener
		service  Service
		server   *telnet.Server
	}

	// Service defines the methods required by the server to associate the
	// service with the server.
	Service interface {
		// Register provides a method to install the service handler with the
		// server.
		Register(*Server)

		// Name returns the name of the service.
		Name() string

		// Shutdown allows the service to stop any long running backgroun processes
		// it may have.
		Shutdown()
	}
)

// NewServer create a telnet server with the given server options.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		addr:   DefaultAddr,
		logger: slog.New(slog.DiscardHandler),
		server: &telnet.Server{},
	}
	s.server.Handler = s.handle

	for _, opt := range opts {
		opt.Apply(s)
	}

	s.logger.Info("telnet server created", "address", s.addr)

	return s
}

// Middleware installs the given middleware with the server.
func (s *Server) Middleware(middleware ...MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware...)
}

// HandleFunc installs a handler function for this server.
func (s *Server) HandleFunc(f HandlerFunc) {
	s.handler = f
}

// Register associates the given service with the server.
func (s *Server) Register(service Service) {
	s.service = service
	s.service.Register(s)
	s.logger.Info(fmt.Sprintf("telnet service registered: '%s'", service.Name()))
}

// Serve accepts...
func (s *Server) Serve() error {
	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return err
	}

	s.logger.Info("begin serving telnet", "address", s.addr, "service,", s.service.Name())
	defer s.logger.Info("serving telnet complete", "address", s.addr, "service,", s.service.Name())

	return s.server.Serve(s.listener)
}

// Shutdown stops the telnet server.
func (s *Server) Shutdown() {
	// Stop the service.
	s.service.Shutdown()
	s.logger.Info("telnet service shutdown", "service", s.service.Name())

	// Stop the telnet server.
	if err := s.server.Shutdown(); err != nil {
		s.logger.Error("failed to shutdown http server", "error", err)
	}

	s.logger.Info("telnet server shutdown")
}

func (s *Server) handle(session *telnet.Session) {
	// If the service handler hasn't been register, log an error.
	if s.handler == nil {
		s.logger.Error("telnet service handler not registered")
		return
	}

	// Build a chain of functions, starting with the registered handler function
	// as the last link in the chain. The goal is to have a single function that
	// will call each middleware function, ending with the registered handler
	// function.
	var chain Handler = s.handler

	// Starting at the end of the middleware slice and working backwards, link
	// the functions together.
	for i := len(s.middleware) - 1; i >= 0; i-- {
		chain = s.middleware[i](chain)
	}

	// Call the chain.
	chain.ServeTELNET(session)
}
