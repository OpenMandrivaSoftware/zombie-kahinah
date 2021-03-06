package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"gitea.tsn.sh/robert/zombie-kahinah/models"
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/xiam/to"
)

var (
	ErrDisabled = errors.New("Mail Service is Disabled")

	outwardUrl    = beego.AppConfig.String("outwardloc")
	outwardPrefix = beego.AppConfig.String("urlprefix")

	maintainer_hours = to.Int64(beego.AppConfig.String("karma::maintainerhours"))

	mail_enabled = to.Bool(beego.AppConfig.String("mail::enabled"))
	mail_user    = beego.AppConfig.String("mail::smtp_user")
	mail_pass    = beego.AppConfig.String("mail::smtp_pass")
	mail_domain  = beego.AppConfig.String("mail::smtp_domain")
	mail_host    = beego.AppConfig.String("mail::smtp_host")
	mail_verify  = beego.AppConfig.String("mail::smtp_tls_verify")
	mail_email   = mail.Address{"Kahinah QA Bot", beego.AppConfig.String("mail::smtp_email")}

	mail_to    = beego.AppConfig.String("mail::to")
	mail_maint = to.Bool(beego.AppConfig.String("mail::maint"))

	model_template       = "emails/regular.tpl"
	maint_model_template = "emails/maint.tpl"
	digest_template      = "emails/digest.tpl"
	mail_template        = template.New("email full template")

	digestModelsIntake = make(chan *models.BuildList, 100)
	digestModelsQueue  = []*models.BuildList{}
)

func init() {
	mail_template = template.Must(mail_template.Parse(`From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Mime-Version: 1.0
Content-type: text/plain

{{.Body}}

-------------------------------
This email was sent by Kahinah, the OpenMandriva QA bot.
Inbound email to this account is not monitored.
`))

	// digest model queue
	go func() {
		timeWait := time.After(12 * time.Hour)
		for {
			select {
			case <-timeWait:
				MailDigest()
				timeWait = time.After(12 * time.Hour)
			case in := <-digestModelsIntake:
				digestModelsQueue = append(digestModelsQueue, in)
			}
		}
	}()
}

func MailDigest() {
	if len(digestModelsQueue) == 0 {
		// nothing to do
		return
	}

	defer func() {
		digestModelsQueue = []*models.BuildList{}
	}()

	if !mail_enabled || mail_to == "" {
		return // Disabled
	}

	// just in case
	o := orm.NewOrm()

	// make sure we have information
	for _, v := range digestModelsQueue {
		if v.Packages == nil {
			o.LoadRelated(v, "Packages")
		}
	}

	kahinahURL := outwardUrl
	if outwardPrefix != "" {
		kahinahURL = outwardUrl + "/" + outwardPrefix
	}

	data := map[string]interface{}{
		"kahinahURL": kahinahURL,
		"Lists":      digestModelsQueue,
	}

	var digestBuf bytes.Buffer
	if err := beego.ExecuteTemplate(&digestBuf, digest_template, data); err != nil {
		log.Printf("[mail] digest template failed: %v\n", err)
		return
	}

	subject := fmt.Sprintf("[kahinah] 12-hour digest at %v", time.Now())

	if err := MailTo(subject, digestBuf.String(), mail_to); err != nil {
		log.Printf("[mail] digest email failed: %v\n", err)
	}
}

func MailTo(subject, content, to string) error {
	if !mail_enabled {
		return ErrDisabled
	}

	data := make(map[string]string)
	data["From"] = mail_email.String()
	data["To"] = to
	data["Subject"] = subject
	data["Body"] = content

	var buf bytes.Buffer
	mail_template.Execute(&buf, data)

	if mail_domain == "" {
		if strings.Contains(mail_user, "@") {
			mail_domain = mail_user[strings.Index(mail_user, "@")+1:]
		} else {
			mail_domain = mail_host[:strings.Index(mail_host, ":")]
		}
	}

	return ourMail(mail_host, smtp.PlainAuth("", mail_user, mail_pass, mail_domain), mail_email.Address, []string{to}, buf.Bytes())

}

func Mail(subject, content string) error {
	return MailTo(subject, content, mail_to)
}

// this function:
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func ourMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	//if err = c.Hello(); err != nil {
	//	return err
	//}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(&tls.Config{
			InsecureSkipVerify: mail_verify == "",
			ServerName:         mail_verify,
		}); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func MailModel(model *models.BuildList) {
	defer func() {
		digestModelsIntake <- model
	}()

	if model.Submitter == nil {
		o := orm.NewOrm()
		o.LoadRelated(model, "Submitter")
	}

	if model.Karma == nil {
		o := orm.NewOrm()
		o.LoadRelated(model, "Karma")
		for _, karma := range model.Karma {
			o.LoadRelated(karma, "User")
		}
	}

	data := make(map[string]interface{})

	action := "lost in an abyss"

	switch model.Status {
	case models.STATUS_TESTING:
		action = "pushed to testing"
	case models.STATUS_PUBLISHED:
		action = "published"
	case models.STATUS_REJECTED:
		action = "rejected/cleared"
	}

	data["Action"] = action

	data["KahinahUrl"] = outwardUrl
	if outwardPrefix != "" {
		data["KahinahUrl"] = outwardUrl + "/" + outwardPrefix
	}
	data["Package"] = model

	// var modelTemplateBuf bytes.Buffer
	// beego.ExecuteTemplate(&modelTemplateBuf, model_template, data)

	if mail_maint {
		var maintTemplateBuf bytes.Buffer
		if err := beego.ExecuteTemplate(&maintTemplateBuf, maint_model_template, data); err != nil {
			log.Printf("[mail] maint template failed: %v\n", err)
			return
		}

		subject := fmt.Sprintf("[kahinah] %v-%v (%v) %v %v", model.Name, model.SourceEVR(), model.Architecture, model.Id, action)

		// err := Mail(subject, modelTemplateBuf.String())
		// if err != nil {
		// 	log.Printf("[mail] model email failed: %s\n", err)
		// }

		err := MailTo(subject, maintTemplateBuf.String(), model.Submitter.Email)
		if err != nil {
			log.Printf("[mail] maint email failed: %v\n", err)
		}
	}
}
