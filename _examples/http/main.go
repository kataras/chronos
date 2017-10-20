package main

import (
	"fmt"
	"time"

	"github.com/kataras/chronos/ext/http"
)

type geoInfo struct {
	IP          string  `json:"query,omitempty"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

func main() {

	// The 'ip-api.com' service has a rate limit of 150 requests per minute.
	var (
		max uint32 = 150
		per        = 1 * time.Minute
	)
	c := http.New(max, per)

	// This should take 2 minutes, 150 per minute.
	for i := 1; i <= 300; i++ {
		var (
			info      geoInfo
			ipToFetch = "165.227.36.244"
		)

		requestURL := fmt.Sprintf("http://ip-api.com/json/%s", ipToFetch)
		// it may pause for 3-4 seconds at 90-100, this is by the http service, don't be worry
		// it will continue until the 150 first, and after a minute from the last fetch
		// it will fetch the 150 more, this guards that you will not get banned by the service
		// even if you have to fetch thousands of millions of IPs at once,
		// chronos will schedule those.
		if err := c.ReadJSON(requestURL, &info); err != nil {
			// c.Do is like net/http.Client's Do but with the feature of
			// chronos' "rate" limiter; maximum operations per x time.
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("[%d] %#+v\n", i, info)
	}

}
