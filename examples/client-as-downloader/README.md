# use go-titan-client to download 

By the end of this tutorial, you will learn how to:

- create Downloader
- download separate block or files from titan or specified gateway with cid

## download principle

first, go to titan network to find the content. if it cannot be found, 

go to the specified gateway or the default gateway to find the content

## set custom gateway
`
titan_client.WithCustomGatewayUrlOption("http://127.0.0.1:5001")
`

this parameter can set the gateway you specified, default is `https://ipfs.io/ipfs/`

but If the default gateway is used, scientific Internet access is required.

