package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SelectOldFolder() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder",
	})

	if err != nil {
	}
	runtime.LogInfo(a.ctx, selection)
	return selection
}

func (a *App) SelectOld() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		/*
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Images (*.png;*.jpg)",
					Pattern:     "*.png;*.jpg",
				}, {
					DisplayName: "Videos (*.mov;*.mp4)",
					Pattern:     "*.mov;*.mp4",
				},
			},
		*/
	})

	if err != nil {
	}
	runtime.LogInfo(a.ctx, selection)
	return selection
}

func (a *App) SelectNew() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpg)",
				Pattern:     "*.png;*.jpg",
			}, {
				DisplayName: "Videos (*.mov;*.mp4)",
				Pattern:     "*.mov;*.mp4",
			},
		},
	})

	if err != nil {
	}

	runtime.LogInfo(a.ctx, selection)
	return selection
}

func (a *App) CallCompare(oldRootPath, newRootPath string) string {
	changed, err := DoCompareFile(oldRootPath, newRootPath)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	//fmt.Println(Yellow + ChangedPrefix + Reset)
	return fmt.Sprintf("%s", changed)
}
