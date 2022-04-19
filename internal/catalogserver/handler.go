package catalogserver

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/transport"
)

func NewHandler(transport *transport.Transport, symlink string) Handler {
	return Handler{
		transport: transport,
		symlink:   symlink,
	}
}

type Handler struct {
	transport *transport.Transport
	symlink   string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ws":
		h.serveWs(w, r)
	default:
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// disable caching of served files
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

		http.FileServer(http.Dir(h.symlink)).ServeHTTP(w, r)
	}
}

func (h Handler) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("remote_addr", r.RemoteAddr).Msg("WebSocket client connected")

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info().Err(err).Msg("WebSocket upgrade failed")
		return
	}

	client := transport.NewClient(h.transport, conn)
	h.transport.Register(&client)
	go client.WritePump()
}
