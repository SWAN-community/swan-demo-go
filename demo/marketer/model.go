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

package marketer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swan-demo-go/demo/shared"
	"github.com/SWAN-community/swan-go"
)

// MarketerModel used with HTML templates.
type MarketerModel struct {
	shared.PageModel
	idNode *owid.Node // The swan.ID as a node and tree associated with the page
}

// Version is a code for cache busting.
func (m *MarketerModel) Version() string {
	t := time.Now().UTC()
	h, _, _ := t.Clock()
	y, mon, d := t.Date()
	return fmt.Sprintf("%d%d%d%d", y, mon, d, h)
}

// Stop returns true if the request includes the key Stop to indicate that the
// advert should no longer be displayed.
func (m *MarketerModel) Stop() bool {
	for k := range m.Request.Form {
		if k == "stop" {
			return true
		}
	}
	return false
}

// TreeAsJSON return the transaction as JSON.
func (m *MarketerModel) TreeAsJSON() (template.HTML, error) {
	if m.idNode == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	b, err := m.idNode.AsJSON()
	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	html.WriteString("<p style=\"word-break:break-all\">")
	html.WriteString(string(b))
	html.WriteString("</p>")
	return template.HTML(html.String()), nil
}

// ID returns the SWAN ID string for the advert.
func (m *MarketerModel) ID() (string, error) {
	if m.idNode == nil {
		return "", nil
	}
	o, err := m.idNode.GetOWID()
	if err != nil {
		return "", err
	}
	return o.AsString(), nil
}

// IDUnpacked returns the unpacked SWAN ID.
func (m *MarketerModel) IDUnpacked() (template.HTML, error) {
	if m.idNode == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}
	o, err := m.idNode.GetOWID()
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}
	s, err := swan.IDFromOWID(o)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	var html bytes.Buffer
	html.WriteString("<table class=\"table\">")
	html.WriteString("<thead><tr>")
	html.WriteString("<th>Field</th>")
	html.WriteString("<th>Value</th>")
	html.WriteString("</tr></thead><tbody>")
	html.WriteString(fmt.Sprintf(
		"<tr><td>Version</td><td>%d</td></tr>",
		o.Version))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Domain</td><td>%s</td></tr>",
		o.Domain))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Created</td><td>%s</td></tr>",
		o.Date.Format("02-01-2006 15:04")))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Signature</td><td style=\"word-break:break-all\">%s</td></tr>",
		convertToString(o.Signature)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>SWID</td><td>%s<br/>%s<td></tr>",
		s.SWIDAsString(),
		owidTitle(s.SWID)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Preferences</td><td>%s<br/>%s<td></tr>",
		s.PreferencesAsString(),
		owidTitle(s.Preferences)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>SID</td><td>%s<br/>%s<td></tr>",
		s.SIDAsString(),
		owidTitle(s.SID)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Pub. domain</td><td>%s</td></tr>",
		s.PubDomain))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Unique</td><td style=\"word-break:break-all\">%s</td></tr>",
		convertToString(s.UUID)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Stopped Ads.</td><td style=\"word-break:break-all\">%s<td></tr>",
		strings.Join(s.StoppedAsArray(), ",")))
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

func owidTitle(o *owid.OWID) string {
	return fmt.Sprintf("%s %s", o.Date.Format("02-01-2006 15:04"), o.Domain)
}

// AuditWinnerHTML returns the audit information from the bid used in the advert
// that resulted in the request to this page.
func (m *MarketerModel) AuditWinnerHTML() (template.HTML, error) {
	var html bytes.Buffer
	if m.idNode == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	w, err := swan.WinningNode(m.idNode)
	if err != nil {
		return "", nil
	}

	if err != nil {
		return "", nil
	}
	htmlAddHeader(&html)
	err = appendParents(m.Domain, &html, w)
	if err != nil {
		return "", err
	}
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

// AuditFullHTML returns the audit information from the bid used in the advert
// that resulted in the request to this page.
func (m *MarketerModel) AuditFullHTML() (template.HTML, error) {
	if m.idNode == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	w, err := swan.WinningNode(m.idNode)
	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	htmlAddHeader(&html)
	err = appendOWIDAndChildren(m.Domain, &html, m.idNode, w, 0)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

func convertToString(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func htmlAddHeader(html *bytes.Buffer) {
	html.WriteString("<table class=\"table\">\r\n")
	html.WriteString("<thead>\r\n<tr>\r\n")
	html.WriteString("<th>Organization</th>\r\n")
	html.WriteString("<th>Audit Result</th>\r\n")
	html.WriteString("<th>\r\n</th>\r\n")
	html.WriteString("<th>\r\n</th>\r\n")
	html.WriteString("<th>\r\n</th>\r\n")
	html.WriteString("</tr>\r\n</thead>\r\n<tbody>\r\n")
}

func htmlAddFooter(html *bytes.Buffer) {
	html.WriteString("</tbody>\r\n</table>\r\n")
}

func appendParents(d *shared.Domain, html *bytes.Buffer, w *owid.Node) error {
	var n []*owid.Node
	p := w
	for p != nil {
		n = append(n, p)
		p = p.GetParent()
	}
	i := len(n) - 1
	for i >= 0 {
		err := appendHTML(d, html, w, n[i], 0)
		if err != nil {
			return err
		}
		i--
	}
	return nil
}

func appendOWIDAndChildren(
	d *shared.Domain,
	html *bytes.Buffer,
	o *owid.Node,
	w *owid.Node,
	level int) error {
	appendHTML(d, html, w, o, level)
	if len(o.Children) > 0 {
		for _, c := range o.Children {
			err := appendOWIDAndChildren(d, html, c, w, level+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func appendHTML(
	d *shared.Domain,
	html *bytes.Buffer,
	w *owid.Node,
	o *owid.Node,
	level int) error {

	s, err := swan.FromNode(o)
	if err != nil {
		return err
	}

	html.WriteString("<tr>\r\n")
	html.WriteString(fmt.Sprintf(
		"<td style=\"padding-left:%dem;\" class=\"text-left\">\r\n"+
			"<script>appendName(document.currentScript.parentNode,\"%s\")</script></td>\r\n",
		level,
		o.GetOWIDAsString()))

	var r string
	if o.GetRoot() != o {
		r = o.GetRoot().GetOWIDAsString()
	}

	html.WriteString(fmt.Sprintf(
		"<td style=\"text-align:center;\">\r\n"+
			"<script>appendAuditMark(document.currentScript.parentNode,\"%s\",\"%s\");</script>\r\n"+
			"<noscript>JavaScript needed to audit</noscript></td>\r\n",
		r,
		o.GetOWIDAsString()))
	if w == o {
		html.WriteString("<td>\r\n<img style=\"width:32px\" src=\"noun_rosette_470370.svg\">\r\n</td>\r\n")
	} else {
		f, fok := s.(*swan.Failed)
		_, bok := s.(*swan.Bid)
		if fok {
			html.WriteString(fmt.Sprintf("<td style=\"color:lightpink\">\r\n%s&nbsp;%s</td>\r\n",
				f.Host,
				f.Error))
		} else if bok {
			html.WriteString("<td>\r\n<img style=\"width:32px\" src=\"noun_movie ticket_1807397.svg\">\r\n</td>\r\n")
		} else {
			html.WriteString("<td>\r\n</td>\r\n")
		}
	}

	if o != nil && r != "" {
		html.WriteString(fmt.Sprintf(
			"<td style=\"text-align:center;\">\r\n"+
				"<script>appendComplaintEmail(document.currentScript.parentNode,\"%s\",\"%s\",\"%s\", \"noun_complaint_376466.svg\");</script>\r\n"+
				"<noscript>JavaScript needed to audit</noscript></td>\r\n",
			d.CMP,
			r,
			o.GetOWIDAsString()))
	} else {
		html.WriteString("<td>\r\n</td>\r\n")
	}

	html.WriteString(fmt.Sprintf(
		"<td style=\"text-align:center;\">\r\n"+
			"<style>img { width: 32px; }</style><script>appendContractURL(document.currentScript.parentNode,\"%s\")</script></td>\r\n",
		o.GetOWIDAsString()))

	html.WriteString("</tr>\r\n")
	return nil
}
