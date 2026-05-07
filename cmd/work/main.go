package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/tscolari/work/internal/cli"
)

func main() {
	if err := cli.RunWork(os.Args[1:], os.Stderr); err != nil {
		var ce *cli.Error
		if errors.As(err, &ce) {
			fmt.Fprintln(os.Stderr, ce.Msg)
			os.Exit(ce.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
