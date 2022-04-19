package catalogpreview

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
)

type Handler struct {
	APIHandler http.Handler
	APIURL     string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// mount catalog api as /catalog-api
	if strings.HasPrefix(r.URL.Path, "/catalog-api") {
		http.StripPrefix("/catalog-api", h.APIHandler).ServeHTTP(w, r)
		return
	}

	// forward api & graphql reqests
	apiURL, err := url.Parse(h.APIURL)
	if err != nil {
		zerolog.Ctx(r.Context()).Err(err).Send()
	}
	proxy := httputil.NewSingleHostReverseProxy(apiURL)
	if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/auth") || strings.HasPrefix(r.URL.Path, "/graphql") {
		proxy.ServeHTTP(w, r)
		return
	}

	// serve static web app assets
	fs := assetsHTTPFS.Chroot("/assets")
	if strings.HasPrefix(r.URL.Path, "/static") {
		staticHandler(fs).ServeHTTP(w, r)
		return
	}

	// redirect anything else to /index.html
	rootHandler(fs).ServeHTTP(w, r)
}

func staticHandler(fs http.FileSystem) http.Handler {
	handler := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set cache-control header to allow client to cache file indefinitely.
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching#Freshness
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control#Revalidation_and_reloading
		w.Header().Set("cache-control", "max-age=31536000, immutable")
		handler.ServeHTTP(w, r)
	})
}

func rootHandler(fs http.FileSystem) http.Handler {
	handler := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fallback to index if path didn't match an asset
		if f, _ := fs.Open(r.URL.Path); f == nil {
			r.URL.Path = "/"
		}

		// Sets the proper headers to prevent any sort of caching for the non-static
		// files such as the index and manifest.
		w.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
		w.Header().Set("pragma", "no-cache")
		w.Header().Set("expires", "0")
		handler.ServeHTTP(w, r)
	})
}
