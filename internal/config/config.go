package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

type TemplateCache map[string]*template.Template

type AppConfig struct {
	UseCache      bool
	TemplateCache TemplateCache
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}
