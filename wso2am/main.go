package main

import (
	"fmt"
	"os"

	"github.com/uphy/go-wso2am"
)

func main() {
	c, err := wso2am.New(&wso2am.Config{
		EndpointCarbon: "https://localhost:9443/",
		EndpointToken:  "https://localhost:8243/",
		ClientName:     "test",
		UserName:       "admin",
		Password:       "admin",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(c.GenerateAccessToken("apim:tier_view"))
	resp, _ := c.APIs(nil)
	fmt.Println(resp.APIs())
}
