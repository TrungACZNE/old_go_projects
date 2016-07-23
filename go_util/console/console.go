package console_helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetUserChoiceFromMenu(text string, menu []string) string {
	for k, v := range menu {
		fmt.Printf("  %3v. %s\n", k, v)
	}
	fmt.Print(text)

	for {
		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		i, err := strconv.Atoi(strings.TrimSpace(choice))
		if err != nil || i < 0 || i > len(menu)-1 {
			fmt.Println("Must be an integer between 0 and", len(menu)-1)
		} else {
			return menu[i]
		}
	}
}
