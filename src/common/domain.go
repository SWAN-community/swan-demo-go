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
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"owid"
	"path/filepath"
	"strings"
	"swan"
)

// Domain represents the information held in the domain configuration file
// commonly represented in the demo in config.json.
type Domain struct {
	Category                 string // Category of the domain
	Name                     string // Common name for the domain
	Bad                      bool   // True if this domain is a bad actor for the demo
	Host                     string // The host name for the domain
	SwanMessage              string // Message if used with SWAN
	SwanBackgroundColor      string // Background color if used with SWAN
	SwanMessageColor         string // Message text color if used with SWAN
	SwanProgressColor        string // Message progress color if used with SWAN
	SwanPostMessage          bool   // True if the publisher gets the results from SWAN as a post message
	SwanDisplayUserInterface bool   // True to display the user interface
	SwanUseHomeNode          bool   // True to use the home node if it has current data
	SwanJavaScript           bool   // True to use JavaScript responses rather than HTML documents
	SwanNodeCount            int    // The number of SWAN nodes to use for operations
	CmpNodeCount             int    // The number of nodes to visit when accessing the CMP
	// The domain of the access node used with SWAN (only set for CMPs)
	SWANAccessNode string
	SWANAccessKey  string // The access key to use when communicating with SWAN.
	// The domain of the CMP that will in turn access the SWAN Network via an Operator
	CMP       string
	Suppliers []string           // Suppliers used by the domain operator
	Adverts   []Advert           // Adverts the domain can serve
	Config    *Configuration     // Configuration for the server
	folder    string             // Location of the directory
	templates *template.Template // HTML templates
	owid      *owid.Creator      // The OWID creator associated with the domain if any
	owidStore owid.Store         // The connection to the OWID store
	swan      *swan.Connection   // The connection to SWAN
	// The HTTP handler to use for this domain
	handler func(d *Domain, w http.ResponseWriter, r *http.Request)
}

// GetConfig returns the configuration from the folder, or nil if the
// configuration does not exist.
func GetConfigFile(folder string) *os.File {
	f, _ := os.Open(filepath.Join(folder, "config.json"))
	return f
}

// NewDomain creates a new instance of domain information from the file
// provided.
func NewDomain(
	c *Configuration,
	folder string,
	configFile *os.File) (*Domain, error) {
	var err error

	// Read the configuration file into the domain data structure.
	var d Domain
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&d)

	// Add some additional parameters.
	d.Config = c
	d.Host = filepath.Base(folder)
	d.folder = folder
	d.templates, err = d.parseHTML()
	if err != nil {
		return nil, err
	}
	d.owidStore = c.owid
	d.swan = swan.NewConnection(swan.Operation{
		Client: swan.Client{
			SWAN: swan.SWAN{
				AccessKey: d.SWANAccessKey,
				Operator:  d.SWANAccessNode,
				Scheme:    d.Config.Scheme}},
		BackgroundColor:       d.SwanBackgroundColor,
		Message:               d.SwanMessage,
		MessageColor:          d.SwanMessageColor,
		ProgressColor:         d.SwanProgressColor,
		NodeCount:             d.SwanNodeCount,
		DisplayUserInterface:  d.SwanDisplayUserInterface,
		PostMessageOnComplete: d.SwanPostMessage,
		JavaScript:            d.SwanJavaScript,
		UseHomeNode:           d.SwanUseHomeNode})
	return &d, nil
}

// SetHandler adds a HTTP handler to the domain.
func (d *Domain) SetHandler(fn func(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request)) {
	d.handler = fn
}

// LookupHTML based on the templates available to the domain.
func (d *Domain) LookupHTML(p string) *template.Template {
	if d.templates == nil {
		return nil
	}

	// Try to find the template that relates to the file path.
	t := d.templates.Lookup(filepath.Base(p))

	// If no template can be found try finding one for the category of the
	// domain.
	if t == nil {
		t = d.templates.Lookup(strings.ToLower(d.Category) + ".html")
	}

	// Finally, if no template is found try the default one.
	if t == nil {
		t = d.templates.Lookup("default.html")
	}
	return t
}

func (d *Domain) SWAN() *swan.Connection {
	return d.swan
}

// GetOWIDCreator returns the OWID creator from the OWID store for the the
// domain.
func (d *Domain) GetOWIDCreator() (*owid.Creator, error) {
	var err error
	if d.owid == nil {
		d.owid, err = d.owidStore.GetCreator(d.Host)
		if err != nil {
			return nil, err
		}
		if d.owid == nil {
			return nil, fmt.Errorf(
				"Domain '%s' is not a registered OWID creator. Register the "+
					"domain for the SWAN demo using http[s]://%s/owid/register",
				d.Host,
				d.Host)
		}
	}
	return d.owid, nil
}

func infoRole(s interface{}) string {
	_, fok := s.(*swan.Failed)
	_, bok := s.(*swan.Bid)
	_, eok := s.(*swan.Empty)
	_, iok := s.(*swan.ID)
	if fok {
		return "Failed"
	}
	if bok {
		return "Bid"
	}
	if eok {
		return "Empty"
	}
	if iok {
		return "ID"
	}
	return ""
}

func (d *Domain) parseHTML() (*template.Template, error) {
	var t *template.Template
	files, err := ioutil.ReadDir(d.folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			s, err := ioutil.ReadFile(filepath.Join(d.folder, file.Name()))
			if err != nil {
				return nil, err
			}
			if t == nil {
				t, err = template.New(file.Name()).Funcs(
					template.FuncMap{"role": infoRole}).Parse(
					string(s))
				if err != nil {
					return nil, err
				}
			} else {
				t, err = t.New(file.Name()).Funcs(
					template.FuncMap{"role": infoRole}).Parse(
					string(s))
				if err != nil {
					return nil, err
				}
			}

		}
	}
	return t, nil
}
