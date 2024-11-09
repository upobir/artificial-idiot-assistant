package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/upobir/artificial-idiot-assistant/internal/assistant"
)

func Run(astnt *assistant.Assistant) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("you > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		startTime := time.Now()

		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Exiting...")
			break
		}

		astnt.AddUserMessage(input)

		fmt.Printf("aia > ")
		for part := range astnt.FetchAssistantMessage() {
			if part.Err != nil {
				return part.Err
			}

			if part.Kind == assistant.Action {
				fmt.Printf("%s \naia > ", part.Content)
			} else if part.Kind == assistant.MessagePart {
				output := strings.ReplaceAll(part.Content, "\n", "\n    ")
				fmt.Printf("%s", output)
			} else {
				return fmt.Errorf("unknown repsonse kind: %v", part.Kind)
			}
		}
		fmt.Println()
		endTime := time.Now()
		fmt.Printf("(%0.3f seconds)\n", endTime.Sub(startTime).Seconds())
		fmt.Println()
	}

	return nil
}
