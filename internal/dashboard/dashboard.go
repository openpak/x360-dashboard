// Package dashboard serves OpenPak's dashboard Home manifest. Resource URIs
// refer to files already supplied by the console; no firmware assets are hosted.
package dashboard

import (
	"embed"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
)

//go:embed dashhome.xml
var content embed.FS

var manifestPath = regexp.MustCompile(`^/epix/[A-Za-z]{2,3}(?:-[A-Za-z0-9]{2,8})*/dashhome\.xml$`)

// Handler serves public dashboard content and logs missing resource requests.
type Handler struct{ Log *slog.Logger }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	defer func() {
		if h.Log != nil {
			h.Log.Info("request", "method", r.Method, "path", r.URL.Path, "status", status, "user_agent", r.UserAgent())
		}
	}()
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		status = http.StatusMethodNotAllowed
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", status)
		return
	}
	if r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !manifestPath.MatchString(r.URL.Path) {
		status = http.StatusNotFound
		http.NotFound(w, r)
		return
	}
	body, err := content.ReadFile("dashhome.xml")
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, "content unavailable", status)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Cache-Control", "public, max-age=60")
	if r.Method == http.MethodGet {
		w.Write(body)
	}
}
