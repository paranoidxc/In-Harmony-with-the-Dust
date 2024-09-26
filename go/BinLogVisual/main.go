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
	USE         = "USE"
	UPDATE      = "UPDATE"
	INSERT      = "INSERT"
	DELETE      = "DELETE"
	HASH_UPDATE = "### UPDATE"
	HASH_INSERT = "### INSERT"
	HASH_DELETE = "### DELETE"
)

const (
	TYPE_TIMESTAMP   = 998
	TYPE_TERMINAL    = 999
	TYPE_USE         = 1000
	TYPE_UPDATE      = 1
	TYPE_INSERT      = 2
	TYPE_DELETE      = 3
	TYPE_HASH_UPDATE = 101
	TYPE_HASH_INSERT = 102
	TYPE_HASH_DELETE = 103
)

var mapCmdStringToIdx map[string]int

func init() {
	mapCmdStringToIdx = make(map[string]int)

	mapCmdStringToIdx[TERMINAL] = TYPE_TERMINAL
	mapCmdStringToIdx[USE] = TYPE_USE

	mapCmdStringToIdx[UPDATE] = TYPE_UPDATE
	mapCmdStringToIdx[HASH_UPDATE] = TYPE_HASH_UPDATE
	mapCmdStringToIdx[INSERT] = TYPE_INSERT
	mapCmdStringToIdx[HASH_INSERT] = TYPE_HASH_INSERT
	mapCmdStringToIdx[DELETE] = TYPE_DELETE
	mapCmdStringToIdx[HASH_DELETE] = TYPE_HASH_DELETE
}

type Command struct {
	Db        string
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

	//isSqlFinish := true
	isStartSql := false
	startTime := ""
	db := ""

	command := Command{}

	line := 0
	for scanner.Scan() {
		line += 1

		/*
			if line > 100 {
				break
			}
		*/
		if len(cmds) > 1000 {
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

		if isStartSql {
			terminalType := checkCommandTerminal(s)
			if terminalType > 0 {
				if terminalType == TYPE_TERMINAL {
					isStartSql = false
					//isSqlFinish = true
					cmds = append(cmds, command)
					command = Command{Line: line, Db: db}
				} else {
					cmds = append(cmds, command)
					isStartSql = true
					//isSqlFinish = false
					command = Command{Line: line, Typ: terminalType, Db: db}
					command.Sql = command.Sql + " " + strings.TrimSpace(s)
				}
			} else {
				command.Sql = command.Sql + " " + strings.TrimSpace(s)
				continue
			}
		} else {
			//if !isStartSql {
			if strings.HasPrefix(s, TIMESTAMP) {
				startTime = strings.ReplaceAll(strings.ReplaceAll(s, TIMESTAMP, ""), TERMINAL, "")
			} else if len(s) > 10 {
				if strings.HasPrefix(strings.ToUpper(s[:3]), USE) {
					db = strings.ReplaceAll(s[3:], TERMINAL, "")
				} else {
					cmdType := getCommandType(prefix)
					if cmdType > 0 {
						isStartSql = true
						command = Command{
							Db:        db,
							Line:      line,
							Typ:       cmdType,
							Sql:       s,
							StartTime: startTime,
						}
						//isSqlFinish = false
					}
				}
			}
		}
	}

	for idx, cmd := range cmds {
		fmt.Printf("--------------------------------- \nCommand %d:%+v\n", idx+1, cmd)
	}
}
