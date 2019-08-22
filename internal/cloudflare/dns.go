package cloudflare

import (
	"log"

	"github.com/pkg/errors"

	"github.com/cloudflare/cloudflare-go"
)

type DNSAer interface {
	List(zoneID string, recordFetcher func(zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error)) ([]cloudflare.DNSRecord, error)
	UpdateIP(ip string, records []cloudflare.DNSRecord, recordUpdater func(zoneID, recordID string, rr cloudflare.DNSRecord) error) error
}

type DNSA struct {
	recordFilter cloudflare.DNSRecord
	recordType   string
}

func NewDNSA() DNSAer {
	filter := cloudflare.DNSRecord{Type: "A"}
	return DNSA{recordFilter: filter, recordType: "A"}
}

func (d DNSA) List(zoneID string, recordFetcher func(zoneID string,
	rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord,
	error)) ([]cloudflare.DNSRecord, error) {
	recs, err := recordFetcher(zoneID, d.recordFilter)
	if err != nil {
		return nil, err
	}
	log.Printf("retrieved %d records to update", len(recs))
	return recs, nil
}

func (d DNSA) UpdateIP(ip string, records []cloudflare.DNSRecord, recordUpdater func(zoneID,
	recordID string, rr cloudflare.DNSRecord) error) error {
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
			return errors.Wrapf(err, "failure to update record: %s", update.Name)
		}
		log.Printf("updated record: %s with ip: %s", update.Name, update.Content)
	}
	return nil
}
