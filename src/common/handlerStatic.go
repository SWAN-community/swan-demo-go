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
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Map of allowed static files and whether request for file should support CORS
var allowList = map[string]bool{
	"swan.json": true,
}

// handlerStatic locates and returns static content if relevant to the HTTP
// request. True is returned if static content was returned, otherwise false.
func handlerStatic(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request) (bool, error) {
	var err error
	folder := d.folder
	found := false

	// If the request is for the favicon.ico then rename the path.
	if r.URL.Path == "/favicon.ico" {
		r.URL.Path = "/noun_Swan_3263882.svg"
	}

	for found == false && strings.Contains(folder, "www") {
		found, err = handleStaticFolder(d, w, r, folder)
		if err != nil {
			return false, err
		}
		if found == false {
			folder = filepath.Dir(folder)
		}
	}
	return found, nil
}

func handleStaticFolder(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request,
	folder string) (bool, error) {

	// Get all the file system items in the folder provided.
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return false, &SWANError{err, nil}
	}

	// Loop through the file system items to find any that match the request.
	for _, f := range files {

		// If this file system item is a file and it matches the name of the
		// file requested in the URL path then serve the file.
		if f.IsDir() == false &&
			f.Name() == filepath.Base(r.URL.Path) {
			return handlerFile(w, r, filepath.Join(folder, f.Name())), nil
		}

		// If this file system item is a directory and it matches the last part
		// of the request URL path then evaluate the folder for the requested
		// static file.
		if f.IsDir() == true &&
			strings.HasSuffix(filepath.Dir(r.URL.Path), f.Name()) {
			return handleStaticFolder(d, w, r, folder+filepath.Dir(r.URL.Path))
		}
	}
	return false, nil
}

func handlerFile(
	w http.ResponseWriter,
	r *http.Request,
	f string) bool {
	switch filepath.Ext(f) {
	case ".ico":
		http.ServeFile(w, r, f)
		return true
	case ".jpeg":
		http.ServeFile(w, r, f)
		return true
	case ".jpg":
		http.ServeFile(w, r, f)
		return true
	case ".png":
		http.ServeFile(w, r, f)
		return true
	case ".css":
		http.ServeFile(w, r, f)
		return true
	case ".js":
		http.ServeFile(w, r, f)
		return true
	case ".svg":
		http.ServeFile(w, r, f)
		return true
	case ".map":
		http.ServeFile(w, r, f)
		return true
	default:
		if c, ok := allowList[filepath.Base(f)]; ok {
			if c {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			http.ServeFile(w, r, f)
			return true
		}
		return false
	}
}
