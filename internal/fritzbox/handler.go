package fritzbox

import (
	"errors"
	"fritzbox-cloudflare-dnsupdater/internal/cloudflare"
	"log"
	"net/http"
	"net/url"

	api "github.com/cloudflare/cloudflare-go"
)

type DNSer interface {
	DNSARecords(zoneID string, recordFetcher func(zoneID string, rr api.DNSRecord) ([]api.DNSRecord,
		error)) ([]api.DNSRecord, error)
}

func NewUpdateHandler(apiProvider func(token string, opts ...api.Option) (*api.API, error),
	dnsAer cloudflare.DNSAer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		zoneID, token, ip, err := validateUpdateQuery(request.URL.RawQuery, url.ParseQuery)
		if err != nil {
			log.Print("Error: ", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("creating api with token: %s", token)
		api, err := apiProvider(token)
		if err != nil {
			log.Print(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
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

func validateUpdateQuery(rawQuery string, parser func(query string) (url.Values, error)) (zoneID string, token string,
	ip string, err error) {
	values, err := parser(rawQuery)
	if err != nil {
		return "", "", "", err
	}
	if len(values) == 0 {
		return "", "", "", errors.New("url query must not be empty")
	}

	token = values.Get("token")
	if token == "" {
		return "", "", "", errors.New("'token' must not be absent")
	}
	zoneID = values.Get("zone_identifier")
	if zoneID == "" {
		return "", "", "", errors.New("'zone_identifier' must not be absent")
	}
	ip = values.Get("ip")
	if ip == "" {
		return "", "", "", errors.New("'ip' must not be absent")
	}
	return
}
