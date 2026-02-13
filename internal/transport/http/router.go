package http

import (
	"html/template"
	"log/slog"
	"net/http"
)

func NewRouter(handler *PackCalculatorHandler, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/calculate", handler.CalculatePacks)
	mux.HandleFunc("GET /api/pack-sizes", handler.GetPackSizes)
	mux.HandleFunc("PUT /api/pack-sizes", handler.UpdatePackSizes)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			slog.Error("failed to render template", "error", err)
		}
	})

	return Chain(mux, Recovery, Logging, CORS)
}
