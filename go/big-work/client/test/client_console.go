package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InputParam struct {
	Command string
	Param   []string
}
type ClientConsole struct {
	chInput chan *InputParam
}

func NewClientConsole() *ClientConsole {
	c := &ClientConsole{}
	return c
}

func (c *ClientConsole) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input err ,check your input and  try again !!!")
			continue
		}
		readString = strings.ReplaceAll(readString, "\n", "")
		split := strings.Split(readString, " ")
		if len(split) == 0 {
			fmt.Println("input err, check your input and  try again !!! ")
			continue
		}
		in := &InputParam{
			Command: split[0],
			Param:   split[1:],
		}
		c.chInput <- in
	}
}
