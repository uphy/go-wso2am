package cli

import (
	"errors"
	"fmt"

	wso2am "github.com/uphy/go-wso2am"
	"github.com/urfave/cli"
)

func (c *CLI) api() cli.Command {
	return cli.Command{
		Name:  "api",
		Usage: "API management command",
		Subcommands: cli.Commands{
			c.apiList(),
			c.apiChangeStatus(),
			c.apiDelete(),
		},
	}
}

func (c *CLI) apiList() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "dir"},
		Usage:   "List APIs",
		Action: func(ctx *cli.Context) error {
			resp, _ := c.client.APIs(nil)
			for _, api := range resp.APIs() {
				fmt.Println(api.ID)
				c.client.ChangeAPIStatus(api.ID, wso2am.APIActionPublish)
			}
			return nil
		},
	}
}

func (c *CLI) apiChangeStatus() cli.Command {
	return cli.Command{
		Name:  "change-status",
		Usage: "Change API status",
		Description: fmt.Sprintf(`Change API status.

Available actions are:
- %s
- %s
- %s
- %s
- %s
- %s
- %s
- %s
`, wso2am.APIActionPublish, wso2am.APIActionDeployAsPrototype, wso2am.APIActionDemoteToCreated, wso2am.APIActionDemoteToPrototyped, wso2am.APIActinBlock, wso2am.APIActinDeprecate, wso2am.APIActionRePublish, wso2am.APIActionRetire),
		ArgsUsage: "ID ACTION",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 2 {
				return errors.New("ID and ACTION are required")
			}
			id := ctx.Args().Get(0)
			action := ctx.Args().Get(1)
			return c.client.ChangeAPIStatus(id, wso2am.APIAction(action))
		},
	}
}

func (c *CLI) apiDelete() cli.Command {
	return cli.Command{
		Name:      "delete",
		Aliases:   []string{"del", "rm"},
		Usage:     "Delete the API",
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			return c.client.DeleteAPI(id)
		},
	}
}