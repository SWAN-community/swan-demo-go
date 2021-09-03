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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"swift"
	"time"
)

/**
 * Command line application used to add or update SWIFT nodes to a SWIFT store.
 */

func main() {
	var settingsFile string
	scanner := bufio.NewScanner(os.Stdin)

	if len(os.Args) >= 2 {
		settingsFile = os.Args[1]
	} else {
		settingsFile = "appsettings.json"
	}

	c := swift.NewConfig(settingsFile)
	s := swift.NewStore(c)
	svc := swift.NewStorageService(c, s...)

	fmt.Println("Add or update a node...")

	var r swift.Register
	r.Store = ""
	r.Network = ""
	r.Domain = ""
	r.Starts = time.Now().UTC().AddDate(0, 0, 1)
	r.Expires = time.Now().UTC().AddDate(0, 3, 0)
	r.Role = 1

	success := false
	isUpdate := false
	for !success {

		// Get the store to use
		if r.StoreError != "" {
			fmt.Println(r.StoreError)
		}
		fmt.Println("Select a store to use in the add or update operation:")
		for _, s := range svc.GetStoreNames() {
			fmt.Printf("\t- %s\r\n", s)
		}
		fmt.Printf("Store [%s]: ", r.Store)
		store, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		r.Store = store

		// Get the network name
		if r.NetworkError != "" {
			fmt.Println(r.NetworkError)
		}
		fmt.Printf("Node network name [%s]: ", r.Network)
		network, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		r.Network = network

		// Get the node domain
		fmt.Printf("Node domain [%s]: ", r.Domain)
		domain, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		r.Domain = domain

		// get the node startdate

		fmt.Printf("Start date [%s]: ", r.Starts.Format(time.RFC3339))
		starts, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		if starts != "" {
			r.Starts, err = time.Parse(time.RFC3339, starts)
			if err != nil {
				panic(err)
			}
		}

		// get the node expiry date
		if r.ExpiresError != "" {
			fmt.Println(r.ExpiresError)
		}
		fmt.Printf("Expires date [%s]: ", r.Expires.Format(time.RFC3339))
		expires, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		if expires != "" {
			r.Expires, err = time.Parse(time.RFC3339, expires)
			if err != nil {
				panic(err)
			}
		}

		// get the node role
		fmt.Println("Select a node role:")
		printNodeTypes()
		fmt.Printf("Role [%d]: ", r.Role)
		role, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		if role != "" {
			i, err := strconv.Atoi(role)
			if err != nil {
				panic(err)
			}

			r.Role = i
		}

		// check if details are correct
		fmt.Println()
		fmt.Println("Confirm the details are correct!:")
		fmt.Printf("\tStore: %s\r\n", r.Store)
		fmt.Printf("\tNetwork Name: %s\r\n", r.Network)
		fmt.Printf("\tDomain Name: %s\r\n", r.Domain)
		fmt.Printf("\tStart Date: %s\r\n", r.Starts.Format(time.RFC3339))
		fmt.Printf("\tExpiry Date: %s\r\n", r.Expires.Format(time.RFC3339))
		fmt.Printf("\tRole: %d\r\n", r.Role)
		fmt.Printf("Correct? (Y/n) [Y]: ")
		correct, err := scan(scanner)
		if err != nil {
			panic(err)
		}
		if correct != "y" &&
			correct != "Y" &&
			correct != "" {
			continue
		}

		// set the node
		success, isUpdate = svc.SetNode(&r)
		if !success {
			fmt.Println("There were some errors, check your values:")
		}
	}

	// confirmation
	if isUpdate {
		fmt.Println("Node updated!")
	} else {
		fmt.Println("Node added!")
	}
}

func scan(scanner *bufio.Scanner) (string, error) {
	scanner.Scan()
	res := scanner.Text()
	if scanner.Err() != nil {
		return "", scanner.Err()
	}
	return res, nil
}

func printNodeTypes() {
	fmt.Println("\t(0) Access Node")
	fmt.Println("\t(1) Storage Node")
	fmt.Println("\t(2) Sharing Node")
}
