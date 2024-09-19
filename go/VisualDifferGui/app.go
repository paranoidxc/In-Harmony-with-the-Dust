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

func (a *App) SelectOld(compareType bool) string {
	if compareType {
		selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Select Folder",
		})

		if err != nil {
		}
		runtime.LogInfo(a.ctx, selection)
		return selection
	} else {
		selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Select File",
		})
		if err != nil {
		}
		runtime.LogInfo(a.ctx, selection)
		return selection
	}
}

func (a *App) SelectNew(compareType bool) string {
	if compareType {
		selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Select Folder",
		})

		if err != nil {
		}
		runtime.LogInfo(a.ctx, selection)
		return selection
	} else {
		selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Select File",
		})
		if err != nil {
		}
		runtime.LogInfo(a.ctx, selection)
		return selection
	}
}

func (a *App) CallCompare(oldRootPath, newRootPath string) string {
	//runtime.LogInfo(a.ctx, oldRootPath)
	//runtime.LogInfo(a.ctx, newRootPath)
	//runtime.LogInfo(a.ctx, strconv.Itoa(compareType))
	//var err error
	var changed string

	isOldRootPathDir, _ := IsDir(oldRootPath)
	isNewRootPathDir, _ := IsDir(newRootPath)

	if isOldRootPathDir && !isNewRootPathDir || !isOldRootPathDir && isNewRootPathDir {
		a.MessageBox(ErrorMsg)
		return "-1"
	}

	if isOldRootPathDir && isNewRootPathDir {
		compare, err := DoCompareFolder(oldRootPath, newRootPath)
		if err != nil {
		}
		changed = LogInfoCompare(compare)
	}

	if !isOldRootPathDir && !isNewRootPathDir {
		changedByte, err := DoCompareFile(oldRootPath, newRootPath)
		if err != nil {
		}
		changed = string(changedByte)
	}

	//fmt.Println(Yellow + ChangedPrefix + Reset)
	//return fmt.Sprintf("%s", changed)
	return changed
}

func (a *App) MessageBox(str string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "提示信息",
		Message: str,
	})
}
