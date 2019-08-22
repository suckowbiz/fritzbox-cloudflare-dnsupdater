package cloudflare

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestDNSARecords(t *testing.T) {
	zoneID := "4711"

	apiMock := *new(APIMock)
	apiMock.On("DNSRecord", zoneID, mock.MatchedBy(func(record cloudflare.DNSRecord) bool {
		return record.Type == "A"
	})).Return([]cloudflare.DNSRecord{}, nil)

	rec, err := NewDNSA().List(zoneID, apiMock.DNSRecord)

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(rec)

	apiMock.AssertExpectations(t)
}

func TestUpdateDNSARecordsIP(t *testing.T) {
	ip := "127.0.0.1"
	zoneID := "4711"
	givenRecord := cloudflare.DNSRecord{
		ID:      "42",
		Name:    "localhost",
		Proxied: true,
		TTL:     1,
		Type:    "A",
		ZoneID:  zoneID,
	}
	updateRecord := cloudflare.DNSRecord{
		Content: ip,
		Name:    "localhost",
		Proxied: true,
		TTL:     1,
		Type:    "A",
	}

	apiMock := *new(APIMock)
	apiMock.On("UpdateDNSRecord", zoneID, givenRecord.ID,
		mock.MatchedBy(func(record cloudflare.DNSRecord) bool {
			return record.Type == updateRecord.Type &&
				record.Name == updateRecord.Name &&
				record.Content == ip &&
				record.TTL == updateRecord.TTL &&
				record.Proxied == updateRecord.Proxied
		})).Return(nil)

	err := NewDNSA().UpdateIP(ip, []cloudflare.DNSRecord{givenRecord, {Type: "AAA"}}, apiMock.UpdateDNSRecord)
	assert.NoError(t, err)

	apiMock.AssertExpectations(t)
}
