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

package cmp

import (
	"common"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/url"
)

// handlerStop matches the path /stop and redirects the response to the
// SWAN stop preferences URL where the parameters include the host that should
// be stopped from displaying adverts on the browser.
func handlerStop(d *common.Domain, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	if r.Form.Get("host") == "" {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("Host to be stopped must be provided"),
			http.StatusBadRequest)
		return
	}

	// Configure the update operation from this demo domain's configuration with
	// the return URL and host.
	returnUrl, err := url.Parse(r.Form.Get("returnUrl"))
	if err != nil {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			err,
			http.StatusBadRequest)
	}
	s := d.SWAN().NewStop(r, returnUrl.String(), r.Form.Get("host"))

	// Use the access node from the form as this will be used by the publisher
	// to decrypt the result.
	if r.Form.Get("accessNode") != "" {
		s.AccessNode = r.Form.Get("accessNode")
	}

	// Demonstrate the CMP can change the SWAN operation messages.
	if r.Form.Get("message") == "" {
		s.Message = fmt.Sprintf(
			"Bye, bye %s. Thanks for letting us know.",
			s.Host)
	}

	// Get the URL to process the stop data.
	u, se := s.GetURL()
	if se != nil {
		common.ReturnProxyError(d.Config, w, se)
		return
	}

	// Return the URL as a text string.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, err = g.Write([]byte(u))
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
}
