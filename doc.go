// Copyright (c) 2017 Gerasimos Maropoulos. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of chronos nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package chronos provides an easy way to limit X operations per Y time in accuracy of nanoseconds.
Chronos is a crucial tool when calling external, remote or even internal APIs with rate limits.

Installation

The only requirement is the Go Programming Language.

    $ go get -u github.com/kataras/chronos

Chronos has no any other dependencies except the standard library.

Chronos is very simple to use tool, you just declare the maximum operations(max) and in every Y time
that those operations are allowed to be executed(per). Chronos has the accuracy of nanoseconds.

    c := chronos.New(max uint32, per time.Duration)
    <- c.Acquire()
    Do()

The Acquire is blocking when needed, no any further actions neeeded by the end-developer.


Using the ext/http subpackage

This package is just a helper which extends the usability of the `chronos.C.Acquire` functionality.
It allows easy and fast calls to external APIs with limits.

Example Code:


    package chronos

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
*/
package chronos
