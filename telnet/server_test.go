package telnet_test

import (
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"arcadium.dev/dmx/assert"
	"arcadium.dev/dmx/telnet"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	var (
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	)

	tests := []struct {
		name   string
		opts   []telnet.ServerOption
		verify func(*testing.T, *telnet.Server)
	}{
		{
			name: "test defaults",
			verify: func(t *testing.T, s *telnet.Server) {
				assert.NotNil(t, s)
				assert.Equal(t, s.Addr(), telnet.DefaultAddr)
				assert.Compare(t, s.Logger().Handler(), slog.DiscardHandler)
			},
		},
		{
			name: "test WithServerAddress option",
			opts: []telnet.ServerOption{
				telnet.WithServerAddress(":2323"),
			},
			verify: func(t *testing.T, s *telnet.Server) {
				assert.NotNil(t, s)
				assert.Equal(t, s.Addr(), ":2323")
				assert.Compare(t, s.Logger().Handler(), slog.DiscardHandler)
			},
		},
		{
			name: "test WithServerLogger option",
			opts: []telnet.ServerOption{
				telnet.WithServerLogger(logger),
			},
			verify: func(t *testing.T, s *telnet.Server) {
				assert.NotNil(t, s)
				assert.Equal(t, s.Addr(), telnet.DefaultAddr)
				assert.Equal(t, s.Logger(), logger)
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s := telnet.NewServer(test.opts...)
			test.verify(t, s)
		})
	}
}

func TestRegister(t *testing.T) {
	s := telnet.NewServer()
	assert.NotNil(t, s)

	m := &mockService{}
	s.Register(m)
	assert.True(t, m.registerCalled)

	mw := &mockMiddleware{}
	s.Middleware(mw.mock)

	session := &telnet.Session{}
	s.Server().Handler.ServeTELNET(session)
	assert.True(t, m.handlerCalled)
	assert.True(t, mw.called)
	assert.True(t, mw.whenHandled.Before(m.whenHandled))
}

func TestServe(t *testing.T) {

	tests := []struct {
		name    string
		opts    []telnet.ServerOption
		service *mockService
		verify  func(*testing.T, *mockService, error)
	}{
		{
			name: "listen failure",
			opts: []telnet.ServerOption{
				telnet.WithServerAddress(":-42"),
			},
			service: &mockService{},
			verify: func(t *testing.T, service *mockService, err error) {
				assert.Error(t, err, "listen tcp: address -42: invalid port")
				assert.True(t, service.shutdownCalled)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := telnet.NewServer(test.opts...)
			s.Register(test.service)

			result := make(chan error, 1)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() { wg.Done(); result <- s.Serve() }()
			wg.Wait()
			s.Shutdown()

			err := <-result
			test.verify(t, test.service, err)
		})
	}
}

/*
func TestShutdown(t *testing.T) {
}
*/

type (
	mockService struct {
		registerCalled, handlerCalled, shutdownCalled bool
		whenHandled                                   time.Time
	}
)

var _ telnet.Service = (*mockService)(nil)

func (m *mockService) Register(s *telnet.Server) {
	m.registerCalled = true
	s.HandleFunc(m.handle)
}

func (m *mockService) Name() string {
	return "mockService"
}

func (m *mockService) Shutdown() {
	m.shutdownCalled = true
}

func (m *mockService) handle(s *telnet.Session) {
	m.handlerCalled = true
	m.whenHandled = time.Now()
}

type (
	mockMiddleware struct {
		called      bool
		whenHandled time.Time
	}
)

func (m *mockMiddleware) mock(next telnet.Handler) telnet.Handler {
	return telnet.HandlerFunc(func(session *telnet.Session) {
		m.called = true
		m.whenHandled = time.Now()
		next.ServeTELNET(session)
	})
}
