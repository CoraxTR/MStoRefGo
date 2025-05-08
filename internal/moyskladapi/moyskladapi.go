package moyskladapi

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"mstorefgo/internal/config"
)

type MoySkladProcessor struct {
	Ratelimiter *Ratelimiter
	Config      *config.Moyskladapiconfig
}

func NewMoySkladProcessor(r *Ratelimiter, c *config.Moyskladapiconfig) *MoySkladProcessor {
	return &MoySkladProcessor{Ratelimiter: r, Config: c}
}

func (m *MoySkladProcessor) GetDeliverableOrders() {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	dayaftertomorrow := tomorrow.AddDate(0, 0, 1)
	tomorrowstart := ">=" + tomorrow.Format("2006-01-02") + " 00:00:00"
	dayaftertomorrowend := "<=" + dayaftertomorrow.Format("2006-01-02") + " 23:59:59"

	baseURL, err := url.Parse(m.Config.URLstart)
	if err != nil {
		panic(err)
	}
	baseURL.Path = path.Join(baseURL.Path, "entity/customerorder")

	filterValue := fmt.Sprintf(
		"deliveryPlannedMoment%s;deliveryPlannedMoment%s;state=%s", tomorrowstart, dayaftertomorrowend, m.Config.Statehref)

	q := baseURL.Query()
	q.Set("filter", filterValue)
	baseURL.RawQuery = q.Encode()
	fmt.Println(baseURL.String())

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+m.Config.APIKEY)
	req.Header.Set("Accept-Encoding", "gzip")
	log.Println("Waiting for RateLimiter")
	m.Ratelimiter.Wait()
	log.Println("Done waiting")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", body)

}
