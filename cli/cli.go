package cli

import (
	wso2am "github.com/uphy/go-wso2am"
	"github.com/urfave/cli"
)

type CLI struct {
	app    *cli.App
	client *wso2am.Client
}

func New() *CLI {
	app := cli.NewApp()
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
	}
	app.Before = func(ctx *cli.Context) error {
		carbonURL := ctx.String("url-carbon")
		tokenURL := ctx.String("url-token")
		user := ctx.String("user")
		password := ctx.String("password")
		client, err := wso2am.New(&wso2am.Config{
			EndpointCarbon: carbonURL,
			EndpointToken:  tokenURL,
			ClientName:     "test",
			UserName:       user,
			Password:       password,
		})
		if err != nil {
			return err
		}
		c.client = client
		return nil
	}

	c.addCommand(c.api())

	return c
}

func (c *CLI) addCommand(cmd cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd)
}

func (c *CLI) Run(args []string) error {
	return c.app.Run(args)
}
