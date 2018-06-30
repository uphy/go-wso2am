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
			return list(func(entryc chan<- interface{}, errc chan<- error, done <-chan struct{}) {
				c.client.SearchAPIsRaw(query, entryc, errc, done)
			}, func(table *TableFormatter) {
				table.Header("ID", "Name", "Version", "Description", "Status")
			}, func(entry interface{}, table *TableFormatter) {
				api := c.client.ConvertToAPI(entry)
				table.Row(api.ID, api.Name, api.Version, api.Description, api.Status)
			})
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
	flags := []cli.Flag{
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
			Name: "provider",
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
		cli.StringSliceFlag{
			Name: "visible-role",
		},
	}
	if update {
		commandName = "update"
		commandUsage = "Update the API"
		commandArgsUsage = "ID"
	} else {
		commandName = "create"
		commandAliases = []string{"new"}
		commandUsage = "Create the API"
		flags = append(flags, cli.BoolFlag{
			Name: "update",
		})
	}
	return cli.Command{
		Name:      commandName,
		Aliases:   commandAliases,
		Usage:     commandUsage,
		ArgsUsage: commandArgsUsage,
		Flags:     flags,
		Action: func(ctx *cli.Context) error {
			if update {
				if ctx.NArg() != 1 {
					return errors.New("APIID is required")
				}
				unmodifiableFlags := []string{"name", "version", "context", "provider", "state"}
				for _, f := range unmodifiableFlags {
					if ctx.IsSet(f) {
						return fmt.Errorf(`"Cannot update %v"`, unmodifiableFlags)
					}
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
			if ctx.IsSet("provider") {
				api.Provider = ctx.String("provider")
			}
			if ctx.IsSet("visible-role") {
				api.Visibility = wso2am.APIVisibilityRestricted
				api.VisibleRoles = ctx.StringSlice("visible-role")
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

			// if "--update" is specified with create command, find the API ID and update it.
			updateOrCreate := ctx.Bool("update")
			if updateOrCreate {
				// find API ID by context and version
				a, err := c.findAPIByContextVersion(api.Context, api.Version)
				if err != nil {
					return err
				}
				api.ID = a.ID
				/*
					apis, err := c.client.SearchAPIs(fmt.Sprintf("context:%s", api.Context), nil)
					if err != nil {
						return err
					}
					for _, a := range apis.APIs() {
						if a.Version == api.Version {
							api.ID = a.ID
							break
						}
					}
				*/
			}

			// call API
			var res *wso2am.APIDetail
			var err error
			if update || (updateOrCreate && api.ID != "") {
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

func (c *CLI) findAPIByContextVersion(context, version string) (*wso2am.API, error) {
	result, err := c.client.SearchResultToSlice(func(entryc chan<- interface{}, errc chan<- error, done <-chan struct{}) {
		c.client.SearchAPIsRaw(fmt.Sprintf("context:%s", context), entryc, errc, done)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		api := c.client.ConvertToAPI(v)
		if api.Version == version {
			return api, nil
		}
	}
	return nil, nil
}
