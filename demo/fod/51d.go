/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package fod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Device is the 51Degrees.com device item returned from calls to
// cloud.51degrees.com.
type Device struct {
	IsCrawler bool `json:"iscrawler"`
}

// FOD all the information returned from the cloud.51degrees.com service.
type FOD struct {
	Device *Device `json:"device"`
}

// GetCrawlerFrom51Degrees used the 51Degrees.com device detection service to
// determine if the request is from a crawler. Needs the FOD_RESOURCE_KEY
// environment variable configured with a valid resource key from
// https://configure.51degrees.com/vXyRZz8B.
func GetCrawlerFrom51Degrees(r *http.Request) (bool, error) {

	key := os.Getenv("FOD_RESOURCE_KEY")
	if key == "" {
		// 51Degrees device detection is not enabled so return false as the
		// default.
		return false, nil
	}

	// Add all the HTTP headers from the request as query string parameters.
	f := url.Values{}
	for n, v := range r.Header {
		if n != "cookie" {
			f.Add(n, v[0])
		}
	}
	f.Add("Host", r.Host)

	// Get the response from the cloud service.
	resp, err := http.PostForm(
		"https://cloud.51degrees.com/api/v4/"+key+".json",
		f)
	if err != nil {
		return false, err
	}

	// There are limited subscriptions that are throttled or have fixed
	// entitlements. There will return a 429 error if usage is exceed. In these
	// situations treat the request as non crawler rather than display an error.
	// Different paid for subscriptions without this restriction can be trialled
	// at http://51degrees.com/pricing
	if resp.StatusCode == http.StatusTooManyRequests {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("status code '%d' returned", resp.StatusCode)
	}

	defer resp.Body.Close()
	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var fod FOD
	err = json.Unmarshal(j, &fod)
	if err != nil {
		return false, err
	}

	return fod.Device.IsCrawler, nil
}
