package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/alexei-led/pumba/pkg/chaos"
	"github.com/alexei-led/pumba/pkg/chaos/docker"
)

type killContext struct {
	context context.Context
}

// NewKillCLICommand initialize CLI kill command and bind it to the killContext
func NewKillCLICommand(ctx context.Context) *cli.Command {
	cmdContext := &killContext{context: ctx}
	return &cli.Command{
		Name: "kill",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "signal, s",
				Usage: "termination signal, that will be sent by Pumba to the main process inside target container(s)",
				Value: docker.DefaultKillSignal,
			},
			cli.IntFlag{
				Name:  "limit, l",
				Usage: "limit to number of container to kill (0: kill all matching)",
				Value: 0,
			},
		},
		Usage:       "kill specified containers",
		ArgsUsage:   fmt.Sprintf("containers (name, list of names, or RE2 regex if prefixed with %q", chaos.Re2Prefix),
		Description: "send termination signal to the main process inside target container(s)",
		Action:      cmdContext.kill,
	}
}

// KILL Command
func (cmd *killContext) kill(c *cli.Context) error {
	message := docker.KillMessage{}
	// get random
	message.Random = c.GlobalBool("random")
	// get dry-run mode
	message.DryRun = c.GlobalBool("dry-run")
	// get interval
	message.Interval = c.GlobalString("interval")
	// get names or pattern
	message.Names, message.Pattern = chaos.GetNamesOrPattern(c)
	// get signal
	message.Signal = c.String("signal")
	// get limit for number of containers to kill
	message.Limit = c.Int("limit")
	// init kill command
	killCommand, err := docker.NewKillCommand(chaos.DockerClient, message)
	if err != nil {
		return err
	}
	// run kill command
	return chaos.RunChaosCommand(cmd.context, killCommand, message.Interval, message.Random)
}
