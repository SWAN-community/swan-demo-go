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
	"html/template"
	"net/http"
	"owid"
	"swan"
)

// infoModel data needed for the advert information interface.
type infoModel struct {
	OWIDs      map[*owid.OWID]interface{}
	Bid        *swan.Bid
	ID         *swan.ID
	Root       *owid.OWID
	ReturnURL  template.HTML
	AccessNode string
}

func (m *infoModel) findID() (*owid.OWID, *swan.ID) {
	for k, v := range m.OWIDs {
		if o, ok := v.(*swan.ID); ok {
			return k, o
		}
	}
	return nil, nil
}

func (m *infoModel) findBid() *swan.Bid {
	for _, v := range m.OWIDs {
		if b, ok := v.(*swan.Bid); ok {
			return b
		}
	}
	return nil
}

func handlerInfo(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	// Get the SWAN OWIDs from the form parameters.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	var m infoModel
	m.OWIDs = make(map[*owid.OWID]interface{})
	for k, vs := range r.Form {
		if k == "owid" {
			for _, v := range vs {
				o, err := owid.FromBase64(v)
				if err != nil {
					common.ReturnServerError(d.Config, w, err)
					return
				}
				m.OWIDs[o], err = swan.FromOWID(o)
				if err != nil {
					common.ReturnServerError(d.Config, w, err)
					return
				}
			}
		}
	}

	// Set the common fields.
	m.Bid = m.findBid()
	m.Root, m.ID = m.findID()
	f, err := common.GetReturnURL(r)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	m.ReturnURL = template.HTML(f.String())
	m.AccessNode = r.Form.Get("accessNode")

	// Display the template form.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	err = d.LookupHTML("info.html").Execute(g, m)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
}
