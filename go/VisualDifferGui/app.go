package main

import (
	"context"
	"encoding/json"
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

type CompareForJs struct {
	Source         string
	Dest           string
	Tpo            int
	Tips           string
	Diff           bool
	SingleFileDiff string
	Del            []string
	Add            []string
	Change         map[string]string
}

func (a *App) xCallCompare(oldRootPath, newRootPath string) string {
	c := CompareForJs{}
	c.Change = make(map[string]string)

	c.Del = append(c.Del, "a.txt")
	c.Del = append(c.Del, "b.txt")

	c.Add = append(c.Add, "b.txt")
	c.Add = append(c.Add, "b.txt")

	c.Change["a.txt"] = "+dsdfs"
	c.Change["b.txt"] = "+dsdfs"

	s, _ := json.Marshal(c)
	return string(s)
}

func (a *App) CallCompare(oldRootPath, newRootPath string) string {
	//func (a *App) xCallCompare(oldRootPath, newRootPath string) map[string]map[string]string {
	//runtime.LogInfo(a.ctx, oldRootPath)
	//runtime.LogInfo(a.ctx, newRootPath)
	//runtime.LogInfo(a.ctx, strconv.Itoa(compareType))
	//var err error
	//var changed string

	c := &CompareForJs{
		Source: oldRootPath,
		Dest:   newRootPath,
	}

	//mapStr := make(map[string]map[string]string)

	isOldRootPathDir, _ := IsDir(oldRootPath)
	isNewRootPathDir, _ := IsDir(newRootPath)

	if isOldRootPathDir && !isNewRootPathDir || !isOldRootPathDir && isNewRootPathDir {
		a.MessageBox(ErrorMsg)
		return ""
	}

	if isOldRootPathDir && isNewRootPathDir {
		compare, err := DoCompareFolder(oldRootPath, newRootPath)
		if err != nil {
		}
		//changed = LogInfoCompare(compare)
		c = LogInfoCompare(compare)
	}

	if !isOldRootPathDir && !isNewRootPathDir {
		c = DoCompareFileWrap(oldRootPath, newRootPath)
	}

	//fmt.Println(Yellow + ChangedPrefix + Reset)
	//return fmt.Sprintf("%s", changed)
	jsonStr, _ := json.Marshal(c)
	return string(jsonStr)
}

func (a *App) MessageBox(str string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "提示信息",
		Message: str,
	})
}
