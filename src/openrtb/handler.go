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

package openrtb

import (
	"bytes"
	"common"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"owid"
	"swan"
	"sync"
)

var empty swan.Empty // Used for empty responses

const openRTBPath = "/demo/api/v1/bid" // The path for this handler

// Handler is responsible for a real time transaction for advertising.
// The body of the request must contain a JSON array of Processor IDs which
// contain the signature of the last entry in the list of Processors.
func Handler(d *common.Domain, w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == openRTBPath && r.Method == "POST" {

		// Unpack the body of the request to form the bid data structure.
		o, err := getOffer(d, r)
		if err != nil {
			common.ReturnStatusCodeError(d.Config, w, err, http.StatusBadRequest)
			return
		}

		// If this domain is a bad actor then change the publisher's domain to
		// one that would generate more money from advertising.
		if d.Bad {
			err = changePubDomain(o, "high-value-pub.com")
			if err != nil {
				common.ReturnStatusCodeError(
					d.Config,
					w,
					err,
					http.StatusInternalServerError)
				return
			}
		}

		// Handle the bid and return if the URL was found.
		t, err := HandleTransaction(d, o)
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}

		// The caller already knows about the rest of the tree. Only return this
		// Processor OWID and the children.
		b, err := t.AsJSON()
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}

		g := gzip.NewWriter(w)
		defer g.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		_, err = g.Write(b)
		if err != nil {
			common.ReturnServerError(d.Config, w, err)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

func getSWANOffer(r *owid.Node) (*swan.Offer, error) {
	f, err := r.GetOWID()
	if err != nil {
		return nil, err
	}
	o, err := swan.OfferFromOWID(f)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func changePubDomain(r *owid.Node, newPubDomain string) error {
	f, err := r.GetOWID()
	if err != nil {
		return err
	}
	o, err := swan.OfferFromOWID(f)
	if err != nil {
		return err
	}
	o.PubDomain = newPubDomain
	f.Payload, err = o.AsByteArray()
	if err != nil {
		return err
	}
	r.OWID, err = f.AsByteArray()
	if err != nil {
		return err
	}
	return nil
}

// HandleTransaction processes an OpenRTB transaction.
func HandleTransaction(d *common.Domain, n *owid.Node) (*owid.Node, error) {

	// Verify that this domain can create OWIDs. Failure to register a domain
	// as an OWID creator is a common setup mistake.
	oc, err := d.GetOWIDCreator()
	if err != nil {
		return nil, err
	}

	// The single leaf is the parent Processor OWID. If there isn't a single
	// leaf then too much information has been sent from the caller.
	parent, err := n.GetLeaf()
	if err != nil {
		return nil, err
	}

	// Create an OWID for this processor.
	t, err := oc.CreateOWID(nil)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("Could not create new OWID")
	}

	// If this domain has adverts then choose one at random. Get a random
	// byte array to use as the payload from the Processor OWID.
	if len(d.Adverts) > 0 {

		// The root node must be the Offer.
		offer, err := swan.OfferFromNode(n.GetRoot())
		if err != nil {
			return nil, err
		}

		// Get a random advert checking that it is not on the stopped list.
		var b swan.Bid
		i := 10
		for i > 0 {
			w := d.Adverts[rand.Intn(len(d.Adverts))]
			if offer.IsStopped(w.AdvertiserURL) == false {
				b.AdvertiserURL = w.AdvertiserURL
				b.MediaURL = w.MediaURL
				t.Payload, err = b.AsByteArray()
				break
			}
			i--
		}
		if i == 0 {
			t.Payload, err = empty.AsByteArray()
		}
	} else {
		t.Payload, err = empty.AsByteArray()
	}
	if err != nil {
		return nil, err
	}

	// Sign the Processor OWID with the root OWID now that it's part of the
	// tree. This can be used by down stream suppliers to verify that this
	// processor was involved in the transaction.
	r, err := n.GetOWID()
	if err != nil {
		return nil, err
	}
	err = oc.Sign(t, r)
	if err != nil {
		return nil, err
	}

	// Add this signed Processor OWID to the children of the parent.
	n, err = parent.AddOWID(t)
	if err != nil {
		return nil, err
	}

	// Send the transaction on to any suppliers.
	if len(d.Suppliers) > 0 {
		return SendToSuppliers(d, n)
	}
	return n, nil
}

func SendToSuppliers(d *common.Domain, n *owid.Node) (*owid.Node, error) {
	var err error

	// Call all the suppliers adding them to this Processor OWID's child
	// transactions.
	var wg sync.WaitGroup
	wg.Add(len(d.Suppliers))
	c := make([]*owid.Node, len(d.Suppliers))
	e := make([]error, len(d.Suppliers))
	for i, s := range d.Suppliers {
		go func(i int, s string) {
			defer wg.Done()
			c[i], e[i] = sendToSupplier(d, s, n)
		}(i, s)
	}
	wg.Wait()

	// Merge the results from the suppliers.
	i := 0
	for i < len(d.Suppliers) {
		if e[i] != nil {
			return nil, e[i]
		}
		if c[i] != nil {
			n.AddChild(c[i])
		}
		i++
	}

	// If there are children then pick one at random for the payload of
	// this processor. Used to determine the winner when the transaction
	// is complete. This also demonstrates how the payload can be changed
	// after the response has been received.
	if len(n.Children) > 0 {
		n.Value, err = chooseWinner(n)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

func chooseWinner(n *owid.Node) (int, error) {
	w := -1
	e := 0
	o := make([]bool, len(n.Children))
	for i, c := range n.Children {
		b, err := isBid(c)
		if err != nil {
			return -1, err
		}
		if b {
			o[i] = true
			e++
		}
	}
	if e > 0 {
		for w < 0 {
			i := rand.Intn(len(o))
			if o[i] {
				w = i
			}
		}
	}
	return w, nil
}

// isBid returns true if the node is related to an eligible bid, otherwise
// false. Eligible bids nodes are either ones that contain bid information
// or the Value of the node is a number that is greater than or equal to 0 and
// there are children.
func isBid(n *owid.Node) (bool, error) {
	if n.Value != nil && n.Value.(float64) >= 0 && len(n.Children) > 0 {
		return true, nil
	}
	b, err := swan.FromNode(n)
	if err != nil {
		return false, err
	}
	_, ok := b.(*swan.Bid)
	return ok, nil
}

func getOffer(d *common.Domain, r *http.Request) (*owid.Node, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if d.Config.Debug {
		fmt.Println(d.Host)
		fmt.Println(string(b))
	}
	return owid.NodeFromJSON(b)
}

func sendToSupplier(
	d *common.Domain,
	s string,
	n *owid.Node) (*owid.Node, error) {

	// Turn the node into a byte array.
	j, err := n.GetRoot().AsJSON()
	if err != nil {
		return nil, err
	}

	// POST the bid to the supplier.
	var up url.URL
	up.Scheme = d.Config.Scheme
	up.Host = s
	up.Path = openRTBPath
	res, err := http.Post(
		up.String(),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return createFailed(d, n, &up, res)
	}

	// Read the response as a byte array.
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Convert the byte array to a tree to append as a child to the current
	// Processor's children
	c, err := owid.NodeFromJSON(b)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func createFailed(
	d *common.Domain,
	n *owid.Node,
	u *url.URL,
	res *http.Response) (*owid.Node, error) {
	var f swan.Failed
	f.Host = u.Host
	f.Error = fmt.Sprintf("%d", res.StatusCode)
	b, err := f.AsByteArray()
	if err != nil {
		return nil, err
	}
	r, err := n.GetRoot().GetOWID()
	if err != nil {
		return nil, err
	}
	oc, err := d.GetOWIDCreator()
	if err != nil {
		return nil, err
	}
	t, err := oc.CreateOWID(b)
	if err != nil {
		return nil, err
	}
	err = oc.Sign(t, r)
	var c owid.Node
	c.OWID, err = t.AsByteArray()
	if err != nil {
		return nil, err
	}
	return &c, nil
}
