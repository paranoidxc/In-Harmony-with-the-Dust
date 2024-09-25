package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Typ       int
	StartTime string
	Sql       string
}

func main() {
	filename := "./binlog.001585.sql"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	cmds := []Command{}

	scanner := bufio.NewScanner(file)
	idx := 0

	isSqlFinish := true
	isStartSql := false
	startTime := ""

	command := Command{}
	for scanner.Scan() {
		if idx > 2000 {
			break
		}

		s := scanner.Text()

		if isStartSql && !isSqlFinish {
			if s == "/*!*/;" {
				//fmt.Printf("--------------------------------- \n cmmand %+v\n", command)
				isStartSql = false
				isSqlFinish = true
				cmds = append(cmds, command)
				command = Command{}
			} else {
				command.Sql = command.Sql + " " + strings.TrimSpace(s)
			}
			continue
		}

		if strings.HasPrefix(s, "SET TIMESTAMP=") {
			startTime = s
		} else if len(s) > 10 {
			prefix := strings.ToUpper(s[:10])

			if strings.HasPrefix(prefix, "INSERT") {
				isStartSql = true
			} else if strings.HasPrefix(prefix, "UPDATE") {
				isStartSql = true
			} else if strings.HasPrefix(prefix, "### UPDATE") || strings.HasPrefix(prefix, "### INSERT") {
				cmds = append(cmds, command)
				command = Command{}
				isStartSql = true
			} else if strings.HasPrefix(prefix, "DELETE") {
				isStartSql = true
			} else if strings.HasPrefix(prefix, "UPDATEx") {
				//fmt.Println(s)
				isStartSql = true
			}

			if isStartSql {
				command.Sql = s
				command.StartTime = startTime
				isSqlFinish = false
			}
		}

		idx++
	}

	for _, cmd := range cmds {
		fmt.Printf("--------------------------------- \n cmmand %+v\n", cmd)
	}
}
