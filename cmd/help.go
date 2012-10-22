package cmd

import (
	"fmt"
	"time"
)

func printAvailableCommands() {
	fmt.Println("commands available: http, generate-migration, help")
}

func MissingCommandHelp() error {
	printAvailableCommands()
	return nil
}

func ExplicitCallForHelp() error {

	fmt.Println("you are beyond help.")
	time.Sleep(2e9)
	fmt.Println("just kidding...")
	time.Sleep(1e9)

	printAvailableCommands()

	return nil
}
