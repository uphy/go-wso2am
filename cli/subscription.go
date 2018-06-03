package cli

import (
	"errors"

	"github.com/uphy/go-wso2am"

	"github.com/urfave/cli"
)

func (c *CLI) subscription() cli.Command {
	return cli.Command{
		Name:    "subscription",
		Aliases: []string{"s"},
		Usage:   "Subscription management command",
		Subcommands: cli.Commands{
			c.subscriptionList(),
			c.subscriptionInspect(),
			c.subscriptionBlock(),
			c.subscriptionUnblock(),
		},
	}
}

func (c *CLI) subscriptionList() cli.Command {
	return cli.Command{
		Name:      "list",
		Aliases:   []string{"ls", "dir"},
		Usage:     "List subscriptions",
		ArgsUsage: "[API ID]",
		Action: func(ctx *cli.Context) error {
			resp, err := c.client.SubscriptionsByAPI(ctx.Args().First(), nil)
			if err != nil {
				return err
			}
			f := newTableFormatter("ID", "ApplicationID", "APIID", "Status")
			for _, s := range resp.Subscriptions() {
				f.Row(s.ID, s.ApplicationID, s.APIIdentifier, s.Status)
			}
			return nil
		},
	}
}

func (c *CLI) subscriptionInspect() cli.Command {
	return cli.Command{
		Name:      "inspect",
		Aliases:   []string{"show", "cat"},
		Usage:     "Inspect the subscription",
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			s, err := c.client.Subscription(id)
			if err != nil {
				return err
			}
			return c.inspect(s)
		},
	}
}

func (c *CLI) subscriptionBlock() cli.Command {
	return cli.Command{
		Name:  "block",
		Usage: "Block the subscription",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "prodonly",
				Usage: "Block production subscription only",
			},
		},
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			var state wso2am.SubscriptionBlockState
			if ctx.Bool("prodonly") {
				state = wso2am.SubscriptionBlockStateProdOnlyBlocked
			} else {
				state = wso2am.SubscriptionBlockStateBlocked
			}
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			_, err := c.client.BlockSubscription(id, state)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func (c *CLI) subscriptionUnblock() cli.Command {
	return cli.Command{
		Name:      "unblock",
		Usage:     "Unblock the subscription",
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			_, err := c.client.UnblockSubscription(id)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
