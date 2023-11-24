package main

import (
	"fmt"
	"os"

	"github.com/YuyaAbo/tellme"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("引数が足りません"))
		os.Exit(1)
	}
	query := os.Args[1]

	if err := tellme.Run(query); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
