package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/reloonfire/rof-language/helpers"
	"github.com/reloonfire/rof-language/rof"
)

var (
	hadError = false
)

func pnc(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	if len(os.Args) > 1 {
		// file passato inizio con l'interpretazione
		// gRof()
		sc := new(rof.Scanner)
		parser := new(rof.Parser)
		s, _ := ioutil.ReadFile("test.rof")
		sc.Source = string(s)
		tokens := sc.Scan()
		if sc.HadError {
			return
		}
		fmt.Println("[DEBUG] Tokens -> ", tokens)
		parser.Tokens = tokens
		expr := parser.Parse()
		fmt.Println(expr)

		//fmt.Println("Args: ", os.Args[1:])
	} else {
		// Interactive mode
		for {
			hadError = false
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			cmd, err := reader.ReadString('\n')
			fmt.Println("[DEBUG] cmd = ", cmd)
			pnc(err)
			switch strings.Split(cmd, " ")[0] {
			case "quit":
				fmt.Println("\nGoodBye!")
				return
			case "run ":
				path := strings.Split(cmd, " ")[1]

				if path == "" {
					helpers.ReportError(1, "cannot find path to executable file")
					break
				}
				fmt.Println("Path to file = ", path)
			}
		}
	}
}
