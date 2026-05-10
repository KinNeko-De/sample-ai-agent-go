package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/kinneko-de/sample-ai-agent-go/internal/agent"
)

const quitCommand = "quit"
const welcomeMessage = "Welcome to the company assistant! How can I help you today?"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	a := agent.New()

	fmt.Println(welcomeMessage)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == quitCommand {
			break
		}

		err := a.Chat(input, os.Stdout)
		if err != nil {
			slog.Error("chat error", slog.Any("error", err))
			fmt.Println("Something went wrong. Please try again.")
			continue
		}
	}
}
