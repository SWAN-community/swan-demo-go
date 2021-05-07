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
	"bytes"
	"common"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/url"
	"owid"
	"strings"
	"swan"
	"text/template"
)

var complaintSubjectTemplate = newComplaintTemplate(
	"subject",
	"SWAN Complaint: {{ .Organization }}")
var complaintBodyTemplate = newComplaintTemplate("body", `
 To whom it may concern,
 
 I believe that {{ .Organization }} used my personal information without a 
 legal basis on {{ .Date }}. 
 
 I provided you the following permissions for use of this data.
 
	 Personalize Marketing: {{ .Preferences }}
 
 You cryptographically signed this information. We therefore agree that you were
 in posession of the information.
 
 As an organization operating in '{{ .Country }}' you are bound by the following 
 rules.
 
	 {{ .DPRURL }}
 
 I would be grateful if you can respond by email to this address within 7 
 working days.
 
 Regards,
 
 [INSERT YOU NAME]
 
 --- DO NOT CHANGE THE TEXT BELOW THIS LINE ---
 {{ .IDAsString }}
 --- DO NOT CHANGE THE TEXT ABOVE THIS LINE ---`)

// Complaint used to format an email template.
type Complaint struct {
	ID           *swan.ID // The swan.ID that the complaint relates to
	DPRURL       string
	Organization string
	Country      string
	idOWID       *owid.OWID // The ID as an OWID
}

// Date to use in the email template.
func (c *Complaint) Date() string {
	return c.idOWID.Date.Format("2006-01-02 15:01")
}

// SWID to use in the email template.
func (c *Complaint) SWID() string {
	return c.ID.SWID.AsString()
}

// SID to use in the email template.
func (c *Complaint) SID() string {
	return c.ID.SID.AsString()
}

// Preferences string to use in the email template.
func (c *Complaint) Preferences() string {
	return c.ID.PreferencesAsString()
}

// ID as a string
func (c *Complaint) IDAsString() (string, error) {
	return c.ID.AsString()
}

// SWANOWID as a string
func (c *Complaint) SWANOWID() string {
	return c.idOWID.AsString()
}

func newComplaintTemplate(n string, b string) *template.Template {
	t, err := template.New(n).Parse(strings.TrimSpace(b))
	if err != nil {
		panic(err)
	}
	return t
}

func newComplaint(
	cfg *common.Configuration,
	swanOWID *owid.OWID,
	partyOWID *owid.OWID) (*Complaint, error) {
	var err error

	// Set the static information associated with the complaint. These are
	var c Complaint
	c.DPRURL = "URL of the DPR"
	c.Country = "Region of the CMP"

	// Set the ID as an OWID.
	c.idOWID = partyOWID

	// Work out the swan.ID from the ID OWID provided.
	c.ID, err = swan.IDFromOWID(swanOWID)
	if err != nil {
		return nil, err
	}

	// Set the organization as the domain for the moment.
	c.Organization = partyOWID.Domain

	// Return the complain data structure ready for the template email.
	return &c, nil
}

func handlerComplain(
	d *common.Domain,
	w http.ResponseWriter,
	r *http.Request) {

	// Get the form values from the input request.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Check that the SWAN ID and the Party ID are present.
	if r.Form.Get("swanid") == "" {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("'swanid' missing"),
			http.StatusBadRequest)
		return
	}
	if r.Form.Get("partyid") == "" {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("'partyid' missing"),
			http.StatusBadRequest)
		return
	}

	swanOWID, err := owid.FromBase64(r.Form.Get("swanid"))
	if err != nil {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("'swanid' not a valid OWID"),
			http.StatusBadRequest)
		return
	}
	partyOWID, err := owid.FromBase64(r.Form.Get("partyid"))
	if err != nil {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			fmt.Errorf("'partyid' not a valid OWID"),
			http.StatusBadRequest)
		return
	}

	// Create the complaint object.
	c, err := newComplaint(d.Config, swanOWID, partyOWID)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Get the strings for the subject and the body.
	var subject bytes.Buffer
	err = complaintSubjectTemplate.Execute(&subject, c)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	var body bytes.Buffer
	err = complaintBodyTemplate.Execute(&body, c)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Create the URL for the email.
	u := fmt.Sprintf("mailto:info@%s?subject=%s&body=%s",
		c.idOWID.Domain,
		url.PathEscape(subject.String()),
		url.PathEscape(body.String()))

	// Return the URL as a text string.
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, err = g.Write([]byte(u))
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
}
