package main

import (
	"code"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format",
				Value:   "stylish",
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			if c.Args().Len() != 2 {
				return errors.New("expected two arguments: filepath1 filepath2")
			}
			result, err := code.GenDiff(c.Args().Get(0), c.Args().Get(1), c.String("format"))
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
