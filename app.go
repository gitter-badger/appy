package appy

import (
	"html/template"
	"net/http"
	"os"
)

type (
	// App is the core of appy framework which determines how an application is driven.
	App struct {
		cmd       *Cmd
		config    *Config
		logger    *Logger
		server    *Server
		dbManager *DbManager
		support   *Support
	}
)

const (
	// DebugBuild tends to be slow as it includes debug lvl logging which is more verbose.
	DebugBuild = "debug"

	// ReleaseBuild tends to be faster as it excludes debug lvl logging.
	ReleaseBuild = "release"

	// VERSION follows semantic versioning to indicate the framework's release status.
	VERSION = "0.1.0"

	_description = "An opinionated productive web framework that helps scaling business easier."
)

var (
	// Build is the current build type for the application, can be `debug` or `release`. Please take note that this
	// value will be updated to `release` by `go run . build` command.
	Build = DebugBuild
)

func init() {
	if os.Getenv("APPY_ENV") == "" {
		os.Setenv("APPY_ENV", "development")
	}
}

// NewApp initializes App instance that comes with:
//
// cmd - provides appy's built-in commands and allow custom command constructing
// config - provides appy's global configuration
// logger - provides logger
// server - provides the capability to serve HTTP/GRPC requests
// dbManager - manages the databases along with their pool connections
// support - provides utility helpers/extensions
func NewApp(assets http.FileSystem, viewHelper template.FuncMap) *App {
	cmd := NewCmd()
	support := NewSupport()
	logger := NewLogger(Build)
	config := NewConfig(Build, logger, support, assets)
	dbManager := NewDbManager(logger, support)
	server := NewServer(config, logger, support, assets, viewHelper)

	if Build == DebugBuild {
		cmd.AddCommand()
	}

	cmd.AddCommand(
		newConfigDecryptCommand(config, logger, support),
		newConfigEncryptCommand(config, logger, support),
		newDbCreateCommand(config, dbManager, logger),
		newDbDropCommand(config, dbManager, logger),
		newSecretCommand(logger),
		newServeCommand( dbManager, server),
		newSSLCleanCommand(logger, server),
		newSSLSetupCommand(logger, server),
	)

	return &App{
		cmd:       cmd,
		config:    config,
		dbManager: dbManager,
		logger:    logger,
		server:    server,
		support:   support,
	}
}

// Cmd returns the app's Cmd instance.
func (a App) Cmd() *Cmd {
	return a.cmd
}

// Config returns the app's Config instance.
func (a App) Config() *Config {
	return a.config
}

// DbManager returns the app's DbManager instance.
func (a App) DbManager() *DbManager {
	return a.dbManager
}

// Logger returns the app's Logger instance.
func (a App) Logger() *Logger {
	return a.logger
}

// Server returns the app's Server instance.
func (a App) Server() *Server {
	return a.server
}

// Support returns the app's Support instance.
func (a App) Support() *Support {
	return a.support
}

// Run starts the application.
func (a App) Run() {
	a.Cmd().Execute()
}
