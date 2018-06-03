# go-wso2am

WSO2 API Manager product api client.

## CLI

### Install

```bash
$ go get -u github.com/uphy/go-wso2am/wso2am-cli
```

### Usage

See command help.

```bash
$ wso2am-cli --help
```

### Examples

List APIs:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api listID
Name                Version             Description                                                       Status90ad3d5b-e535-4cb7-abfe-608e28de16b6 PizzaShackAPI       1.0.0               This is a simple API for Pizza Shack online pizza delivery store. PUBLISHED
```

Inspect API:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api inspect 90ad3d5b-e535-4cb7-abfe-608e28de16b6
{
  "id": "90ad3d5b-e535-4cb7-abfe-608e28de16b6",
  "name": "PizzaShackAPI",
  "description": "This is a simple API for Pizza Shack online pizza delivery store.",
  "context": "/pizzashack",
  "version": "1.0.0",
  "provider": "admin",
  "status": "PUBLISHED",
  "thumbnailUri": "/apis/90ad3d5b-e535-4cb7-abfe-608e28de16b6/thumbnail",
...
}
```

Create and publish API:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api create \
    --name myapi \
    --context myapi \
    --version 1.0 \
    --definition ./swagger.json
    --production-url http://localhost/ \
    --sandbox-url http://localhost/ \
    --gateway-env "Production and Sandbox" \
    --publish
f9b058f7-af45-4973-91c9-5de510b71f39
```

Update the published API(change the gateway environment):

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api update \
    --gateway-env "Production" \
    f9b058f7-af45-4973-91c9-5de510b71f39
```

Update the swagger definition:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api update-swagger f9b058f7-af45-4973-91c9-5de510b71f39 ./swagger.json
```

Upload the thumbnail:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api upload-thumbnail f9b058f7-af45-4973-91c9-5de510b71f39 ./icon.jpeg
```

Download the thumbnail

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api thumbnail f9b058f7-af45-4973-91c9-5de510b71f39 > icon.jpeg
```

Delete the API:

```bash
$ WSO2_USERNAME=user1 WSO2_PASSWORD=user1 wso2am-cli api delete f9b058f7-af45-4973-91c9-5de510b71f39
```