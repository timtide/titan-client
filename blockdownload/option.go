package blockdownload

type Option func(bg *blockGetter)

// WithLocatorAddressOption set your favorite locator url, address or domain name
func WithLocatorAddressOption(locatorAddr string) Option {
	return func(dg *blockGetter) {
		dg.locatorAddr = locatorAddr
	}
}
