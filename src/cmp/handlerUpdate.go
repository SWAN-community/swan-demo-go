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
)

func handlerUpdate(
	d *common.Domain,
	w http.ResponseWriter,
	r *http.Request) {

	// Get the form values from the input request.
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(d.Config, w, err)
		return
	}

	// Create the update operation.
	o, err := getUpdate(d, r, &r.Form)

	// Add the values from the form parameters.
	err = o.SetPrefFromOWID(r.Form.Get("pref"))
	if err != nil {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			err,
			http.StatusBadRequest)
		return
	}
	err = o.SetEmailFromOWID(r.Form.Get("email"))
	if err != nil {
		common.ReturnStatusCodeError(
			d.Config,
			w,
			err,
			http.StatusBadRequest)
		return
	}
	err = o.SetSaltFromOWID(r.Form.Get("salt"))
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

	// Redirect the response to the return URL.
	http.Redirect(w, r, u, 303)
}
