package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
  app := cli.NewApp()
  app.Name = "mdcli"
  app.Usage = "Say hello"
  app.Action = func(c *cli.Context) error {
    println("Hello World!")
    return nil
  }

  app.Run(os.Args)
}
