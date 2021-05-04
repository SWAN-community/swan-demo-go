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

package marketer

import (
	"common"
	"compress/gzip"
	"encoding/base64"
	"net/http"
	"owid"
)

// Handler for the marketer features.
func Handler(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	// Get the template for the URL path.
	t := d.LookupHTML(r.URL.Path)
	if t == nil {
		http.NotFound(w, r)
		return
	}

	// Get the SWAN ID that relates to the advert.
	o, err := getID(r)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Create the model for use with the page template.
	var m MarketerModel
	m.Domain = d
	m.Request = r
	m.idNode = o

	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	err = t.Execute(g, &m)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
	}
}

func getID(r *http.Request) (*owid.Node, error) {

	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	// If the bid data does not exist return warning.
	if r.Form.Get("transaction") == "" {
		return nil, nil
	}

	// Get the transaction from the form data.
	d, err := base64.RawStdEncoding.DecodeString(
		r.Form.Get("transaction"))
	if err != nil {
		return nil, err
	}

	return owid.NodeFromJSON(d)
}
