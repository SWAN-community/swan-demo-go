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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"owid"
	"reflect"
	"swift"

	uuid "github.com/satori/go.uuid"
)

type dialogModel struct {
	url.Values
	update bool // True if the update should be performed
}

// Title for the SWAN storage operation.
func (m *dialogModel) Title() string { return m.Get("title") }

// SWID as a base64 OWID.
func (m *dialogModel) SWIDAsOWID() string { return m.Get("swid") }

// Email as a string.
func (m *dialogModel) Email() string { return m.Get("email") }

// Pref as a string.
func (m *dialogModel) Pref() string { return m.Get("pref") }

// BackgroundColor for the SWAN storage operation.
func (m *dialogModel) BackgroundColor() string {
	return m.Get("backgroundColor")
}

// PublisherHost the domain from the returnUrl.
func (m *dialogModel) PublisherHost() string {
	u, _ := url.Parse(m.Get("returnUrl"))
	if u != nil {
		return u.Host
	}
	return ""
}

// SWIDAsString returns the SWID as a readable string without the OWID data.
func (m *dialogModel) SWIDAsString() (string, error) {
	o, err := owid.FromBase64(m.Get("swid"))
	if err != nil {
		return "", err
	}
	u, err := uuid.FromBytes(o.Payload)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func handlerDialog(d *common.Domain, w http.ResponseWriter, r *http.Request) {
	var m dialogModel
	m.Values = make(url.Values)

	// Parse the form variables.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Get the SWAN data from the request path.
	s := common.GetSWANDataFromRequest(r)
	if s == "" {
		redirectToSWANDialog(d, w, r)
		return
	}

	// Call the SWAN access node for the CMP to turn the data provided in the
	// URL into usable data for the dialog.
	e := decryptAndDecode(d, s, &m)
	if e != nil {

		// If the data can't be decrypted rather than another type of error
		// then redirect via SWAN to the dialog.
		if e.StatusCode() >= 400 && e.StatusCode() < 500 {
			redirectToSWANDialog(d, w, r)
			return
		}
		common.ReturnStatusCodeError(d.Config, w, e.Err, http.StatusBadRequest)
		return
	}

	// If this is a close request then don't update the values and just return
	// to the return URL.
	if r.Form.Get("close") != "" {
		http.Redirect(w, r, m.Get("returnUrl"), 303)
		return
	}

	// If the method is POST then update the model with the data from the form.
	if r.Method == "POST" {
		se := dialogUpdateModel(d, r, &m)
		if se != nil {
			common.ReturnProxyError(d.Config, w, se)
			return
		}
	}

	// If the redirect URL has been set then redirect, otherwise display the
	// HTML template.
	if m.update == true {

		// The user has request that the data be updated in the SWAN network.
		// Set the redirection URL for the operation to store the data. The web
		// browser will then be redirected to that URL, the data saved and the
		// return URL for the publisher returned to.
		u, err := getRedirectUpdateURL(d, r, m.Values)
		if err != nil {
			common.ReturnProxyError(d.Config, w, err)
		}
		http.Redirect(w, r, u, 303)

	} else {

		// The dialog needs to be displayed. Use the cmp.html template for the
		// user interface.
		g := gzip.NewWriter(w)
		defer g.Close()
		w.Header().Set("Content-Encoding", "gzip")
		err := d.LookupHTML("cmp.html").Execute(g, &m)
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}
	}
}

func dialogUpdateModel(
	d *common.Domain,
	r *http.Request,
	m *dialogModel) *common.SWANError {

	// Copy the field values from the form.
	m.Values.Set("swid", r.Form.Get("swid"))
	m.Values.Set("email", r.Form.Get("email"))
	m.Values.Set("pref", r.Form.Get("pref"))

	// Check to see if the post is as a result of the SWID reset.
	if r.Form.Get("reset-swid") != "" {

		// Replace the SWID with a new random value.
		return setNewSWID(d, m)
	}

	// Check to see if the post is as a result for all data.
	if r.Form.Get("reset-all") != "" {

		// Replace the data.
		m.Set("email", "")
		m.Set("pref", "")
		return setNewSWID(d, m)
	}

	// The data should be updated in the SWAN network.
	m.update = true

	return nil
}

func setNewSWID(d *common.Domain, m *dialogModel) *common.SWANError {
	c, se := createSWID(d)
	if se != nil {
		return se
	}
	o, err := owid.FromByteArray(c)
	if err != nil {
		return &common.SWANError{Err: err}
	}
	m.Set("swid", o.AsString())
	return nil
}

func getRedirectUpdateURL(
	d *common.Domain,
	r *http.Request,
	m url.Values) (string, *common.SWANError) {
	c, err := d.GetOWIDCreator()
	if err != nil {
		return "", &common.SWANError{Err: err}
	}
	b, se := d.CallSWANStorageURL(r, "update", func(q url.Values) error {
		var err error

		// Loop through all the key value pairs in the model values. If the key
		// relates to SWAN data then turn the value into an OWID with this UIP
		// as the signatory.
		for k, v := range m {
			switch k {
			case "pref":
				a := v[0]
				if a == "" {
					a = "off"
				}
				err = setSWANData(c, &q, k, []byte(a))
				break
			case "email":
				err = setSWANData(c, &q, k, []byte(v[0]))
				break
			default:
				for _, i := range v {
					q.Add(k, i)
				}
				break
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	if se != nil {
		return "", se
	}
	return string(b), nil
}

func setSWANData(c *owid.Creator, q *url.Values, k string, v []byte) error {
	o, err := c.CreateOWIDandSign(v)
	if err != nil {
		return err
	}
	q.Set(k, o.AsString())
	return nil
}

func createSWID(d *common.Domain) ([]byte, *common.SWANError) {
	b, e := d.CallSWANURL("create-swid", nil)
	if e != nil {
		return nil, e
	}
	return b, nil
}

func decryptAndDecode(
	d *common.Domain,
	v string,
	m *dialogModel) *common.SWANError {
	b, e := d.CallSWANURL("decrypt-raw", func(q url.Values) error {
		q.Set("encrypted", v)
		return nil
	})
	if e != nil {
		return e
	}
	r := make(map[string]interface{})
	err := json.Unmarshal(b, &r)
	if err != nil {
		return &common.SWANError{Err: err}
	}
	for k, v := range r {
		switch reflect.TypeOf(v) {
		case reflect.TypeOf([]interface{}(nil)):
			for i, a := range v.([]interface{}) {
				switch i {
				case 0:
					m.Set("returnUrl", a.(string))
					break
				case 1:
					m.Set("accessNode", a.(string))
					break
				case 2:
					m.Set("displayUserInterface", a.(string))
					break
				case 3:
					m.Set("postMessageOnComplete", a.(string))
					break
				}
			}
			break
		case reflect.TypeOf(""):
			m.Set(k, v.(string))
			break
		}
	}
	return nil
}

func redirectToSWANDialog(
	d *common.Domain,
	w http.ResponseWriter,
	r *http.Request) {
	f, err := common.GetReturnURL(r)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	a := r.Form.Get("accessNode")
	if a == "" {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("SWAN accessNode parameter required for CMP operation"),
			http.StatusBadRequest)
		return
	}
	u, e := d.CreateSWANURL(
		r,
		// Use this CMP page as the return URL for fetching the SWAN data.
		common.GetCurrentPage(d.Config, r).String(),
		"fetch",
		func(q url.Values) {

			// Use the return URL provided in the request to this URL as the
			// final return URL after the update has occurred. Store in the
			// state for use when the CMP dialogue updates.
			q.Add("state", f.String())

			// Also also add the access node to the state store.
			q.Add("state", a)

			// Add the flags.
			q.Add("state", r.Form.Get("displayUserInterface"))
			q.Add("state", r.Form.Get("postMessageOnComplete"))
		})
	if e != nil {
		common.ReturnProxyError(d.Config, w, e)
		return
	}
	http.Redirect(w, r, u, 303)
}

// All the values from the SWIFT storage operation are OWIDs. Returns the OWID
// for the key provided.
func getOWID(r *swift.Results, k string) (*owid.OWID, error) {
	v := r.Get(k)
	if v != nil {
		return nil, fmt.Errorf("Key '%s' does not exist", k)
	}
	if len(v.Value()) != 1 {
		return nil, fmt.Errorf(
			"Key '%s' contains '%d' values",
			k,
			len(v.Values()))
	}
	return owid.FromByteArray(v.Values()[0])
}
