package cli

import (
	"errors"
	"fmt"
	"os"

	wso2am "github.com/uphy/go-wso2am"
	"github.com/urfave/cli"
)

func (c *CLI) api() cli.Command {
	return cli.Command{
		Name:    "api",
		Aliases: []string{"a"},
		Usage:   "API management command",
		Subcommands: cli.Commands{
			c.apiList(),
			c.apiChangeStatus(),
			c.apiDelete(),
			c.apiInspect(),
			c.apiSwagger(),
			c.apiUpdateSwagger(),
			c.apiUploadThumbnail(),
			c.apiThumbnail(),
			c.apiCreate(true),
			c.apiCreate(false),
		},
	}
}

func (c *CLI) apiList() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "dir"},
		Usage:   "List APIs",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "query,q",
				Value: "",
			},
		},
		Action: func(ctx *cli.Context) error {
			var query = ctx.String("query")
			resp, err := c.client.SearchAPIs(query, nil)
			if err != nil {
				return err
			}
			f := newTableFormatter("ID", "Name", "Version", "Description", "Status")
			for _, api := range resp.APIs() {
				f.Row(api.ID, api.Name, api.Version, api.Description, api.Status)
			}
			f.Flush()
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

func (c *CLI) apiInspect() cli.Command {
	return cli.Command{
		Name:      "inspect",
		Aliases:   []string{"show", "cat"},
		Usage:     "Inspect the API",
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			api, err := c.client.API(id)
			if err != nil {
				return err
			}
			return c.inspect(api)
		},
	}
}

func (c *CLI) apiSwagger() cli.Command {
	return cli.Command{
		Name:      "swagger",
		Usage:     "Inspect the API definition",
		ArgsUsage: "ID",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			def, err := c.client.APIDefinition(id)
			if err != nil {
				return err
			}
			return c.inspect(def)
		},
	}
}

func (c *CLI) apiUpdateSwagger() cli.Command {
	return cli.Command{
		Name:      "update-swagger",
		Usage:     "Update the API definition",
		ArgsUsage: "ID SWAGGERFILE",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 2 {
				return errors.New("ID and SWAGGERFILE are required")
			}
			id := ctx.Args().Get(0)
			def, err := wso2am.NewAPIDefinitionFromFile(ctx.Args().Get(1))
			if err != nil {
				return err
			}
			if _, err := c.client.UpdateAPIDefinition(id, def); err != nil {
				return err
			}
			return nil
		},
	}
}

func (c *CLI) apiThumbnail() cli.Command {
	return cli.Command{
		Name:  "thumbnail",
		Usage: "Download the thumbnail",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("ID is required")
			}
			id := ctx.Args().Get(0)
			return c.client.Thumbnail(id, os.Stdout)
		},
	}
}

func (c *CLI) apiUploadThumbnail() cli.Command {
	return cli.Command{
		Name:      "upload-thumbnail",
		Usage:     "Upload the thumbnail",
		ArgsUsage: "ID FILE",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 2 {
				return errors.New("ID and FILE are required")
			}
			id := ctx.Args().Get(0)
			file := ctx.Args().Get(1)
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := c.client.UploadThumbnail(id, f); err != nil {
				return err
			}
			return nil
		},
	}
}

func (c *CLI) apiCreate(update bool) cli.Command {
	var commandName string
	var commandAliases []string
	var commandUsage string
	var commandArgsUsage string
	if update {
		commandName = "update"
		commandUsage = "Update the API"
		commandArgsUsage = "ID"
	} else {
		commandName = "create"
		commandAliases = []string{"new"}
		commandUsage = "Create the API"
	}
	return cli.Command{
		Name:      commandName,
		Aliases:   commandAliases,
		Usage:     commandUsage,
		ArgsUsage: commandArgsUsage,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "definition",
			},
			cli.StringFlag{
				Name: "name",
			},
			cli.StringFlag{
				Name: "context",
			},
			cli.StringFlag{
				Name: "version",
			},
			cli.StringFlag{
				Name:  "production-url",
				Value: "http://localhost/",
			},
			cli.StringFlag{
				Name:  "sandbox-url",
				Value: "http://localhost/",
			},
			cli.StringFlag{
				Name: "gateway-env",
			},
			cli.BoolFlag{
				Name: "publish,P",
			},
		},
		Action: func(ctx *cli.Context) error {
			if update {
				if ctx.NArg() != 1 {
					return errors.New("APIID is required")
				}
			} else {
				if err := c.checkRequiredParameters(ctx, "definition", "name", "context", "version", "production-url", "gateway-env"); err != nil {
					return err
				}
			}

			var api *wso2am.APIDetail
			if update {
				id := ctx.Args().First()
				a, err := c.client.API(id)
				if err != nil {
					return err
				}
				api = a
			} else {
				api = c.client.NewAPI()
			}

			if ctx.IsSet("definition") {
				swaggerFile := ctx.String("definition")
				def, err := wso2am.NewAPIDefinitionFromFile(swaggerFile)
				if err != nil {
					return err
				}
				api.Definition = def
			}
			if ctx.IsSet("name") {
				api.Name = ctx.String("name")
			}
			if ctx.IsSet("context") {
				api.Context = ctx.String("context")
			}
			if ctx.IsSet("version") {
				api.Version = ctx.String("version")
			}
			if ctx.IsSet("gateway-env") {
				api.GatewayEnvironments = ctx.String("gateway-env")
			}

			// endpoint config
			if ctx.IsSet("production-url") || ctx.IsSet("sandbox-url") {
				endpointConfig := &wso2am.APIEndpointConfig{
					Type: "http",
				}
				var productionURL = ctx.String("production-url")
				var sandboxURL = ctx.String("sandbox-url")
				endpointConfig.ProductionEndpoints = &wso2am.APIEndpoint{
					URL: productionURL,
				}
				if sandboxURL != "" {
					endpointConfig.SandboxEndpoints = &wso2am.APIEndpoint{
						URL: sandboxURL,
					}
				}
				api.SetEndpointConfig(endpointConfig)
			}

			// call API
			var res *wso2am.APIDetail
			var err error
			if update {
				res, err = c.client.UpdateAPI(api)
			} else {
				res, err = c.client.CreateAPI(api)
			}
			if err != nil {
				return err
			}

			// print the ID of the created API
			if !update {
				fmt.Println(res.ID)
			}

			// publish
			if ctx.Bool("publish") {
				return c.client.ChangeAPIStatus(res.ID, wso2am.APIActionPublish)
			}
			return nil
		},
	}
}
