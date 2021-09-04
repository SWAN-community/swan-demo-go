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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/SWAN-community/swan-go"
)

// Handler for all HTTP requests to domains controlled by the demo.
func Handler(d []*Domain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set to true if a domain is found and handled.
		found := false

		// Default domain to use if another domain can't be found.
		var def *Domain

		// r.Host may include the port number or www prefixes or other
		// characters dependent on the environment. Using strings.Contains
		// rather than testing for equality eliminates these issues for a demo
		// where the domain names are not sub strings of one another.
		for _, domain := range d {
			if strings.EqualFold(r.Host, domain.Host) {

				handleDomain(domain, w, r)

				// Mark as the domain being found and then break.
				found = true
				break
			}

			// Record the default domain if flagged.
			if def == nil && domain.Default {
				def = domain
			}
		}

		// All handlers have been tried and nothing has been found. Return the
		// default domain if available.
		if found == false && def != nil {
			handleDomain(def, w, r)
			found = true
		}

		// Return not found.
		if found == false {
			http.NotFound(w, r)
		}
	}
}

// NewError creates an error instance that includes the details of the
// response returned. This is needed to pass the correct status codes and
// context back to the caller.
func NewError(c *Configuration, r *http.Response) *swan.Error {
	var u string
	in, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &swan.Error{Err: err}
	}
	if c.Debug {
		u = r.Request.URL.String()
	} else {
		u = r.Request.Host
	}
	return &swan.Error{
		Err: fmt.Errorf("SWAN '%s' status '%d' message '%s'",
			u,
			r.StatusCode,
			strings.TrimSpace(string(in))),
		Response: r}
}

// ReturnProxyError returns an error where the request is related to a proxy
// request being passed to another end point.
func ReturnProxyError(c *Configuration, w http.ResponseWriter, e *swan.Error) {
	s := http.StatusInternalServerError
	if e.Response != nil {
		s = e.Response.StatusCode
	}
	ReturnStatusCodeError(c, w, e.Err, s)
}

// ReturnServerError returns an internal server error.
func ReturnServerError(c *Configuration, w http.ResponseWriter, err error) {
	ReturnStatusCodeError(c, w, err, http.StatusInternalServerError)
}

// ReturnStatusCodeError returns the HTTP status code specified.
func ReturnStatusCodeError(
	c *Configuration,
	w http.ResponseWriter,
	e error,
	code int) {
	http.Error(w, e.Error(), code)
	if c.Debug {
		println(e.Error())
	}
}

// GetCleanURL returns a URL with the SWAN data removed and no query string
// parameters.
func GetCleanURL(c *Configuration, r *http.Request) *url.URL {
	var u url.URL
	u.Scheme = c.Scheme
	u.Host = r.Host
	u.Path = strings.ReplaceAll(
		r.URL.Path,
		GetSWANDataFromRequest(r),
		"")
	u.RawQuery = ""
	return &u
}

// GetReturnURL returns a parsed URL from the query string, or if not present
// from the referer HTTP header.
func GetReturnURL(r *http.Request) (*url.URL, error) {
	u, err := url.Parse(r.Form.Get("returnUrl"))
	if err != nil {
		return nil, err
	}
	if u == nil || u.String() == "" {
		u, err = url.Parse(r.Header.Get("Referer"))
		if err != nil {
			return nil, err
		}
	}
	u.RawQuery = ""
	return u, nil
}

// GetCurrentPage returns the current request URL.
func GetCurrentPage(c *Configuration, r *http.Request) *url.URL {
	var u url.URL
	u.Scheme = c.Scheme
	u.Host = r.Host
	u.Path = r.URL.Path
	return &u
}

func handleDomain(d *Domain, w http.ResponseWriter, r *http.Request) {

	// Try static resources first.
	f, err := handlerStatic(d, w, r)
	if err != nil {
		ReturnServerError(d.Config, w, err)
		return
	}

	// If not found then use the domain handler.
	if f == false {
		d.handler(d, w, r)
	}
}
