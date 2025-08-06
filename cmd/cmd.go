package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Config struct {
	Args []string
}

func commandUser(config *Config) {
	if len(config.Args) == 0 {
		fmt.Println("Expected Username")
	}
}
