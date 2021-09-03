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
	"net/http"
	"strings"
	"swan"
)

// HandlerHTML returns HTML that does not require a model for the template.
func HandlerHTML(d *Domain, w http.ResponseWriter, r *http.Request) {

	// If the request is for the OWID creator then direct to that handler.
	if r.URL.Path == "/create-owid" {
		handlerCreateOWID(d, w, r)
		return
	}

	// If the request is being proxied to SWAN then pass to SWAN access node.
	if strings.HasPrefix(r.URL.Path, "/swan-proxy") {
		handlerSWANProxy(d, w, r)
		return
	}

	// Get the template for the URL path.
	t := d.LookupHTML(r.URL.Path)
	if t == nil {
		http.NotFound(w, r)
		return
	}

	// Execute the template without a model.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	err := t.Execute(g, nil)
	if err != nil {
		ReturnServerError(d.Config, w, &swan.Error{err, nil})
	}
}
