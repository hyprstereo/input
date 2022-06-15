package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hyprstereo/go-dao/encoding/json"
	"github.com/hyprstereo/input"
)

func main() {
	if len(os.Args) > 1 {
		in := strings.Join(os.Args[1:], " ")
		if in != "" {
			inp := input.NewInput()
			if _, er := inp.Read("${command:String}: ${args:Int}, name:${n:String}", strings.NewReader(in)); er == nil {

				fmt.Println(json.Encode(inp.All(), true).String())
			}
		}
	}
}
