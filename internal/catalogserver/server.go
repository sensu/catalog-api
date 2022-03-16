package catalogserver

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/transport"
)

func NewCatalogServer(listenAddr string, symlink string) Server {
	t := transport.NewTransport()

	handler := &Handler{
		symlink:   symlink,
		transport: &t,
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: handler,
	}

	return NewServer(server, &t)
}

type Server struct {
	server    *http.Server
	transport *transport.Transport
}

func NewServer(server *http.Server, transport *transport.Transport) Server {
	return Server{
		server:    server,
		transport: transport,
	}
}

func (c *Server) Start(ctx context.Context) {
	// start the transport server
	go c.transport.Start(ctx)

	// start the tcp listener
	listener, err := net.Listen("tcp", c.server.Addr)
	if err != nil {
		panic(err)
	}
	log.Info().Str("address", c.server.Addr).Msg("API server started")

	// serve http requests over the tcp listener
	if err := c.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (c *Server) Stop(ctx context.Context) error {
	return c.server.Shutdown(ctx)
}

func (c *Server) HandleWatchEvent() {
	c.transport.Broadcast([]byte("refresh"))
}
