package cloudflare

import (
	"log"

	"github.com/pkg/errors"

	"github.com/cloudflare/cloudflare-go"
)

// DNSAer provides functionality to work with DNS records.
type DNSAer interface {
	List(zoneID string, recordFetcher func(zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord,
		error)) ([]cloudflare.DNSRecord, error)
	UpdateIP(ip string, records []cloudflare.DNSRecord, recordUpdater func(zoneID, recordID string,
		rr cloudflare.DNSRecord) error) error
}

// DNSA is a representation of a DNS address record (aka "A" record).
type DNSA struct {
	recordFilter cloudflare.DNSRecord
	recordType   string
}

// NewDNSA constructs a DNSAer that is able to handle DNS a records.
func NewDNSA() DNSAer {
	filter := cloudflare.DNSRecord{Type: "A"}
	return DNSA{recordFilter: filter, recordType: "A"}
}

// List fetches and returns all DNS A records given vor the zone specified.
func (d DNSA) List(zoneID string, recordFetcher func(zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord,
	error)) ([]cloudflare.DNSRecord, error) {
	recs, err := recordFetcher(zoneID, d.recordFilter)
	if err != nil {
		return nil, err
	}
	log.Printf("fetched: %d DNS records of type: %q", len(recs), d.recordFilter.Type)
	return recs, nil
}

// UpdateIP performs an update of the IP of the given DNS records in case those are of type 'A'.
func (d DNSA) UpdateIP(ip string, records []cloudflare.DNSRecord, recordUpdater func(zoneID, recordID string,
	rr cloudflare.DNSRecord) error) error {
	for _, record := range records {
		if record.Type != d.recordType {
			continue
		}
		update := cloudflare.DNSRecord{
			Name:    record.Name,
			Content: ip,
			Proxied: record.Proxied,
			TTL:     record.TTL,
			Type:    record.Type,
		}
		err := recordUpdater(record.ZoneID, record.ID, update)
		if err != nil {
			return errors.Wrapf(err, "failure to update DNS record: %q", update.Name)
		}
		log.Printf("updated DNS record: %q to: %q", update.Name, update.Content)
	}
	return nil
}
