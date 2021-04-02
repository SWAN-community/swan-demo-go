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
	"cmp"
	"common"
	"fmt"
	"io/ioutil"
	"log"
	"marketer"
	"openrtb"
	"os"
	"path/filepath"
	"publisher"
	"swan"
)

// AddHandlers and outputs configuration information.
func AddHandlers(settingsFile string) {

	// Get the demo configuration.
	dc := common.NewConfig(settingsFile)

	// Get the example simple access control implementations.
	swa := swan.NewAccessSimple(dc.AccessKeys)

	// Get all the domains for the SWAN demo.
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	domains, err := parseDomains(&dc, filepath.Join(wd, "www"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dc.Domains = domains

	// Add the SWAN handlers, with the demo handler being used for any
	// malformed storage requests.
	err = swan.AddHandlers(
		settingsFile,
		swa,
		common.Handler(domains))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Output details for information.
	log.Printf("Demo scheme: %s\n", dc.Scheme)
	for _, d := range domains {
		log.Printf("%s:%s:%s", d.Category, d.Host, d.Name)
	}
}

// parseDomains returns an array of domains (e.g. swan-demo.uk) with all the
// information needed to server static, API and HTML requests.
// c is the general server configuration.
// path provides the root folder where the child folders are the names of the
// domains that the demo responds to.
func parseDomains(
	c *common.Configuration,
	path string) ([]*common.Domain, error) {
	var domains []*common.Domain
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {

		if f.Name() == "owid-js" {
			continue
		}

		// Domains are the directories of the folder provided.
		if f.IsDir() {
			domain, err := common.NewDomain(c, filepath.Join(path, f.Name()))
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
	return domains, nil
}

// Set the HTTP handler for the domain.
func addHandler(d *common.Domain) error {
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
		d.SetHandler(openrtb.Handler)
		break
	case "SSP":
		d.SetHandler(openrtb.Handler)
		break
	case "DMP":
		d.SetHandler(openrtb.Handler)
		break
	case "Exchange":
		d.SetHandler(openrtb.Handler)
		break
	case "Demo":
		d.SetHandler(common.HandlerHTML)
		break
	default:
		return fmt.Errorf("Category '%s' invalid for domain '%s'",
			d.Category,
			d.Host)
	}
	return nil
}
