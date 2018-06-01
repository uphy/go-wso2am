package main

import (
	"fmt"
	"os"

	"github.com/uphy/go-wso2am/cli"
)

func main() {
	c := cli.New()
	if err := c.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "failed to execute the command: ", err)
		os.Exit(1)
	}
}
