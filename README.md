# go-titan-client

> download data from titan network or gateway

By the end of this tutorial, you will learn how to:

- download complete file
- one by one download block

## complete file

### download principle

first, go to titan network to find the content. if it cannot be found,

go to the specified gateway or the default gateway to find the content.

default the gateway address is `https://ipfs.io/ipfs`,

default the locator address is `http://39.108.143.56:5000`.

### download step description

if you want to download the complete file according to the cid，

just call `NewDownloader` production `Downloader` object,

then call `Download()` method or call `GetReader()` method.

of course, you can also customize the gateway address and locator address,

`WithCustomGatewayAddressOption()` method to customize the gateway address,

`WithLocatorAddressOption()` method to customize the locator address.

example:
```
    c, err := cid.Decode("QmUbaDBz6YKn3dVzoKrLDyupMmyWk5am2QSdgfKsU1RN3N")
    if err != nil {
        return
    }
    err = NewDownloader().Download(context.Background(), c, false, gzip.NoCompression, "./titan.mp4")
    if err != nil {
        return
    }
```

## separate blocks

### download principle

go to titan network to find the content.

default the locator address is `http://39.108.143.56:5000`.

### download step description

if you want to download separate blocks according to the cid，

just call `NewBlockGetter` production `BlockGetter` object,

then, if you just want to get one block, call `GetBlock()` method,

if you want to get many blocks, call `GetBlocks()` method. of course, 

you can also customize locator address, please call `blockdownload` 

folder `WithLocatorAddressOption()` method to customize the locator address.

this function can be added to the ipfs code, 

and titan can be used as a cache to speed up downloading data

example:
```
	c, err := cid.Decode("QmPgaP4SiadmrtFzEVY5aGTCRou5vbMDJCgEaJwuN9Lk4H")
	if err != nil {
		return
	}
	block, err := NewBlockGetter().GetBlock(context.Background(), c)
	if err != nil {
		return
	}
```

## License

MIT license
