//nolint:unparam
package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/mock"
)

type APIMock struct {
	mock.Mock
}

func (a *APIMock) DNSRecord(zoneID string, record cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error) {
	args := a.Called(zoneID, record)
	return args.Get(0).([]cloudflare.DNSRecord), args.Error(1)
}

// UpdateDNSRecord updates a DNS address records with the given record.
func (a *APIMock) UpdateDNSRecord(zoneID, recordID string, record cloudflare.DNSRecord) error {
	args := a.Called(zoneID, recordID, record)
	return args.Error(0)
}
