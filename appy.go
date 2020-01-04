package appy

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/appist/appy/cmd"
	ah "github.com/appist/appy/http"
	"github.com/appist/appy/mailer"
	"github.com/appist/appy/support"
)

type (
	// App is the framework core that drives the application.
	App struct {
		assets  *support.Assets
		command *cmd.Command
		config  *support.Config
		i18n    *support.I18n
		logger  *support.Logger
		mailer  *mailer.Mailer
		server  *ah.Server
	}

	// Context contains the request information and is meant to be passed through the entire HTTP request.
	Context = ah.Context

	// Email defines the email headers/body/attachments.
	Email = mailer.Email

	// Mailer provides the capability to parse/render email template and send it out via SMTP protocol.
	Mailer = mailer.Mailer

	// H is a shortcut for map[string]interface{}.
	H = support.H

	// Server serves the HTTP requests.
	Server = ah.Server
)

func init() {
	if os.Getenv("APPY_ENV") == "" {
		os.Setenv("APPY_ENV", "development")
	}
}

// NewApp initializes an app instance.
func NewApp(static http.FileSystem, viewFuncs map[string]interface{}) *App {
	command := cmd.NewCommand()
	logger := support.NewLogger()
	assets := support.NewAssets(nil, "", static)
	config := support.NewConfig(assets, logger)
	i18n := support.NewI18n(assets, config, logger)
	server := ah.NewServer(assets, config, logger)
	mailer := mailer.NewMailer(assets, config, i18n, viewFuncs)

	// Setup the default middleware.
	server.Use(ah.CSRF(config, logger))
	server.Use(ah.RequestID())
	server.Use(ah.RequestLogger(config, logger))
	server.Use(ah.RealIP())
	server.Use(ah.ResponseHeaderFilter())
	server.Use(ah.HealthCheck(config.HTTPHealthCheckURL))
	server.Use(ah.Prerender(config, logger))
	server.Use(ah.Gzip(config))
	server.Use(ah.Secure(config))
	server.Use(ah.I18n(i18n))
	server.Use(ah.ViewEngine(assets, viewFuncs))
	server.Use(ah.Mailer(i18n, mailer, server))
	server.Use(ah.SessionMngr(config))
	server.Use(ah.Recovery(logger))

	// Setup the default commands.
	if support.IsDebugBuild() {
		command.AddCommand(cmd.NewBuildCommand(assets, logger, server))
		command.AddCommand(cmd.NewStartCommand(logger, server))
	}

	command.AddCommand(cmd.NewConfigDecryptCommand(config, logger))
	command.AddCommand(cmd.NewConfigEncryptCommand(config, logger))
	command.AddCommand(cmd.NewDcDownCommand(assets, logger))
	command.AddCommand(cmd.NewDcRestartCommand(assets, logger))
	command.AddCommand(cmd.NewDcUpCommand(assets, logger))
	command.AddCommand(cmd.NewMiddlewareCommand(config, logger, server))
	command.AddCommand(cmd.NewRoutesCommand(config, logger, server))
	command.AddCommand(cmd.NewSecretCommand(logger))
	command.AddCommand(cmd.NewServeCommand(logger, server))
	command.AddCommand(cmd.NewSSLSetupCommand(logger, server))
	command.AddCommand(cmd.NewSSLTeardownCommand(logger, server))

	return &App{
		assets:  assets,
		command: command,
		config:  config,
		i18n:    i18n,
		logger:  logger,
		mailer:  mailer,
		server:  server,
	}
}

// Cmd returns the cmd instance.
func (a App) Cmd() *cmd.Command {
	return a.command
}

// Config returns the config instance.
func (a App) Config() *support.Config {
	return a.config
}

// I18n returns the I18n instance.
func (a App) I18n() *support.I18n {
	return a.i18n
}

// Logger returns the logger instance.
func (a App) Logger() *support.Logger {
	return a.logger
}

// Mailer returns the mailer instance.
func (a App) Mailer() *Mailer {
	return a.mailer
}

// Server returns the server instance.
func (a App) Server() *Server {
	return a.server
}

// Run starts running the app instance.
func (a App) Run() error {
	a.server.ServeSPA("/", a.assets.Static())
	a.server.ServeNoRoute()

	return a.command.Execute()
}

// Bootstrap initializes the project layout.
func Bootstrap() {
	_, path, _, _ := runtime.Caller(0)
	appTplPath := filepath.Dir(path) + "/templates/app"

	err := filepath.Walk(appTplPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		log.Fatal(err)
	}
}
