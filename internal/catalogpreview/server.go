package catalogpreview

import (
	"net/http"

	"github.com/sensu/catalog-api/internal/catalogserver"
	"github.com/sensu/catalog-api/internal/transport"
)

func NewPreviewServer(listenAddr, symlink, apiURL string) catalogserver.Server {
	t := transport.NewTransport()

	server := &http.Server{
		Addr: listenAddr,
		Handler: Handler{
			APIHandler: catalogserver.NewHandler(&t, symlink),
			APIURL:     apiURL,
		},
	}
	return catalogserver.NewServer(server, &t)
}
