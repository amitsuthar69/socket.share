package main

import (
	"context"
	"socket-share/internal/discovery"
	"socket-share/internal/registry"
	fs "socket-share/internal/share"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	go app.fr.SyncRead(ctx)
	app.dm.Start()
}

// OpenFilePicker will open the native file explorer with all file type.
func (app *App) OpenFilePicker() (string, error) {
	return runtime.OpenFileDialog(app.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
}

func (app *App) CreateNewFile(path string) registry.File {
	return app.fr.NewFile(path)
}

func (app *App) DownloadFile(ip, path string) {
	fs.StartFileClient(ip, path)
}

// domReady is called after front-end resources have been loaded
// func (app App) domReady(ctx context.Context) {
// // Add your action here
// }

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
// func (app *App) beforeClose(ctx context.Context) (prevent bool) {
// 	return false
// }

// shutdown is called at application termination
// func (app *App) shutdown(ctx context.Context) {
// // Perform your teardown here
// }
