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

package main

import (
	"demo"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type HTTPSHandler struct {
}

func getPortHTTP() string {
	var port string
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		// Get the port environment variable from Azure App Services.
		port = os.Getenv("HTTP_PLATFORM_PORT")
	} else if os.Getenv("PORT") != "" {
		// Get the port environment variable from Amazon Web Services.
		port = os.Getenv("PORT")
	}
	return port
}

func getPortHTTPS() string {
	var port string
	if os.Getenv("HTTPS_PLATFORM_PORT") != "" {
		// Get the port environment variable from Azure App Services.
		port = os.Getenv("HTTPS_PLATFORM_PORT")
	}
	return port
}

func main() {
	var err error
	var settingsFile string

	// Get the path to the settings file.
	if len(os.Args) >= 2 {
		settingsFile = os.Args[1]
	} else {
		settingsFile = "appsettings.json"
	}

	// Get the ports for HTTP or HTTPS.
	portHttp := getPortHTTP()
	portHttps := getPortHTTPS()

	// Add the SWAN handlers.
	demo.AddHandlers(settingsFile)

	// Start the HTTPS proxy if there is a provided port.
	if portHttps != "" {
		go func() {
			log.Printf("Listenning on HTTPS port: %s\n", portHttps)
			err := http.ListenAndServeTLS(
				fmt.Sprintf(":%s", portHttps),
				"uk.crt",
				"uk.key",
				&HTTPSHandler{})
			if err != nil {
				log.Println(err)
			}
		}()
	}

	// Start the HTTP web server on the port provided.
	log.Printf("Listenning on HTTP port: %s\n", portHttp)
	err = http.ListenAndServe(fmt.Sprintf(":%s", portHttp), nil)
	if err != nil {
		log.Println(err)
	}
}

func (h HTTPSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var u url.URL
	u.Scheme = "http"
	u.Host = r.Host
	p := httputil.NewSingleHostReverseProxy(&u)
	p.ServeHTTP(w, r)
}
