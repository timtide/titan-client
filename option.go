package titan_client

type Option func(td *titanDownloader)

// WithCustomGatewayAddressOption custom set gateway url
// eg: http://127.0.0.1:5001 or https://ipfs.io/ipfs/
// If you use the local port as the gateway,
// you need to enable the ipfs node locally
func WithCustomGatewayAddressOption(addr string) Option {
	return func(td *titanDownloader) {
		td.customGatewayAddr = addr
	}
}

func WithLocatorAddressOption(locatorAddr string) Option {
	return func(td *titanDownloader) {
		td.locatorAddr = locatorAddr
	}
}
