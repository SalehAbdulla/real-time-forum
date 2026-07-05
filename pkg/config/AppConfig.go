package config

import (
	"log/slog"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	Logger        *slog.Logger
	InProduction  bool
	Session       *scs.SessionManager
	LogLevel      string
}
