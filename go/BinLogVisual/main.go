package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	TIMESTAMP   = "SET TIMESTAMP="
	TERMINAL    = "/*!*/;"
	UPDATE      = "UPDATE"
	INSERT      = "INSERT"
	DELETE      = "DELETE"
	HASH_UPDATE = "### UPDATE"
	HASH_INSERT = "### INSERT"
	HASH_DELETE = "### DELETE"
)

const (
	IDX_TIMESTAMP   = 998
	IDX_TERMINAL    = 999
	IDX_UPDATE      = 1
	IDX_INSERT      = 2
	IDX_DELETE      = 3
	IDX_HASH_UPDATE = 101
	IDX_HASH_INSERT = 102
	IDX_HASH_DELETE = 103
)

var mapCmdStringToIdx map[string]int

func init() {
	mapCmdStringToIdx = make(map[string]int)

	mapCmdStringToIdx[TERMINAL] = IDX_TERMINAL

	mapCmdStringToIdx[UPDATE] = IDX_UPDATE
	mapCmdStringToIdx[HASH_UPDATE] = IDX_HASH_UPDATE
	mapCmdStringToIdx[INSERT] = IDX_INSERT
	mapCmdStringToIdx[HASH_INSERT] = IDX_HASH_INSERT
	mapCmdStringToIdx[DELETE] = IDX_DELETE
	mapCmdStringToIdx[HASH_DELETE] = IDX_HASH_DELETE
}

type Command struct {
	Line      int
	Typ       int
	StartTime string
	Sql       string
}

func getCommandType(s string) int {
	r := 0
	switch {
	case strings.HasPrefix(s, UPDATE):
		r = mapCmdStringToIdx[UPDATE]
	case strings.HasPrefix(s, HASH_UPDATE):
		r = mapCmdStringToIdx[HASH_UPDATE]
	case strings.HasPrefix(s, INSERT):
		r = mapCmdStringToIdx[INSERT]
	case strings.HasPrefix(s, HASH_INSERT):
		r = mapCmdStringToIdx[HASH_DELETE]
	case strings.HasPrefix(s, DELETE):
		r = mapCmdStringToIdx[DELETE]
	case strings.HasPrefix(s, HASH_DELETE):
		r = mapCmdStringToIdx[HASH_DELETE]
	}

	return r
}

func checkCommandTerminal(s string) int {
	r := 0
	switch {
	case s == TERMINAL:
		r = mapCmdStringToIdx[TERMINAL]
	case strings.HasPrefix(s, HASH_UPDATE):
		r = mapCmdStringToIdx[HASH_UPDATE]
	case strings.HasPrefix(s, HASH_INSERT):
		r = mapCmdStringToIdx[HASH_INSERT]
	case strings.HasPrefix(s, HASH_DELETE):
		r = mapCmdStringToIdx[HASH_DELETE]
	}

	return r
}

func main() {
	filename := "./binlog.001585.sql"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	cmds := []Command{}

	scanner := bufio.NewScanner(file)

	isSqlFinish := true
	isStartSql := false
	startTime := ""

	command := Command{}

	line := 0
	for scanner.Scan() {
		line += 1

		/*
			if line > 100 {
				break
			}
		*/
		if len(cmds) > 50 {
			break
		}

		s := scanner.Text()

		/*
			fmt.Println("line:", line)
			fmt.Println("string:", s)
		*/

		prefix := ""
		if len(s) >= 10 {
			prefix = strings.ToUpper(s[:10])
		}

		if isStartSql && !isSqlFinish {
			terminalType := checkCommandTerminal(s)
			if terminalType > 0 {
				if terminalType == IDX_TERMINAL {
					isStartSql = false
					isSqlFinish = true
					cmds = append(cmds, command)
					command = Command{Line: line}
				} else {
					cmds = append(cmds, command)
					isStartSql = true
					isSqlFinish = false
					command = Command{Line: line, Typ: terminalType}
					command.Sql = command.Sql + " " + strings.TrimSpace(s)
				}
			} else {
				command.Sql = command.Sql + " " + strings.TrimSpace(s)
				continue
			}
		}

		if !isStartSql {
			if strings.HasPrefix(s, TIMESTAMP) {
				startTime = strings.ReplaceAll(strings.ReplaceAll(s, TIMESTAMP, ""), TERMINAL, "")
			} else if len(s) > 10 {
				cmdType := getCommandType(prefix)
				if cmdType > 0 {
					isStartSql = true
					command = Command{
						Line:      line,
						Typ:       cmdType,
						Sql:       s,
						StartTime: startTime,
					}
					isSqlFinish = false
				}
			}
		}
	}

	for _, cmd := range cmds {
		fmt.Printf("--------------------------------- \nCommand %+v\n", cmd)
	}
}
