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

package publisher

import (
	"bytes"
	"common"
	"encoding/base64"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"openrtb"
	"owid"
	"strings"
	"swan"
	"time"

	"github.com/google/uuid"
)

// Model used with HTML templates.
type Model struct {
	common.PageModel
	swanData []*swan.Pair // The SWAN data for display
}

// CMPURL returns the URL for the CMP dialog.
func (m Model) CMPURL() string {
	return getCMPURL(m.Domain, m.Request, m.swanData)
}

// SWANURL returns the URL for the SWAN operation. Used when the SWAN
// transaction is completed in a pop up window, iFrame or is a JavaScript
// include.
func (m Model) SWANURL() string {
	u, _ := getSWANURL(m.Domain, m.Request, m.swanData)
	return u
}

// HomeNode returns the internet domain name of the home node.
func (m Model) HomeNode() string {
	h, _ := getHomeNode(m.Domain, m.Request)
	return h
}

// IsNew returns true if the SWID is newly created, otherwise false.
func (m Model) IsNew() bool {
	if m.swid() != nil {
		o, _ := m.swid().AsOWID()
		if o != nil {
			return o.Age() <= 1
		}
	}
	return false
}

// Personalized returns a boolean to indicate if personalized marketing is enabled.
func (m Model) Personalized() bool { return m.PrefAsString() == "on" }

// SWIDAsString Secure Web IDentifier
func (m Model) SWIDAsString() string { return common.AsStringFromUUID(m.swid()) }

// SIDAsString Signed in IDentifier
func (m Model) SIDAsString() string { return common.AsPrintable(m.sid()) }

// PrefAsString true if personalized marketing allowed, otherwise false
func (m Model) PrefAsString() string { return common.AsString(m.pref()) }

// SWIDDomain returns the domain that created the SWID OWID
func (m Model) SWIDDomain() string { return common.OWIDDomain(m.swid()) }

// SIDDomain returns the domain that created the SID OWID
func (m Model) SIDDomain() string { return common.OWIDDomain(m.sid()) }

// PrefDomain returns the domain that created the Allow OWID
func (m Model) PrefDomain() string { return common.OWIDDomain(m.pref()) }

// SWIDDate returns the date SWID OWID was created
func (m Model) SWIDDate() string { return common.OWIDDate(m.swid()) }

// SIDDate returns the date SID OWID was created
func (m Model) SIDDate() string { return common.OWIDDate(m.sid()) }

// PrefDate returns the date Allow OWID was created
func (m Model) PrefDate() string { return common.OWIDDate(m.pref()) }

// Stopped returns a list of the domains that have been stopped for advertising.
func (m Model) Stopped() []string {
	return strings.Split(common.AsString(m.stop()), "\r\n")
}

// DomainsByCategory returns all the domains that match the category.
func (m Model) DomainsByCategory(category string) []*common.Domain {
	var domains []*common.Domain
	for _, domain := range m.Domain.Config.Domains {
		if domain.Category == category {
			domains = append(domains, domain)
		}
	}
	return domains
}

// NewAdvertHTML provides the HTML for the advert that will be displayed on the
// web page at the placement provided.
func (m Model) NewAdvertHTML(placement string) (template.HTML, error) {

	// Check that the preference information has a value and is not empty.
	if m.PrefAsString() == "" {
		return template.HTML("<p>Preferences not set</p>"), nil
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Use the SWAN network to generate the Offer ID.
	r, err := m.newOfferNode()
	if err != nil {
		return "", err
	}

	// Add the publishers signature and then process the supply chain.
	_, err = openrtb.SendToSuppliers(m.Domain, r)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	// Get the OWID tree as a base 64 string.
	e, err := r.AsJSON()
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	// Get the winning bid node.
	w, err := swan.WinningNode(r)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	// Get the winning bid.
	b, err := swan.WinningBid(r)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	// Get the return URL.
	t, err := common.GetReturnURL(m.Request)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	// Get the URL for the info icon.
	var i url.URL
	i.Scheme = m.Config().Scheme
	i.Host = m.Domain.CMP
	i.Path = "/info"
	q := i.Query()
	n := w
	for n != nil {
		q.Add("owid", n.GetOWIDAsString())
		n = n.GetParent()
	}
	q.Set("returnUrl", t.String())
	i.RawQuery = q.Encode()

	// Return a FORM HTML element with a button for the advert. The OWID tree
	// is a base 64 string added as a hidden field to the form.
	var html bytes.Buffer
	html.WriteString(fmt.Sprintf("<form method=\"POST\" action=\"//%s\">"+
		"<div class=\"form-group\">"+
		"<input type=\"hidden\" id=\"transaction\" name=\"transaction\" value=\"%s\">"+
		"<button type=\"submit\" id=\"view\" name=\"view\" class=\"advert-button\">"+
		"<img src=\"//%s\">"+
		"</button>"+
		"<a href=\"%s\" class=\"advert-stop\" title=\"Info about this advert\">"+
		"<img src=\"%s\">"+
		"</a>"+
		"</div>"+
		"</form>",
		b.AdvertiserURL,
		base64.RawStdEncoding.EncodeToString(e),
		b.MediaURL,
		i.String(),
		"noun_Info_1582932.svg"))
	return template.HTML(html.String()), nil
}

// SWID Secure Web IDentifier
func (m Model) swid() *swan.Pair { return m.findResult("swid") }

// SID Signed in IDentifier
func (m Model) sid() *swan.Pair { return m.findResult("sid") }

// Allow true if personalized marketing allowed, otherwise false
func (m Model) pref() *swan.Pair { return m.findResult("pref") }

// Stop the list of domains that should not have adverts displayed form.
func (m Model) stop() *swan.Pair { return m.findResult("stop") }

func (m Model) findResult(k string) *swan.Pair {
	for _, n := range m.swanData {
		if strings.EqualFold(k, n.Key) {
			return n
		}
	}
	return nil
}

// newOfferNode returns a new Offer OWID Node.
func (m *Model) newOfferNode() (*owid.Node, error) {
	o, err := m.newOfferOWID()
	if err != nil {
		return nil, err
	}
	b, err := o.AsByteArray()
	if err != nil {
		return nil, err
	}
	return &owid.Node{OWID: b}, nil
}

// Creates a new Offer OWID from the form parameters of the request. If values
// are missing or are invalid then an error is returned.
func (m *Model) newOfferOWID() (*owid.OWID, error) {
	of, err := m.newOffer()
	if err != nil {
		return nil, err
	}
	b, err := of.AsByteArray()
	if err != nil {
		return nil, err
	}
	oc, err := m.Domain.GetOWIDCreator()
	if err != nil {
		return nil, err
	}
	o, err := oc.CreateOWIDandSign(b)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Returns a new unsigned swan.Offer ready to be used the byte array payload in
// an OWID that the caller is generating a a Root Party for the commencement of
// an advertising request.
func (m *Model) newOffer() (*swan.Offer, error) {
	var err error
	o := swan.NewOffer()

	// Get the page placement from the form parameters.
	o.Placement, err = getValue(m.Request, "placement")
	if err != nil {
		return nil, err
	}

	// Set the publisher domain from the request.
	o.PubDomain = m.Request.Host

	// Get the SWID as an OWID.
	o.SWID, err = getOWID(m.Config(), m.Request, m.swid())
	if err != nil {
		return nil, err
	}

	// Get the Signed in Identifier (SID) as an OWID.
	o.SID, err = getOWID(m.Config(), m.Request, m.sid())
	if err != nil {
		return nil, err
	}

	// Get the preferences as an OWID.
	o.Preferences, err = getOWID(m.Config(), m.Request, m.pref())
	if err != nil {
		return nil, err
	}

	// Get the stopped adverts string.
	o.Stopped = offerGetStopped(m.Request)

	// Random one time data is used to ensure the Offer ID is unique for all
	// time.
	o.UUID, err = uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &o, nil
}

// Returns an array of stopped advert IDs. As the parameter is optional no error
// is returned.
func offerGetStopped(r *http.Request) []string {
	return strings.Split(r.FormValue("stop"), " ")
}

// Returns the value as a string associated with the form parameter k. If the
// value is missing then an error is returned.
func getValue(r *http.Request, k string) (string, error) {
	v := r.FormValue(k)
	if v == "" {
		return "", fmt.Errorf("missing '%s' parameter", k)
	}
	return v, nil
}

// Returns a verified OWID associated with the form parameter k. Returns and
// error if the value is missing or is not a verified OWID.
func getOWID(
	c *common.Configuration,
	r *http.Request,
	p *swan.Pair) (*owid.OWID, error) {
	o, err := owid.FromBase64(p.Value)
	if err != nil {
		return nil, err
	}
	e, err := o.Verify(c.Scheme)
	if err != nil {
		return nil, err
	}
	if e == false {
		return nil, fmt.Errorf("'%s' not a valid OWID", p.Key)
	}
	return o, nil
}
