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
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
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
	templateFile string,
	data interface{}) error {

	err := canSend(s)
	if err != nil {
		fmt.Println("Sending failed:", err)
	}

	b, err := build(templateFile, subject, data)
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

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := s.Host
	smtpPort := s.Port

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func build(
	templateFile string,
	subject string,
	data interface{}) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Email Protection Reminder \n%s\n\n", mimeHeaders)))

	t.Execute(&body, data)

	return &body, nil
}
