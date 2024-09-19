package main

import (
	"fmt"
	"strings"
)

var (
	Red     = ""
	Green   = ""
	Yellow  = ""
	Reset   = ""
	Magenta = ""

	CreatedPrefix = "=== 新建 ==="
	RemovedPrefix = "=== 删除 ==="
	ChangedPrefix = "=== 修改 ==="

	ErrorMsg = "只能比对文件夹或者文件"
)

type Compare struct {
	RemovedFiles []string
	CreatedFiles []string
	Changed      [][]byte
}

// DoCompareFolder compare folder
func DoCompareFolder(oldRootPath, newRootPath string) (*Compare, error) {
	compare := &Compare{}

	oldTree, err := GetTree(oldRootPath)
	if err != nil {
		return nil, err
	}

	newTree, err := GetTree(newRootPath)
	if err != nil {
		return nil, err
	}

	var sameTree = make(Tree)
	for path := range newTree {
		var _, isExit = oldTree[path]
		if isExit {
			sameTree[path] = ""

			delete(newTree, path)
			delete(oldTree, path)
		}
	}

	for path := range oldTree {
		compare.RemovedFiles = append(compare.RemovedFiles, path)
	}

	for path := range newTree {
		compare.CreatedFiles = append(compare.CreatedFiles, path)
	}

	for path := range sameTree {
		pathOldFile := oldRootPath + path
		pathNewFile := newRootPath + path

		diff, err := DoCompareFile(pathOldFile, pathNewFile)
		if err != nil {
			return nil, err
		}
		if diff != nil {
			compare.Changed = append(compare.Changed, diff)
		}
	}

	return compare, nil
}

// DoCompareFile compare file
func DoCompareFile(pathOldFile, pathNewFile string) ([]byte, error) {
	oldFile, err := GetFileInfo(pathOldFile)
	if err != nil {
		return nil, err
	}
	if oldFile.IsDir {
		return nil, nil
	}

	newFile, err := GetFileInfo(pathNewFile)
	if err != nil {
		return nil, err
	}
	if newFile.IsDir {
		return nil, nil
	}

	diff := Diff(pathOldFile, oldFile.Data, pathNewFile, newFile.Data)
	return diff, nil
}

// LogInfoCompare logInfo compare
func LogInfoCompare(compare *Compare) string {
	strs := []string{}
	if compare.RemovedFiles != nil {
		//fmt.Println(Red + RemovedPrefix + Reset)
		strs = append(strs, Red+RemovedPrefix+Reset)
		for _, path := range compare.RemovedFiles {
			//fmt.Printf("+ %s\n", path)
			strs = append(strs, fmt.Sprintf("+ %s\n", path))
		}
	}

	if compare.CreatedFiles != nil {
		//fmt.Println()
		//fmt.Println(Green + CreatedPrefix + Reset)

		strs = append(strs, Green+CreatedPrefix+Reset)
		for _, path := range compare.CreatedFiles {
			//fmt.Printf("+ %s\n", path)
			strs = append(strs, fmt.Sprintf("+ %s\n", path))
		}
	}

	if compare.Changed != nil {
		//fmt.Println()
		//fmt.Println(Yellow + ChangedPrefix + Reset)
		strs = append(strs, Yellow+ChangedPrefix+Reset)
		for _, change := range compare.Changed {
			//fmt.Printf("%s", change)
			strs = append(strs, fmt.Sprintf("+ %s\n", change))
		}
	}

	return strings.Join(strs, "\r\n")
}
