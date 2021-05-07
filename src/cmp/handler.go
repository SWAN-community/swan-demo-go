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
	"net/http"
	"strings"
)

// Handler for the CMP features.
func Handler(d *common.Domain, w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/preferences") {
		handlerDialog(d, w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/stop") {
		handlerStop(d, w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/info") {
		handlerInfo(d, w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/complain") {
		handlerComplain(d, w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/update") {
		handlerUpdate(d, w, r)
		return
	}
}
