package main

import (
	"context"
	"fmt"
	"socket-share/internal/discovery"
	"socket-share/internal/registry"
	fs "socket-share/internal/share"
)

// App struct
type App struct {
	ctx context.Context
	dm  *discovery.DiscoveryModule
	fr  *registry.FileRegistry
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		dm: discovery.NewDiscoveryModule(),
		fr: registry.NewFileRegistry(),
	}
}

// startup is called at application startup
func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
	go fs.StartFileServer()
	go app.fr.SyncRead()
	app.dm.Start()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Greet2(name string) string {
	return fmt.Sprintf("Hello %s", name)
}
