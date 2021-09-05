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
	"net/http"
	"regexp"
	"strings"
)

// Regular expression to check for none base64 encoded characters.
var swanDataRegex *regexp.Regexp

func init() {
	swanDataRegex, _ = regexp.Compile("[^\\w|\\-|\\=]+")
}

// GetSWANDataFromRequest returns the base 64 SWAN data from the request if
// present, otherwise an empty string.
func GetSWANDataFromRequest(r *http.Request) string {
	s := ""

	// Get the SWAN data from the last segment of the PATH or the query string.
	i := strings.LastIndex(r.URL.Path, "/")
	if i >= 0 {
		s = r.URL.Path[i+1:]
	} else {
		s = r.URL.RawQuery
	}

	// Validate that the string is base 64 characters.
	if swanDataRegex.Match([]byte(s)) {
		return ""
	}

	// Return the string.
	return s
}
