package cmd

import (
	"time"
	"fmt"
)

func Help() error {
	fmt.Println("you are beyond help.")
	time.Sleep(16e8)
	fmt.Println("just kidding...")
	time.Sleep(8e8)
	fmt.Println("commands available: http, generate-migration, help(this)")
	return nil
}