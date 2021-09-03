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
	"os"
	"owid"
)

// Configuration maps to the appsettings.json settings file.
type Configuration struct {
	AccessKeys []string   `json:"accessKeys"` // Array of valid keys for SWAN access
	Scheme     string     `json:"scheme"`     // The scheme to use for requests
	Debug      bool       `json:"debug"`      // True if debug HTML output should be provided
	Domains    []*Domain  // All the domains that form the demo
	owid       owid.Store // The OWID store for use with domains
}

// NewConfig creates a new instance of configuration from the file provided.
func NewConfig(settingsFile string) Configuration {
	var c Configuration
	configFile, err := os.Open(settingsFile)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	c.owid = getOWIDStore(settingsFile)
	return c
}

func getOWIDStore(settingsFile string) owid.Store {
	owidConfig := owid.NewConfig(settingsFile)
	err := owidConfig.Validate()
	if err != nil {
		panic(err)
	}
	return owid.NewStore(owidConfig)
}
