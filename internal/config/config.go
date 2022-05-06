package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/igor6629/booking/internal/models"
	"html/template"
	"log"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
