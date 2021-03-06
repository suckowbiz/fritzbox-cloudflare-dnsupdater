package fritzbox

import (
	"errors"
	"fritzbox-cloudflare-dnsupdater/internal/cloudflare"
	"log"
	"net/http"
	"net/url"

	api "github.com/cloudflare/cloudflare-go"
)

// DNSer combines functionality to handle DNS type "A" records.
type DNSer interface {
	DNSARecords(zoneID string, recordFetcher func(zoneID string, rr api.DNSRecord) ([]api.DNSRecord,
		error)) ([]api.DNSRecord, error)
}

// NewUpdateHandler creates a new instance of a http handler to answer DNS type "A" update requests.
func NewUpdateHandler(apiProvider func(token string, opts ...api.Option) (*api.API, error),
	dnsAer cloudflare.DNSAer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		zoneIDs, token, ip, err := validateUpdateQuery(request.URL.RawQuery, url.ParseQuery)
		if err != nil {
			log.Print("Error: ", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		api, err := apiProvider(token)
		if err != nil {
			log.Print(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, zoneID := range zoneIDs {
			log.Printf("Updating zone: %q", zoneID)
			rr, err := dnsAer.List(zoneID, api.DNSRecords)
			if err != nil {
				log.Print(err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			err = dnsAer.UpdateIP(ip, rr, api.UpdateDNSRecord)
			if err != nil {
				log.Print(err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func validateUpdateQuery(rawQuery string, parser func(query string) (url.Values, error)) (zoneID []string, token string,
	ip string, err error) {
	values, err := parser(rawQuery)
	if err != nil {
		return []string{}, "", "", err
	}
	if len(values) == 0 {
		return []string{}, "", "", errors.New("url query must not be empty")
	}

	token = values.Get("token")
	if token == "" {
		return []string{}, "", "", errors.New("'token' must not be absent")
	}
	zoneID, ok := values["zone_id"]
	if !ok {
		return []string{}, "", "", errors.New("'zone_id' must not be absent")
	}
	ip = values.Get("ip")
	if ip == "" {
		return []string{}, "", "", errors.New("'ip' must not be absent")
	}
	return
}
