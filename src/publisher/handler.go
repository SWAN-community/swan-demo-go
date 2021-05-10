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
	"common"
	"compress/gzip"
	"fmt"
	"fod"
	"log"
	"net/http"
	"net/url"
	"swan"
	"swanop"
	"time"
)

// Handler for publisher web pages.
func Handler(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	// Check to see if this request is for an advert.
	if r.URL.Path == "/advert" {
		HandlerAdvert(d, w, r)
		return
	}

	// Try the URL path for the preference values.
	p, ae := newSWANDataFromPath(d, r)
	if ae != nil {

		// If the data can't be decrypted rather than another type of error then
		// redirect to the CMP dialog.
		if ae.StatusCode() >= 400 && ae.StatusCode() < 500 {
			if d.SwanPostMessage == false {
				http.Redirect(w, r, getCMPURL(d, r, nil), 303)
			} else {
				handlerPublisherPage(d, w, r, p)
			}
			return
		}
		common.ReturnServerError(d.Config, w, ae)
		return
	}
	if p != nil {
		redirectToCleanURL(d.Config, w, r, p)
		return
	}

	// If the path does not contain any values then get them from the cookies.
	if p == nil {
		var err error
		p, err = newSWANDataFromCookies(r)
		if err != nil && d.Config.Debug {
			log.Println(err.Error())
		}
	}

	// If the request is from a crawler than ignore SWAN.
	c, err := fod.GetCrawlerFrom51Degrees(r)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}
	if c {
		handlerPublisherPage(d, w, r, p)
		return
	}

	// If there is valid SWAN data then display the page using the page handler.
	// If the SWAN data is not complete, valid, or needs revalidating because it
	// might be old then ask the user to verify or add the required data via the
	// User Interface Provider redirect action.
	// If the SWAN data is not present or invalid then redirect to SWAN to
	// get the latest data.
	if p != nil && len(p) > 0 {
		if isSet(p) {

			// Check to see if the values need to be revalidated.
			if d.SwanPostMessage == false &&
				d.SwanJavaScript == false &&
				revalidateNeeded(p) {
				redirectToSWANFetch(d, w, r, p)
			} else {
				handlerPublisherPage(d, w, r, p)
			}
		} else {
			http.Redirect(w, r, getCMPURL(d, r, p), 303)
		}
	} else {
		if d.SwanPostMessage == false && d.SwanJavaScript == false {
			redirectToSWANFetch(d, w, r, p)
		} else {
			handlerPublisherPage(d, w, r, p)
		}
	}
}

func handlerPublisherPage(
	d *common.Domain,
	w http.ResponseWriter,
	r *http.Request,
	p []*swan.Pair) {
	t := d.LookupHTML(r.URL.Path)
	if t == nil {
		http.NotFound(w, r)
		return
	}
	var m Model
	m.Domain = d
	m.Request = r
	m.swanData = p
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	err := t.Execute(g, &m)
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
	}
}

func newSWANDataFromCookies(r *http.Request) ([]*swan.Pair, error) {
	var p []*swan.Pair
	for _, c := range r.Cookies() {
		if swan.IsSWANCookie(c) {
			p = append(p, swan.NewPairFromCookie(c))
		}
	}
	return p, nil
}

func newSWANData(
	d *common.Domain,
	v string) ([]*swan.Pair, *swan.Error) {
	return d.SWAN().Decrypt(v)
}

// Get the section of the URL that has the SWAN data.
func newSWANDataFromPath(
	d *common.Domain,
	r *http.Request) ([]*swan.Pair, *swan.Error) {
	b := common.GetSWANDataFromRequest(r)
	if b == "" {
		return nil, nil
	}
	return newSWANData(d, b)
}

// SWAN data could be obtained from the URL. Remove the SWAN data string from
// the URL and redirect back to the page. Set cookies in the redirect so that
// the data is persisted.
func redirectToCleanURL(
	c *common.Configuration,
	w http.ResponseWriter,
	r *http.Request,
	p []*swan.Pair) {
	u := common.GetCleanURL(c, r).String()
	if c.Debug {
		log.Printf("Redirecting to '%s'\n", u)
	}
	setCookies(r, w, p)
	http.Redirect(w, r, u, 303)
}

// Redirect back to the current URL after fetching the SWAN data. If SWAN data
// does not exist then use the values contained in the swan pairs provided in
// parameter p.
func redirectToSWANFetch(
	d *common.Domain,
	w http.ResponseWriter,
	r *http.Request,
	p []*swan.Pair) {
	u, err := getSWANURL(d, r, p)
	if err != nil {
		common.ReturnProxyError(d.Config, w, err)
		return
	}
	http.Redirect(w, r, u, 303)
}

func getSWANURL(
	d *common.Domain,
	r *http.Request,
	p []*swan.Pair) (string, *swan.Error) {
	return d.SWAN().NewFetch(
		r,
		common.GetCleanURL(d.Config, r).String(),
		p).GetURL()
}

func getHomeNode(
	d *common.Domain,
	r *http.Request) (string, *swan.Error) {
	return d.SWAN().HomeNode(r)
}

func setCookies(r *http.Request, w http.ResponseWriter, p []*swan.Pair) {
	s := r.URL.Scheme == "https"
	for _, i := range p {
		if i.Value != "" {
			c := i.AsCookie(r, w, s)
			http.SetCookie(w, c)
		}
	}
}

// Returns the CMP preferences URL.
func getCMPURL(d *common.Domain, r *http.Request, p []*swan.Pair) string {
	var u url.URL
	u.Scheme = d.Config.Scheme
	u.Host = d.CMP
	u.Path = "/preferences/"
	q := u.Query()
	q.Set("returnUrl", common.GetCleanURL(d.Config, r).String())
	q.Set("accessNode", d.SWANAccessNode)
	if d.CmpNodeCount > 0 {
		q.Set("nodeCount", fmt.Sprintf("%d", d.CmpNodeCount))
	}
	addSWANParams(r, &q, p)
	setFlags(d, &q)
	u.RawQuery = q.Encode()
	return u.String()
}

// Add the SWAN data values known the to publisher. Used for default values if
// others do not already exist in the SWAN network.
func addSWANParams(r *http.Request, q *url.Values, p []*swan.Pair) {
	if p != nil {
		for _, i := range p {
			q.Set(i.Key, i.Value)
		}
	}
}

// isSet returns true if all three of the values are present in the results and
// are valid OWIDs.
func isSet(d []*swan.Pair) bool {
	c := 0
	for _, e := range d {
		if e.Key == "pref" || e.Key == "swid" || e.Key == "sid" {
			o, err := e.AsOWID()
			if err != nil {
				return false
			}
			if len(o.Payload) > 0 {
				c++
			}
		}
	}
	return c == 3
}

// Get the revalidation time from the swan validation cookie if present. Then
// check to see if the time has elapsed. If so then return true to indicate the
// SWAN data needs to be revalidated with the SWAN Operator.
func revalidateNeeded(d []*swan.Pair) bool {
	for _, e := range d {
		if e.Key == "val" {
			t, err := time.Parse(swanop.ValidationTimeFormat, e.Value)
			if err != nil {
				return true
			}
			return time.Now().UTC().After(t)
		}
	}
	return false
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
	if d.SwanJavaScript {
		q.Set("javaScript", "true")
	} else {
		q.Set("javaScript", "false")
	}
}
