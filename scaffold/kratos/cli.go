package kratos

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

type CliManager[T any] struct {
	app      *cli.App
	provider *T
}

func NewCli[T any](provider *T) CliManager[T] {
	return CliManager[T]{
		app:      cli.NewApp(),
		provider: provider,
	}
}

type CliRegisterOption struct {
	Timeout time.Duration
}

func (c *CliManager[T]) RegisterCmd(cmd func(provider *T) *cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd(c.provider))
}

func (c *CliManager[T]) Run() {
	if err := c.app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
