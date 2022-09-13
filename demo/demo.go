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

package demo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/SWAN-community/swan-demo-go/demo/cmp"
	"github.com/SWAN-community/swan-demo-go/demo/marketer"
	"github.com/SWAN-community/swan-demo-go/demo/publisher"
	"github.com/SWAN-community/swan-demo-go/demo/shared"
	"github.com/SWAN-community/swan-demo-go/demo/swanopenrtb"

	swanop "github.com/SWAN-community/swan-op-go"
	"github.com/SWAN-community/access-go"
)

// AddHandlers and outputs configuration information.
func AddHandlers(settingsFile string) error {

	// Get the demo configuration.
	dc := shared.NewConfig(settingsFile)

	// Used a fixed set of access keys for the demo.
	swa := access.NewFixed(dc.AccessKeys)

	// Get all the domains for the SWAN demo.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	domains, err := parseDomains(&dc, filepath.Join(wd, "www"))
	if err != nil {
		return err
	}
	dc.Domains = domains

	// Add the SWAN handlers, with the demo handler being used for any
	// malformed storage requests.
	err = swanop.AddHandlers(
		settingsFile,
		swa,
		shared.Handler(domains))
	if err != nil {
		return err
	}

	// Output details for information.
	log.Printf("Demo scheme: %s\n", dc.Scheme)
	for _, d := range domains {
		log.Printf("%s:%s:%s", d.Category, d.Host, d.Name)
	}
	return nil
}

// parseDomains returns an array of domains (e.g. swan-demo.uk) with all the
// information needed to server static, API and HTML requests. Folder names
// relate to the domain name and must contain a config.json file to be valid.
// c is the general server configuration.
// path provides the root folder where the child folders are the names of the
// domains that the demo responds to.
func parseDomains(
	c *shared.Configuration,
	path string) ([]*shared.Domain, error) {
	var domains []*shared.Domain
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			p := filepath.Join(path, f.Name())
			g := shared.GetConfigFile(p)
			if g != nil {
				domain, err := shared.NewDomain(c, p, g)
				if err != nil {
					return nil, err
				}
				err = addHandler(domain)
				if err != nil {
					return nil, err
				}
				domains = append(domains, domain)
			}
		}
	}
	return domains, nil
}

// Set the HTTP handler for the domain.
func addHandler(d *shared.Domain) error {
	switch d.Category {
	case "CMP":
		d.SetHandler(cmp.Handler)
		break
	case "Publisher":
		d.SetHandler(publisher.Handler)
		break
	case "Advertiser":
		d.SetHandler(marketer.Handler)
		break
	case "DSP":
		d.SetHandler(swanopenrtb.Handler)
		break
	case "SSP":
		d.SetHandler(swanopenrtb.Handler)
		break
	case "DMP":
		d.SetHandler(swanopenrtb.Handler)
		break
	case "Exchange":
		d.SetHandler(swanopenrtb.Handler)
		break
	case "Demo":
		d.SetHandler(shared.HandlerHTML)
		break
	default:
		return fmt.Errorf("Category '%s' invalid for domain '%s'",
			d.Category,
			d.Host)
	}
	return nil
}
