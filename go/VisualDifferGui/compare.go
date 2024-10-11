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

type CompareFile struct {
	OldFileData    string
	NewFileData    string
	SingleFileDiff string
}

type Compare struct {
	RemovedFiles []string
	CreatedFiles []string
	Changed      map[string]string
	FilesChanged map[string]CompareFile
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

	compare.Changed = make(map[string]string)
	compare.FilesChanged = make(map[string]CompareFile)

	for path := range sameTree {
		pathOldFile := oldRootPath + path
		pathNewFile := newRootPath + path

		diffStr, oldFileData, newFileData, err := DoCompareFile(pathOldFile, pathNewFile)
		if err != nil {
			return nil, err
		}
		if diffStr != "" {
			compare.Changed[path] = diffStr
			//compare.Changed = append(compare.Changed, diff)
			compare.FilesChanged[path] = CompareFile{
				OldFileData:    oldFileData,
				NewFileData:    newFileData,
				SingleFileDiff: diffStr,
			}
		}
	}

	return compare, nil
}

// DoCompareFile compare file

func DoCompareFileWrap(pathOldFile, pathNewFile string) *CompareForJs {
	c := &CompareForJs{
		Tpo:    0,
		Source: pathOldFile,
		Dest:   pathNewFile,
	}

	diffStr, old, new, err := DoCompareFile(pathOldFile, pathNewFile)
	if err != nil {
	}
	c.Old = old
	c.New = new
	c.SingleFileDiff = diffStr
	c.Change = make(map[string]string)
	c.Tpo = 0
	if diffStr == "" {
		c.Tips = "文件内容相同"
	} else {
		c.Diff = true
	}

	return c
}

func DoCompareFile(pathOldFile, pathNewFile string) (string, string, string, error) {
	oldFile, err := GetFileInfo(pathOldFile)
	if err != nil {
		return "", "", "", err
	}
	if oldFile.IsDir {
		return "", "", "", nil
	}

	newFile, err := GetFileInfo(pathNewFile)
	if err != nil {
		return "", "", "", err
	}
	if newFile.IsDir {
		return "", "", "", nil
	}

	diff := Diff(pathOldFile, oldFile.Data, pathNewFile, newFile.Data)

	diffStr := string(diff)
	diffStr = strings.Replace(diffStr, "\r\n", "<br/>", -1)
	diffStr = strings.Replace(diffStr, "\n", "<br/>", -1)

	return diffStr, string(oldFile.Data), string(newFile.Data), nil
}

// LogInfoCompare logInfo compare
// func LogInfoCompare(compare *Compare) map[string]map[string]string {
func LogInfoCompare(compare *Compare) *CompareForJs {
	strs := []string{}
	mapStr := make(map[string]map[string]string)

	c := CompareForJs{Tpo: 1}
	c.Change = make(map[string]string)
	c.FilesChange = make(map[string]CompareFile)
	isDiff := false

	if compare.RemovedFiles != nil {
		mapStr["DEL"] = make(map[string]string)
		//fmt.Println(Red + RemovedPrefix + Reset)
		strs = append(strs, Red+RemovedPrefix+Reset)
		for _, path := range compare.RemovedFiles {
			isDiff = true
			//fmt.Printf("+ %s\n", path)
			strs = append(strs, fmt.Sprintf("+ %s\n", path))
			mapStr["DEL"][path] = path
			c.Del = append(c.Del, path)
		}
	}

	if compare.CreatedFiles != nil {
		//fmt.Println()
		//fmt.Println(Green + CreatedPrefix + Reset)
		mapStr["NEW"] = make(map[string]string)
		strs = append(strs, Green+CreatedPrefix+Reset)
		for _, path := range compare.CreatedFiles {
			isDiff = true
			//fmt.Printf("+ %s\n", path)
			strs = append(strs, fmt.Sprintf("+ %s\n", path))
			//mapStr["NEW"] = append(mapStr["NEW"], path)
			mapStr["NEW"][path] = path
			c.Add = append(c.Add, path)
		}
	}

	if compare.Changed != nil {
		//fmt.Println()
		//fmt.Println(Yellow + ChangedPrefix + Reset)
		mapStr["CHANGE"] = make(map[string]string)
		strs = append(strs, Yellow+ChangedPrefix+Reset)
		isDiff = true
		c.FilesChange = compare.FilesChanged
		// for path, change := range compare.Changed {
		// 	isDiff = true
		// 	//fmt.Printf("%s", change)
		// 	changeStr := fmt.Sprintf("+ %s\n", change)
		// 	//strs = append(strs, fmt.Sprintf("+ %s\n", change))
		// 	//mapStr["CHANGE"][path] = changeStr //= append(mapStr["CHANGE"], path)
		// 	c.Change[path] = changeStr
		// 	c.FilesChange[path] = CompareFile{
		// 		OldFileData:    "",
		// 		NewFileData:    "",
		// 		SingleFileDiff: changeStr,
		// 	}
		// }
	}

	c.Diff = isDiff

	if !isDiff {
		c.Tips = "文件夹内容相同"
	}

	//return strings.Join(strs, "\r\n")
	//return mapStr
	return &c
}

func DetectBinary(path string) bool {
    file, err := os.Open(path)
    if err != nil {
        return false
    }
    defer file.Close()

    r := bufio.NewReader(file)
    buf := make([]byte, 1024)
    n, err := r.Read(buf)

    var white_byte int = 0
    for i := 0; i < n; i++ {
        if (buf[i] >= 0x20 && buf[i] <= 0xff) ||
            buf[i] == 9  ||
            buf[i] == 10 ||
            buf[i] == 13 {
            white_byte++
        } else if buf[i] <= 6 || (buf[i] >= 14 && buf[i] <= 31) {
            return true
        }
    }

    if white_byte >= 1 {
        return false
    }
    return true
}