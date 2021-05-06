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
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"owid"
	"reflect"
	"strconv"
	"strings"
	"swan"

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

// Salt as a string
func (m *dialogModel) Salt() string { return m.Get("salt") }

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

// HiddenFields turns the parameters from the storage operation into hidden
// fields so they are available when the form is posted.
func (m *dialogModel) HiddenFields() template.HTML {
	b := strings.Builder{}
	for k, v := range m.Values {
		if k != "salt" && k != "swid" && k != "email" && k != "pref" {
			b.WriteString(fmt.Sprintf(
				"<input type=\"hidden\" id=\"%s\" name=\"%s\" value=\"%s\"/>",
				k, k, v[0]))
		}
	}
	return template.HTML(b.String())
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

	// If all the form variables are present in the query string then these can
	// be used with the form and the update initiated. A new SWID will be used.
	// TODO: This values will need to be encrypted as a part of a long lived
	// SWIFT storage transaction before production use.
	if r.Method == "GET" &&
		r.Form.Get("email") != "" &&
		r.Form.Get("salt") != "" &&
		r.Form.Get("pref") != "" &&
		r.Form.Get("accessNode") != "" &&
		r.Form.Get("returnUrl") != "" {

		// Copy the key values.
		m.Values.Set("email", r.Form.Get("email"))
		m.Values.Set("salt", r.Form.Get("salt"))
		m.Values.Set("pref", r.Form.Get("pref"))
		m.Values.Set("accessNode", r.Form.Get("accessNode"))
		m.Values.Set("returnUrl", r.Form.Get("returnUrl"))

		// Get a new SWID.
		se := setNewSWID(d, &m)
		if se != nil {
			common.ReturnProxyError(d.Config, w, se)
			return
		}

		// Automatically trigger the update with the values provided.
		m.update = true

	} else {

		// Not parameters were provided so get the SWAN data from the request
		// path.
		s := common.GetSWANDataFromRequest(r)
		if s == "" {
			redirectToSWANDialog(d, w, r)
			return
		}

		// Call the SWAN access node for the CMP to turn the data provided in
		// the URL into usable data for the dialog.
		e := decryptAndDecode(d, s, &m)
		if e != nil {

			// If the data can't be decrypted rather than another type of
			// error then redirect via SWAN to the dialog.
			if e.StatusCode() >= 400 && e.StatusCode() < 500 {
				redirectToSWANDialog(d, w, r)
				return
			}
			common.ReturnStatusCodeError(
				d.Config,
				w,
				e.Err,
				http.StatusBadRequest)
			return
		}
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
		if err != nil {
			common.ReturnProxyError(d.Config, w, err)
		}

		// Get the CMP URL for the email.
		eu := getCMPURL(d, r, &m)

		// Send the email if the SMTP server is setup.
		e := sendReminderEmail(d, m.Values, eu)
		if e != nil {
			fmt.Println(err)
		}

		// Redirect the response to the return URL.
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

func sendReminderEmail(d *common.Domain, m url.Values, u string) error {
	smtp := common.NewSMTP()

	e := m.Get("email")

	if e == "" {
		return nil
	}

	s := "SWAN Demo: Email Reminder"
	t := d.LookupHTML("email-template.html")

	b := m.Get("salt")

	if b == "" {
		return nil
	}

	var a []byte
	a, err := base64.RawStdEncoding.DecodeString(b)
	if err != nil {
		return err
	}

	if len(a) != 2 {
		return nil
	}

	s1, s2 := a[0]>>4, a[0]&0xF
	s3, s4 := a[1]>>4, a[1]&0xF

	var arr = []byte{s1, s2, s3, s4}

	td := EmailTemplate{
		Salt:           arr,
		PreferencesUrl: u,
	}

	err = smtp.Send(e, s, t, td)
	if err != nil {
		return err
	}

	return nil

}

func dialogUpdateModel(
	d *common.Domain,
	r *http.Request,
	m *dialogModel) *swan.Error {

	// Copy the field values from the form.
	m.Values.Set("swid", r.Form.Get("swid"))
	m.Values.Set("email", r.Form.Get("email"))
	m.Values.Set("salt", r.Form.Get("salt"))
	m.Values.Set("pref", r.Form.Get("pref"))

	// Check to see if the post is as a result of the SWID reset.
	if r.Form.Get("reset-swid") != "" {

		// Replace the SWID with a new random value.
		return setNewSWID(d, m)
	}

	// Check to see if the email and salt are being reset.
	if r.Form.Get("reset-email-salt") != "" {
		m.Set("email", "")
		m.Set("salt", "")
		return nil
	}

	// Check to see if the post is as a result for all data.
	if r.Form.Get("reset-all") != "" {

		// Replace the data.
		m.Set("email", "")
		m.Set("salt", "")
		m.Set("pref", "")
		return setNewSWID(d, m)
	}

	// The data should be updated in the SWAN network.
	m.update = true

	return nil
}

func setNewSWID(d *common.Domain, m *dialogModel) *swan.Error {
	o, err := d.SWAN().CreateSWID()
	if err != nil {
		return err
	}
	m.Set("swid", o.AsString())
	return nil
}

func getRedirectUpdateURL(
	d *common.Domain,
	r *http.Request,
	m url.Values) (string, *swan.Error) {

	// Get the OWID creator which is needed to sign the data just captured.
	c, err := d.GetOWIDCreator()
	if err != nil {
		return "", &swan.Error{Err: err}
	}

	// Configure the update operation from this demo domain's configuration.
	returnUrl, err := url.Parse(r.Form.Get("returnUrl"))
	if err != nil {
		return "", &swan.Error{Err: err}
	}
	u := d.SWAN().NewUpdate(r, returnUrl)

	// Use the form to get any information from the initial storage operation
	// to configure the update storage operation.
	if r.Form.Get("accessNode") != "" {
		u.AccessNode = r.Form.Get("accessNode")
	}
	if r.Form.Get("backgroundColor") != "" {
		u.BackgroundColor = r.Form.Get("backgroundColor")
	}
	if r.Form.Get("displayUserInterface") != "" {
		u.DisplayUserInterface = r.Form.Get("displayUserInterface") == "true"
	}
	if r.Form.Get("javaScript") != "" {
		u.JavaScript = r.Form.Get("javaScript") == "true"
	}
	if r.Form.Get("message") != "" {
		u.Message = r.Form.Get("message")
	}
	if r.Form.Get("messageColor") != "" {
		u.MessageColor = r.Form.Get("messageColor")
	}
	if r.Form.Get("postMessageOnComplete") != "" {
		u.PostMessageOnComplete = r.Form.Get("postMessageOnComplete") == "true"
	}
	if r.Form.Get("progressColor") != "" {
		u.ProgressColor = r.Form.Get("progressColor")
	}
	if r.Form.Get("title") != "" {
		u.Title = r.Form.Get("title")
	}
	if r.Form.Get("useHomeNode") != "" {
		u.UseHomeNode = r.Form.Get("useHomeNode") == "true"
	}
	u.Pref = m.Get("pref") == "on"
	u.Email = m.Get("email")
	u.Salt = []byte(m.Get("salt"))
	u.SWID = m.Get("swid")
	return u.GetURL(c)
}

func decodeOWID(
	k string,
	r *http.Request,
	m *dialogModel,
	payloadAsString func(*owid.OWID) string) error {
	o, err := owid.FromBase64(r.Form.Get(k))
	if err != nil {
		return err
	}
	m.Set(k, payloadAsString(o))
	return nil
}

func decode(
	d *common.Domain,
	r *http.Request,
	m *dialogModel) *swan.Error {

	err := decodeOWID("email", r, m, func(o *owid.OWID) string {
		return o.PayloadAsString()
	})
	if err != nil {
		return &swan.Error{Err: err}
	}

	err = decodeOWID("salt", r, m, func(o *owid.OWID) string {
		return o.PayloadAsString()
	})
	if err != nil {
		return &swan.Error{Err: err}
	}

	err = decodeOWID("pref", r, m,
		func(o *owid.OWID) string {
			return o.PayloadAsString()
		})
	if err != nil {
		return &swan.Error{Err: err}
	}

	setNewSWID(d, m)

	m.Set("returnUrl", r.Form.Get("returnUrl"))
	m.Set("accessNode", r.Form.Get("accessNode"))
	m.update = true

	return nil
}

func decryptAndDecode(
	d *common.Domain,
	v string,
	m *dialogModel) *swan.Error {
	r, err := d.SWAN().DecryptRaw(v)
	if err != nil {
		return err
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

	// Create the fetch function returning to this URL.
	f := d.SWAN().NewFetch(r, common.GetCleanURL(d.Config, r))

	// User Interface Provider fetch operations only need to consider
	// one node if the caller will have already recently accessed SWAN.
	// This will be true for callers that have not used third party
	// cookies to fetch data from SWAN prior to calling this API. if the
	// request has a node count then use that, otherwise use 1 to get
	// the data from the home node.
	if r.Form.Get("nodeCount") != "" {
		i, err := strconv.ParseInt(r.Form.Get("nodeCount"), 10, 32)
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
			return
		}
		f.NodeCount = int(i)
	} else {
		f.NodeCount = 1
	}

	f.State = make([]string, 4)

	// Use the return URL provided in the request to this URL as the
	// final return URL after the update has occurred. Store in the
	// state for use when the CMP dialogue updates.
	returnUrl, err := common.GetReturnURL(r)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	f.State[0] = returnUrl.String()

	// Also also add the access node to the state store.
	f.State[1] = r.Form.Get("accessNode")
	if f.State[1] == "" {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("SWAN accessNode parameter required for CMP operation"),
			http.StatusBadRequest)
		return
	}

	// Add the flags.
	f.State[2] = r.Form.Get("displayUserInterface")
	f.State[3] = r.Form.Get("postMessageOnComplete")

	// Get the URL.
	u, se := f.GetURL()
	if se != nil {
		common.ReturnProxyError(d.Config, w, se)
		return
	}
	http.Redirect(w, r, u, 303)
}

// Returns the CMP preferences URL.
func getCMPURL(d *common.Domain, r *http.Request, m *dialogModel) string {
	var u url.URL
	u.Scheme = d.Config.Scheme
	u.Host = d.Host
	u.Path = "/preferences/"
	q := u.Query()
	q.Set("returnUrl", r.Form.Get("returnUrl"))
	q.Set("accessNode", r.Form.Get("accessNode"))
	addSWANParams(r, &q, m)
	setFlags(d, &q)
	u.RawQuery = q.Encode()
	return u.String()
}

func setFlags(d *common.Domain, q *url.Values) {
	if d.SwanPostMessage {
		q.Set("postMessageOnComplete", "true")
	} else {
		q.Set("postMessageOnComplete", "false")
	}
	if d.SwanDisplayUserInterface {
		q.Set("displayUserInterface", "true")
	} else {
		q.Set("displayUserInterface", "false")
	}
	if d.SwanUseHomeNode {
		q.Set("useHomeNode", "true")
	} else {
		q.Set("useHomeNode", "false")
	}
}

func addSWANParams(r *http.Request, q *url.Values, m *dialogModel) {
	if m != nil {
		for k, v := range m.Values {
			q.Set(k, v[0])
		}
	}
}
