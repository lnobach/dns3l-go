package common

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// The ResolveTester can be used to test if a set DNS01 acme challenge has been actually
// placed on a DNS service so it can be successfully validated. It also provides functionality
// to wait for an active DNS challenge record.
type ResolveTester struct {
	DNSCheckServers   []string      `yaml:"DNSCheckServers"`
	ResolveMaxRetries int           `yaml:"ResolveMaxRetries"`
	ResolveTimeout    time.Duration `yaml:"ResolveTimeout"`
}

// WaitForAActive checks for the A record for being placed. If it is not
// yet placed, it waits until a timeout for a valid A record,
// as the placement could be delayed.
func (re *ResolveTester) WaitForAActive(name string, ipv4 net.IP) error {
	return re.waitForActive(name, dns.TypeA, func(dname string, rr []dns.RR) (bool, error) {
		result := rr[0].(*dns.A).A
		if result.Equal(ipv4) {
			return true, nil
		}
		log.Warnf("Domain %s: A record should be %s, but is %s. Retrying", dname, ipv4.String(), result.String())
		return false, nil
	})
}

// WaitForTXTActive checks for the DNS01 challenge for being placed. If it is not
// yet placed, it waits until a timeout for a valid DNS01 challenge,
// as the placement could be delayed.
func (re *ResolveTester) WaitForChallengeActive(name, expectedChallenge string) error {

	dName, err := EnsureAcmeChallengeFormat(name)
	if err != nil {
		return err
	}

	return re.waitForActive(dName, dns.TypeTXT, func(ddname string, rr []dns.RR) (bool, error) {
		result := rr[0].(*dns.TXT).Txt[0]
		if result == expectedChallenge {
			return true, nil
		}
		log.Warnf("Domain %s: Challenge should be %s, but is %s. Retrying", ddname, expectedChallenge, result)
		return false, nil
	})
}

func (re *ResolveTester) waitForActive(dName string, dnstype uint16,
	ah func(dName string, rr []dns.RR) (bool, error)) error {

	err := ValidateDomainName(dName)
	if err != nil {
		return err
	}

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dName, dnstype)

	for i := 0; i < re.ResolveMaxRetries; i++ {
		dnsSocket := getDNSSocketForTry(i, re.DNSCheckServers)
		log.Debugf("Trying to resolve %s with DNS server %s, attempt %d..", dName, dnsSocket, i)
		r, t, err := c.Exchange(&m, dnsSocket)
		if err != nil {
			log.Errorf("Error when contacting DNS server for validation: %v\n", err)
			time.Sleep(re.ResolveTimeout)
			continue
		}
		if len(r.Answer) > 0 {
			log.Infof("Successfully resolved %s, time=%s", dName, t)
			success, err := ah(dName, r.Answer)
			if err != nil {
				return err
			}
			if success {
				return nil
			}
		}
		time.Sleep(re.ResolveTimeout)
	}

	return fmt.Errorf("too many retries to resolve record after setting")
}

func getDNSSocketForTry(tryNum int, dnsSockets []string) string {
	//Round-robin trial of DNS servers
	return dnsSockets[tryNum%len(dnsSockets)]
}

// ConfigureFromEnv configures a ResolveTester object based on environment
// variables set. No fields on this object need to be pre-initialized.
func (re *ResolveTester) ConfigureFromEnv() error {
	var err error
	cfResolveMaxRetriesStr := os.Getenv("DNS3LD_TESTRESOLVE_RETRIES")
	var cfResolveMaxRetries int
	if cfResolveMaxRetriesStr != "" {
		cfResolveMaxRetries, err = strconv.Atoi(cfResolveMaxRetriesStr)
		if err != nil {
			return err
		}
	} else {
		cfResolveMaxRetries = 10
	}
	cfResolveTimeoutStr := os.Getenv("DNS3LD_TESTRESOLVE_TIMEOUT")
	var cfResolveTimeout time.Duration
	if cfResolveTimeoutStr != "" {
		cfResolveTimeout, err = time.ParseDuration(cfResolveTimeoutStr)
		if err != nil {
			return err
		}
	} else {
		cfResolveTimeout = 5 * time.Second
	}
	cfDNSCheckServers := os.Getenv("DNS3LD_DNSCHECKSERVERS")
	if cfDNSCheckServers == "" {
		cfDNSCheckServers = "80.158.48.19:53,93.188.242.252:53"
	}
	re.ResolveMaxRetries = cfResolveMaxRetries
	re.ResolveTimeout = cfResolveTimeout
	re.DNSCheckServers = strings.Split(cfDNSCheckServers, ",")
	return nil
}
