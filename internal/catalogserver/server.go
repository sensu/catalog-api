package catalogserver

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/transport"
)

type CatalogServer struct {
	server    *http.Server
	symlink   string
	transport *transport.Transport
}

func NewCatalogServer(listenAddr string, symlink string) CatalogServer {
	t := transport.NewTransport()

	handler := &handler{
		symlink:   symlink,
		transport: &t,
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: handler,
	}

	return CatalogServer{
		server:    server,
		symlink:   symlink,
		transport: &t,
	}
}

func (c *CatalogServer) Start(ctx context.Context) {
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

func (c *CatalogServer) Stop(ctx context.Context) error {
	return c.server.Shutdown(ctx)
}

func (c *CatalogServer) HandleWatchEvent() {
	c.transport.Broadcast([]byte("refresh"))
}
