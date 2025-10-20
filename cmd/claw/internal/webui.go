package internal

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/olivere/vite"
)

var indexTemplate = template.Must(template.New("").Parse( /*html*/
	`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script>window.Otel = {{ .OtelConfig  }};</script>
	<title>{{if .WebUI.Title}}{{.WebUI.Title}}{{else}}Claw{{end}}</title>
	{{- .Vite.Tags -}}
</head>
<body>
	<div id="app"></div>
</body>
`))

type WebUIConfig struct {
	Fragment *vite.Fragment
	DistFS   fs.FS
	Logger   *slog.Logger
}

func CreateWebuiHandler(config WebUIConfig) http.Handler {
	fileServer := http.FileServer(http.FS(config.DistFS))
	mux := http.NewServeMux()
	if cfg.Server.WebUI.DevMode {
		mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	} else {
		mux.Handle("/assets", http.StripPrefix("/assets", fileServer))
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := indexTemplate.Execute(w, map[string]any{
			"Vite":       config.Fragment,
			"WebUI":      cfg.Server.WebUI,
			"OtelConfig": GetOtelConfig().ToJS(),
		}); err != nil {
			config.Logger.Error("failed to execute http template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))
	return mux
}

//go:embed all:webui
var webuiFiles embed.FS
var embeddedWebUIFS, _ = fs.Sub(webuiFiles, "webui")

// CreateViteFragment creates a Vite fragment from the given directory.
// If dir is empty, it uses the embedded webui files.
func CreateViteFragment() (*vite.Fragment, fs.FS, error) {
	var files fs.FS
	if cfg.Server.WebUI.DevMode {
		files = os.DirFS("./webui") // Relative to root working directory. Assuming this golang server runs from project root.
	} else if cfg.Server.WebUI.Path == "" {
		files = embeddedWebUIFS
	} else {
		files = os.DirFS(cfg.Server.WebUI.Path)
	}
	fragment, err := vite.HTMLFragment(vite.Config{
		FS:        files,
		IsDev:     cfg.Server.WebUI.DevMode,
		ViteEntry: "src/main.ts", // This is used in dev mode to load the Vite dev server script.
	})
	if err != nil {
		return nil, files, err
	}
	if cfg.Server.WebUI.DevMode {
		slog.Info("running in development mode; make sure the Vite dev server is running before accessing the web UI")
	}
	return fragment, files, nil
}

type OtelConfig struct {
	OTEL_EXPORTER_OTLP_ENDPOINT        string `json:"OTEL_EXPORTER_OTLP_ENDPOINT,omitempty"`
	OTEL_EXPORTER_OTLP_PROTOCOL        string `json:"OTEL_EXPORTER_OTLP_PROTOCOL,omitempty"`
	OTEL_EXPORTER_OTLP_LOGS_ENDPOINT   string `json:"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT,omitempty"`
	OTEL_EXPORTER_OTLP_LOGS_PROTOCOL   string `json:"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL,omitempty"`
	OTEL_EXPORTER_OTLP_TRACES_ENDPOINT string `json:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT,omitempty"`
	OTEL_EXPORTER_OTLP_TRACES_PROTOCOL string `json:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL,omitempty"`

	DEV_MODE bool `json:"DEV_MODE,omitempty"`
}

func GetOtelConfig() OtelConfig {
	return OtelConfig{
		OTEL_EXPORTER_OTLP_ENDPOINT:        os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTEL_EXPORTER_OTLP_PROTOCOL:        os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"),
		OTEL_EXPORTER_OTLP_LOGS_ENDPOINT:   os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT"),
		OTEL_EXPORTER_OTLP_LOGS_PROTOCOL:   os.Getenv("OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"),
		OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"),
		OTEL_EXPORTER_OTLP_TRACES_PROTOCOL: os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"),
		DEV_MODE:                           cfg.Server.WebUI.DevMode,
	}
}

func (c OtelConfig) ToJS() template.JS {
	b, _ := json.Marshal(c)
	return template.JS(b)
}
