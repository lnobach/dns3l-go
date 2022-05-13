package types

//Can spawn a new instance of a DNS provider.
//Provider-specific config is unmarshaled into the builder's implementing struct.
type DNSProviderBuilder interface {
	NewInstance() (DNSProvider, error)
}

type DNSProviderInfo struct {
	Name        string
	Feature     []string
	ZoneNesting bool
}

type DNSProvider interface {
	GetInfo() *DNSProviderInfo
}