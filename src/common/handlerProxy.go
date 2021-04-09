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

package common

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"swift"
)

// handlerProxy takes an incoming request, adds the access key to the parameters
// and then passes on the result to the SWAN access node.
func handlerSWANProxy(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		ReturnServerError(d.Config, w, err)
		return
	}

	// Get the converted URL for the access node.
	var u url.URL
	u.Scheme = d.Config.Scheme
	u.Host = d.SWANAccessNode
	u.Path = strings.Replace(r.URL.Path, "swan-proxy", "swan", 1)

	// Add the access key to the incoming parameters.
	r.Form.Set("accessKey", d.SWANAccessKey)

	// Add any additional parameters that might be important from HTTP headers.
	swift.SetHomeNodeHeaders(r, &r.Form)

	// Post the data to the SWAN endpoint.
	res, err := http.PostForm(u.String(), r.Form)
	if err != nil {
		ReturnServerError(d.Config, w, err)
	}

	// Get the body.
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ReturnServerError(d.Config, w, err)
	}

	// Return the OWID as a base 64 string.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", "no-cache")
	_, err = g.Write(b)
	if err != nil {
		ReturnServerError(d.Config, w, err)
		return
	}
}
