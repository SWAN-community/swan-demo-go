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

package publisher

import (
	"common"
	"compress/gzip"
	"fmt"
	"net/http"
	"swan"
)

// HandlerAdvert for the request for adverts for the publisher web pages.
func HandlerAdvert(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	// Create the model for publishers.
	var m Model
	m.Domain = d
	m.Request = r

	// Get the form parameters which will include the encrypted data.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// See if there is also SWAN data in the request. If so then use it for this
	// advert and set cookies to store it in the response.
	if r.Form.Get("encrypted") != "" {
		var e *swan.Error
		m.swanData, e = newSWANData(d, r.Form.Get("encrypted"))
		if e != nil {
			common.ReturnProxyError(d.Config, w, e)
			return
		}
		setCookies(r, w, m.swanData)
	} else {
		m.swanData, err = newSWANDataFromCookies(r)
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
		}
		if m.swanData == nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				fmt.Errorf("SWAN data cookies missing for request"),
				http.StatusBadRequest)
			return
		}
	}

	// Use the new advert HTML to request the advert.
	t, err := m.NewAdvertHTML(r.Form.Get("placement"))
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Respond with the HTML.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, err = g.Write([]byte(t))
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
}
