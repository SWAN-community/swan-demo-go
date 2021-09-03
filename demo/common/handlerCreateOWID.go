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
	"owid"
)

// handlerCreateOWID takes an input payload and returns a new OWID.
func handlerCreateOWID(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ReturnServerError(d.Config, w, err)
		return
	}

	// Get the OWID creator for this domain.
	oc, err := d.GetOWIDCreator()
	if err != nil {
		ReturnServerError(d.Config, w, err)
	}

	ot, err := decodeOthers(r.Form["others"])
	if err != nil {
		ReturnServerError(d.Config, w, err)
	}

	// Create and sign the OWID for this domain.
	o, err := oc.CreateOWIDandSign([]byte(r.Form.Get("payload")), ot...)
	if err != nil {
		ReturnServerError(d.Config, w, err)
	}

	// Return the OWID as a base 64 string.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, err = g.Write([]byte(o.AsString()))
	if err != nil {
		ReturnServerError(d.Config, w, err)
		return
	}
}

// decodeOthers decodes an array of other OWID Base 64 strings and returns an
// array of OWID instances.
func decodeOthers(v []string) ([]*owid.OWID, error) {
	var os []*owid.OWID

	for _, o := range v {
		o, err := owid.FromBase64(o)
		if err != nil {
			return nil, err
		}

		os = append(os, o)
	}

	return os, nil
}
