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
	"html/template"
	"log"
	"net/http"
	"net/url"
	"owid"
	"reflect"
	"salt"
	"strconv"
	"strings"
	"swan"

	uuid "github.com/satori/go.uuid"
)

// dialogModel key value pairs with functions to interpret them.
type dialogModel struct {
	url.Values
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

	// Parse the form variables.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// All GET requests are find encrypted data from the URL and redirect to get
	// data if none is found.
	if r.Method == "GET" {

		// No parameters were provided so get the SWAN data from the request
		// path. If no data is present then redirect to SWAN.
		s := common.GetSWANDataFromRequest(r)
		if s == "" {
			redirectToSWAN(d, w, r)
			return
		}

		// Call the SWAN access node for the CMP to turn the data provided in
		// the URL into usable data for the dialog.
		e := decryptAndDecode(d, s, &r.Form)
		if e != nil {

			// If the data can't be decrypted rather than another type of
			// error then redirect via SWAN to the dialog.
			if e.StatusCode() >= 400 && e.StatusCode() < 500 {
				redirectToSWAN(d, w, r)
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

	// If there is no SWID then add a new one.
	if len(r.Form["swid"]) == 0 {
		se := setNewSWID(d, &r.Form)
		if se != nil {
			common.ReturnProxyError(d.Config, w, se)
			return
		}
	}

	// If this is a close request then don't update the values and just return
	// to the return URL.
	if r.Form.Get("close") != "" {
		http.Redirect(w, r, r.Form.Get("returnUrl"), 303)
		return
	}

	// If the method is POST is used then check for any resets.
	if r.Method == "POST" {
		se := dialogReset(d, &r.Form)
		if se != nil {
			common.ReturnProxyError(d.Config, w, se)
			return
		}
	}

	// If the update action is requested them start that process.
	if len(r.Form["update"]) != 0 {

		// The user has request that the data be updated in the SWAN network.

		// Get the OWID creator which is needed to sign the data just captured.
		c, err := d.GetOWIDCreator()
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}

		// Prepare the SWAN update operation.
		o, err := getUpdate(d, r, &r.Form)

		// Set the parameters for the update.
		err = o.SetPref(c, r.Form.Get("pref") == "on")
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
			return
		}
		err = o.SetEmail(c, r.Form.Get("email"))
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
			return
		}
		var s *salt.Salt
		if sf := r.Form.Get("salt"); sf != "" {
			s, err = salt.FromBase64(sf)
			if err != nil {
				common.ReturnStatusCodeError(
					d.Config,
					w,
					err,
					http.StatusBadRequest)
				return
			}
		}
		err = o.SetSalt(c, s)
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
			return
		}
		err = o.SetSWID(r.Form.Get("swid"))
		if err != nil {
			common.ReturnStatusCodeError(
				d.Config,
				w,
				err,
				http.StatusBadRequest)
			return
		}

		// Set the redirection URL for the operation to store the data. Web
		// browser will then be redirected to that URL, the data saved and the
		// return URL for the publisher returned to.
		u, se := o.GetURL()
		if se != nil {
			common.ReturnProxyError(d.Config, w, se)
			return
		}

		// Send the email if the SMTP server is setup.
		if o.Email().PayloadAsString() != "" &&
			strings.Contains(o.Email().PayloadAsString(), "@") {
			err = sendReminderEmail(d, o)
			if err != nil {
				log.Println(err)
			}
		}

		// Redirect the response to the return URL.
		http.Redirect(w, r, u, 303)

	} else {

		// The dialog needs to be displayed. Use the cmp.html template for the
		// user interface.
		g := gzip.NewWriter(w)
		defer g.Close()
		w.Header().Set("Content-Encoding", "gzip")
		err := d.LookupHTML("cmp.html").Execute(g, &dialogModel{Values: r.Form})
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}
	}
}

// sendReminderEmail sends the reminder email with a link to setup other
// browsers.
func sendReminderEmail(d *common.Domain, o *swan.Update) error {

	// Get the salt to display the grid in the email.
	s, err := salt.FromByteArray(o.Salt().Payload)
	if err != nil {
		return err
	}

	// Set the URL using the parameters contained in the update operation.
	u := url.URL{
		Scheme: d.Config.Scheme,
		Host:   d.Host,
		Path:   "/update"}
	q, err := o.GetValues()
	if err != nil {
		return err
	}
	u.RawQuery = q.Encode()

	// Set the email with the model populated.
	err = common.NewSMTP().Send(
		o.Email().PayloadAsString(),
		"SWAN Demo: Email Reminder",
		d.LookupHTML("email-template.html"),
		ModelEmail{Salt: s, PreferencesUrl: u.String()})
	if err != nil {
		return err
	}

	return nil
}

// dialogReset checks for any reset keys and removes other keys if present. If
// these keys are present they are removed from the collection to avoid being
// added as hidden fields.
func dialogReset(d *common.Domain, m *url.Values) *swan.Error {

	// Check to see if the post is as a result of the SWID reset. If so then
	// replace the SWID with a new random value.
	if m.Get("reset-swid") != "" {
		m.Del("reset-swid")
		return setNewSWID(d, m)
	}

	// Check to see if the email and salt are being reset.
	if m.Get("reset-email-salt") != "" {
		m.Set("email", "")
		m.Set("salt", "")
		m.Del("reset-email-salt")
		return nil
	}

	// Check to see if the post is as a result for all data.
	if m.Get("reset-all") != "" {
		m.Set("email", "")
		m.Set("salt", "")
		m.Set("pref", "")
		m.Del("reset-all")
		return setNewSWID(d, m)
	}

	return nil
}

// setNewSWID creates a new SWID and adds to the key values.
func setNewSWID(d *common.Domain, m *url.Values) *swan.Error {
	o, err := d.SWAN().CreateSWID()
	if err != nil {
		return err
	}
	m.Set("swid", o.AsString())
	return nil
}

// getUpdate returns a populated SWAN Update operation.
func getUpdate(
	d *common.Domain,
	r *http.Request,
	m *url.Values) (*swan.Update, error) {

	// Configure the update operation from this demo domain's configuration.
	returnUrl, err := url.Parse(m.Get("returnUrl"))
	if err != nil {
		return nil, err
	}
	u := d.SWAN().NewUpdate(r, returnUrl.String())

	// Use the form to get any information from the initial storage operation
	// to configure the update storage operation.
	if m.Get("accessNode") != "" {
		u.AccessNode = m.Get("accessNode")
	}
	if m.Get("backgroundColor") != "" {
		u.BackgroundColor = m.Get("backgroundColor")
	}
	if m.Get("displayUserInterface") != "" {
		u.DisplayUserInterface = m.Get("displayUserInterface") == "true"
	}
	if m.Get("javaScript") != "" {
		u.JavaScript = m.Get("javaScript") == "true"
	}
	if m.Get("message") != "" {
		u.Message = m.Get("message")
	}
	if m.Get("messageColor") != "" {
		u.MessageColor = m.Get("messageColor")
	}
	if m.Get("postMessageOnComplete") != "" {
		u.PostMessageOnComplete = m.Get("postMessageOnComplete") == "true"
	}
	if m.Get("progressColor") != "" {
		u.ProgressColor = m.Get("progressColor")
	}
	if m.Get("title") != "" {
		u.Title = m.Get("title")
	}
	if m.Get("useHomeNode") != "" {
		u.UseHomeNode = m.Get("useHomeNode") == "true"
	}
	return u, nil
}

// decryptAndDecode the encrypted data returned from SWAN.
func decryptAndDecode(
	d *common.Domain,
	v string,
	m *url.Values) *swan.Error {
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

// redirectToSWAN redirects the request to SWAN to return to this URL with the
// current SWAN data.
func redirectToSWAN(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	// Create the fetch function returning to this URL.
	f := d.SWAN().NewFetch(r, common.GetCleanURL(d.Config, r).String(), nil)

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
