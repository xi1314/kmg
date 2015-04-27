package kmgEmail

import (
	"bytes"
	"crypto/tls"
	"github.com/bronze1man/InlProxy"
	"github.com/bronze1man/kmg/kmgNet"
	"net"
	"net/smtp"
	"strings"
	"text/template"
)

func SmtpSendEmail(req *SmtpRequest) (err error) {
	parameters := &struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		req.From,
		strings.Join([]string(req.To), ","),
		req.Subject,
		req.Message,
	}

	buffer := new(bytes.Buffer)

	t := template.Must(template.New("emailTemplate").Parse(_EmailScript()))
	t.Execute(buffer, parameters)

	auth := smtp.PlainAuth("", req.From, req.SmtpPassword, req.SmtpHost)

	var conn net.Conn
	if req.Socks4aProxyAddr != "" {
		conn, err = InlProxy.Socks4aDial(req.Socks4aProxyAddr, kmgNet.JoinHostPortInt(req.SmtpHost, req.SmtpPort))
	} else {
		conn, err = net.Dial("tcp", kmgNet.JoinHostPortInt(req.SmtpHost, req.SmtpPort))
	}
	if err != nil {
		return
	}
	c, err := smtp.NewClient(conn, req.SmtpHost)
	if err != nil {
		return
	}
	err = smtpSendMailPart2(c,
		req.SmtpHost,
		auth,
		req.From,
		req.To,
		buffer.Bytes())
	return err
}

func smtpSendMailPart2(c *smtp.Client, smtpHost string, a smtp.Auth, from string, to []string, msg []byte) (err error) {
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: smtpHost}
		if err = c.StartTLS(config); err != nil {
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

// _EmailScript returns a template for the email message to be sent
func _EmailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}
