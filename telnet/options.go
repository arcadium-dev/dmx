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
	"log/slog"
)

type (
	// ServerOption provides options for configuring the creation of a telnet server.
	ServerOption interface {
		Apply(*Server)
	}
)

// WithServerAddress will configure the server with the listen address.
func WithServerAddress(addr string) ServerOption {
	return newServerOption(func(s *Server) {
		s.addr = addr
	})
}

// WithServerLogger provides a logger to the server.
func WithServerLogger(logger *slog.Logger) ServerOption {
	return newServerOption(func(s *Server) {
		s.logger = logger
	})
}

type (
	serverOption struct {
		f func(*Server)
	}
)

func newServerOption(f func(*Server)) serverOption {
	return serverOption{f: f}
}

func (o serverOption) Apply(s *Server) {
	o.f(s)
}
