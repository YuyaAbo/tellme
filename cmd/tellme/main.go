package main

import (
	"fmt"
	"os"

	"github.com/YuyaAbo/tellme"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("引数が足りません")
	}
	object := os.Args[1]

	ai, err := tellme.NewAI()
	if err != nil {
		return err
	}

	answer, err := ai.Describes(object)
	if err != nil {
		return err
	}

	fmt.Println(answer)
	return nil
}
