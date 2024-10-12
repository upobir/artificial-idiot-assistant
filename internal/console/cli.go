package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("you > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Exiting...")
			break
		}

		fmt.Printf("aia > %s\n", input)
		fmt.Println()
	}
}
