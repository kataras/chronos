package http

import (
	"fmt"
	"time"
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

// This service has a rate limit of 150 requests per minute.
func geoServiceFromIP(ip string) string {
	return fmt.Sprintf("http://ip-api.com/json/%s", ip)
}

func ExampleClient_ReadJSON() {
	// c :=  New(3, 4*time.Second)
	//
	// But we can't wait a minute to finish the example
	// and we don't have 150 addresses available, so:
	var (
		max uint32 = 3
		per        = 4 * time.Second
	)

	c := New(max, per)

	ipToFind := "165.227.36.244"
	time.AfterFunc(time.Duration(per.Nanoseconds()-200*time.Millisecond.Nanoseconds()), func() {
		fmt.Println("waiting for fetch-4 to be fired after 4 seconds the last one (or immediately if the limit time passed)...")
	})

	var lastFetched time.Time
	var i uint32
	for i = 1; i <= max+1; i++ {
		lastFetched = time.Now()
		var info geoInfo
		reqURL := geoServiceFromIP(ipToFind)
		if err := c.ReadJSON(reqURL, &info); err != nil {
			panic(err)
		}
		fmt.Printf("fetch-%d\n", i)
		fmt.Printf("%#+v\n", info)

		if i == max+1 {
			since := time.Since(lastFetched)
			// the last one will be executed after 4 seconds from the last executed.
			fmt.Printf("finish-after-%.0f-seconds\n", since.Seconds())
		}
	}

	// Output:
	// fetch-1
	// http.geoInfo{IP:"165.227.36.244", Country:"Canada", CountryCode:"CA", Region:"ON", RegionName:"Ontario", City:"Toronto", Zip:"M5A", Lat:43.6555, Lon:-79.3626, Timezone:"America/Toronto", ISP:"Digital Ocean", Org:"Digital Ocean", As:"AS14061 Digital Ocean, Inc."}
	// fetch-2
	// http.geoInfo{IP:"165.227.36.244", Country:"Canada", CountryCode:"CA", Region:"ON", RegionName:"Ontario", City:"Toronto", Zip:"M5A", Lat:43.6555, Lon:-79.3626, Timezone:"America/Toronto", ISP:"Digital Ocean", Org:"Digital Ocean", As:"AS14061 Digital Ocean, Inc."}
	// fetch-3
	// http.geoInfo{IP:"165.227.36.244", Country:"Canada", CountryCode:"CA", Region:"ON", RegionName:"Ontario", City:"Toronto", Zip:"M5A", Lat:43.6555, Lon:-79.3626, Timezone:"America/Toronto", ISP:"Digital Ocean", Org:"Digital Ocean", As:"AS14061 Digital Ocean, Inc."}
	// waiting for fetch-4 to be fired after 4 seconds the last one (or immediately if the limit time passed)...
	// fetch-4
	// http.geoInfo{IP:"165.227.36.244", Country:"Canada", CountryCode:"CA", Region:"ON", RegionName:"Ontario", City:"Toronto", Zip:"M5A", Lat:43.6555, Lon:-79.3626, Timezone:"America/Toronto", ISP:"Digital Ocean", Org:"Digital Ocean", As:"AS14061 Digital Ocean, Inc."}
	// finish-after-4-seconds
}
