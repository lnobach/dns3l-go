package types

import (
	dnstypes "github.com/dta4/dns3l-go/dns/types"
	"github.com/dta4/dns3l-go/state"
)

type CertificateResources struct {
	Certificate string
	Key         string
	Chain       string
	FullChain   string
}

type CertificateClaimInfo struct {
	Name            string
	SubjectAltNames []string
	//TODO implement hints
}

type CAConfigurationContext interface {
	GetStateProvider() state.StateProvider
	GetDNSProvider(provID string) (dnstypes.DNSProvider, bool)
}
