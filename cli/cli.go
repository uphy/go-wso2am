package cli

import (
	"encoding/json"
	"fmt"

	wso2am "github.com/uphy/go-wso2am"
	"github.com/urfave/cli"
)

const Version = "0.0.1"

type CLI struct {
	app    *cli.App
	client *wso2am.Client
}

func New() *CLI {
	app := cli.NewApp()
	app.Version = Version
	app.Usage = "WSO2 API Manager product API client"
	c := &CLI{
		app: app,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "url-carbon,ec",
			EnvVar: "WSO2_CARBON_URL",
			Value:  "https://localhost:9443/",
		},
		cli.StringFlag{
			Name:   "url-token,et",
			EnvVar: "WSO2_TOKEN_URL",
			Value:  "https://localhost:8243/",
		},
		cli.StringFlag{
			Name:   "user,u",
			EnvVar: "WSO2_USERNAME",
			Value:  "admin",
		},
		cli.StringFlag{
			Name:   "password,p",
			EnvVar: "WSO2_PASSWORD",
			Value:  "admin",
		},
		cli.StringFlag{
			Name:   "client,c",
			EnvVar: "WSO2_CLIENT_NAME",
			Value:  "wso2am-cli-client",
		},
		cli.StringFlag{
			Name:   "client-id,ci",
			EnvVar: "WSO2_CLIENT_ID",
			Value:  "", // Automatically register client
		},
		cli.StringFlag{
			Name:   "client-secret,cs",
			EnvVar: "WSO2_CLIENT_SECRET",
			Value:  "", // Automatically register client
		},
		cli.StringFlag{
			Name:   "apiversion,av",
			EnvVar: "WSO2_API_VERSION",
			Value:  wso2am.DefaultAPIVersion, // Automatically register client
		},
	}
	app.Before = func(ctx *cli.Context) error {
		carbonURL := ctx.String("url-carbon")
		tokenURL := ctx.String("url-token")
		user := ctx.String("user")
		password := ctx.String("password")
		clientName := ctx.String("client")
		apiVersion := ctx.String("apiversion")
		client, err := wso2am.New(&wso2am.Config{
			EndpointCarbon: carbonURL,
			EndpointToken:  tokenURL,
			ClientName:     clientName,
			UserName:       user,
			Password:       password,
			APIVersion:     apiVersion,
		})
		if err != nil {
			return err
		}
		c.client = client
		return nil
	}

	c.addCommand(c.api())
	c.addCommand(c.subscription())

	return c
}

func (c *CLI) addCommand(cmd cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd)
}

func (c *CLI) checkRequiredParameters(ctx *cli.Context, parameters ...string) error {
	for _, p := range parameters {
		if !ctx.IsSet(p) {
			return fmt.Errorf(`"%s" is not set.  (required flags: %v)`, p, parameters)
		}
	}
	return nil
}

func (c *CLI) inspect(v interface{}) error {
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(d))
	return nil
}

func (c *CLI) Run(args []string) error {
	return c.app.Run(args)
}
