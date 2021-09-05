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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type SMTP struct {
	Sender   string
	Host     string
	Port     string
	Password string
}

func NewSMTP() *SMTP {
	p := new(SMTP)

	p.Sender = os.Getenv("SMTP_SENDER")
	p.Host = os.Getenv("SMTP_HOST")
	p.Port = os.Getenv("SMTP_PORT")
	p.Password = os.Getenv("SMTP_PASSWORD")

	return p
}

func (s *SMTP) Send(
	email string,
	subject string,
	emailTemplate *template.Template,
	data interface{}) error {
	err := canSend(s)
	if err != nil {
		return err
	}

	b, err := build(emailTemplate, subject, data)
	if err != nil {
		return err
	}

	err = send(s, email, b)
	if err != nil {
		return err
	}

	return nil
}

func canSend(s *SMTP) error {
	if s.Sender == "" ||
		s.Host == "" ||
		s.Port == "" ||
		s.Password == "" {
		return errors.New(
			"cannot send email, make sure the following environment " +
				"variables are configured: SMTP_SENDER, SMTP_HOST, SMTP_PORT, " +
				"SMTP_PASSWORD")
	}
	return nil
}

func send(s *SMTP, email string, body *bytes.Buffer) error {
	// Sender data.
	from := s.Sender
	password := s.Password

	// smtp server configuration.
	smtpHost := s.Host
	smtpPort := s.Port

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Print(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Print(err)
	}

	// To && From
	if err = c.Mail(s.Sender); err != nil {
		log.Print(err)
	}

	if err = c.Rcpt(email); err != nil {
		log.Print(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Print(err)
	}

	_, err = w.Write(body.Bytes())
	if err != nil {
		log.Print(err)
	}

	err = w.Close()
	if err != nil {
		log.Print(err)
	}

	c.Quit()

	return nil
}

func build(
	emailTemplate *template.Template,
	subject string,
	data interface{}) (*bytes.Buffer, error) {

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Email Protection Reminder \n%s\n\n", mimeHeaders)))

	emailTemplate.Execute(&body, data)

	return &body, nil
}
