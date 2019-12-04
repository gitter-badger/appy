package mailer

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"net/textproto"
	"regexp"
	"strings"

	appyhttp "github.com/appist/appy/internal/http"
	appysupport "github.com/appist/appy/internal/support"
	"github.com/jordan-wright/email"
)

type (
	// Mailer provides the capability to parse/render email template and send it out via SMTP protocol.
	Mailer struct {
		addr      string
		plainAuth smtp.Auth
		previews  map[string]Mail
		server    *appyhttp.Server
	}

	// Mail defines the email headers/body/attachments.
	Mail struct {
		From         string
		To           []string
		ReplyTo      []string
		Bcc          []string
		Cc           []string
		Sender       string
		Subject      string
		Headers      textproto.MIMEHeader
		Template     string
		TemplateData interface{}
		Attachments  []string
		ReadReceipt  []string
	}
)

// NewMailer initializes Mailer instance.
func NewMailer(c *appysupport.Config, l *appysupport.Logger, s *appyhttp.Server) *Mailer {
	mailer := &Mailer{
		addr: c.MailerAddr,
		plainAuth: smtp.PlainAuth(
			c.MailerPlainAuthIdentity,
			c.MailerPlainAuthUsername,
			c.MailerPlainAuthPassword,
			c.MailerPlainAuthHost,
		),
		previews: map[string]Mail{},
		server:   s,
	}

	if appysupport.IsDebugBuild() {
		s.HTMLRenderer().AddFromString("mailer/preview", previewTpl())

		s.Router().GET(s.Config().MailerPreviewBaseURL, func(ctx *appyhttp.Context) {
			name := ctx.DefaultQuery("name", "")
			if name == "" && len(mailer.previews) > 0 {
				for _, preview := range mailer.previews {
					name = preview.Template
					break
				}
			}

			locales := appyhttp.I18nLocales()
			locale := ctx.DefaultQuery("locale", "")
			if locale == "" {
				if len(locales) > 0 {
					locale = locales[0]
				}
			}

			ctx.HTML(http.StatusOK, "mailer/preview", appyhttp.H{
				"baseURL":  s.Config().MailerPreviewBaseURL,
				"previews": mailer.previews,
				"title":    "Mailer Preview",
				"name":     name,
				"ext":      ctx.DefaultQuery("ext", "html"),
				"locale":   locale,
				"locales":  locales,
				"mail":     mailer.previews[name],
			})
		})

		s.Router().GET(s.Config().MailerPreviewBaseURL+"/preview", func(ctx *appyhttp.Context) {
			name := ctx.Query("name")
			preview, exists := mailer.previews[name]
			if !exists {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}

			locale := ctx.Query("locale")
			appyhttp.SetI18nLocale(ctx, locale)
			preview.TemplateData.(appyhttp.H)["t"] = func(key string, count int, str string, args ...string) string {
				data := make(map[string]interface{})
				_ = json.Unmarshal([]byte(str), &data)

				return appyhttp.T(ctx, key, count, data, args...)
			}

			var (
				contentType string
				content     []byte
				err         error
			)
			switch ctx.Query("ext") {
			case "html":
				contentType = "text/html"
				content, err = mailer.Content("html", name, preview.TemplateData)
			case "txt":
				contentType = "text/plain"
				content, err = mailer.Content("txt", name, preview.TemplateData)
			}

			if err != nil {
				panic(err)
			}

			ctx.Writer.Header().Del(http.CanonicalHeaderKey("x-frame-options"))
			ctx.Data(http.StatusOK, contentType, content)
		})
	}

	return mailer
}

// Preview sets up the preview for the mail HTML/text template.
func (m *Mailer) Preview(mail Mail) {
	m.previews[mail.Template] = mail
}

// Send delivers the email via SMTP protocol without TLS.
func (m Mailer) Send(mail Mail) error {
	email, err := m.composeEmail(mail)
	if err != nil {
		return err
	}

	return email.Send(m.addr, m.plainAuth)
}

// SendWithTLS delivers the email via SMTP protocol with TLS.
func (m Mailer) SendWithTLS(mail Mail, tls *tls.Config) error {
	email, err := m.composeEmail(mail)
	if err != nil {
		return err
	}

	return email.SendWithTLS(m.addr, m.plainAuth, tls)
}

// Content returns the content for the named email template.
func (m Mailer) Content(ext, name string, obj interface{}) ([]byte, error) {
	var objCopy interface{}
	if err := appysupport.DeepClone(&objCopy, &obj); err != nil {
		return nil, err
	}

	if _, ok := obj.(appyhttp.H)["layout"]; !ok {
		objCopy.(appyhttp.H)["layout"] = "mailer_default"
	}

	objCopy.(appyhttp.H)["layout"] = strings.TrimSuffix((objCopy.(appyhttp.H)["layout"]).(string), ".html")
	objCopy.(appyhttp.H)["layout"] = strings.TrimSuffix((objCopy.(appyhttp.H)["layout"]).(string), ".txt")
	objCopy.(appyhttp.H)["layout"] = (objCopy.(appyhttp.H)["layout"]).(string) + "." + ext

	recorder := httptest.NewRecorder()
	renderer := m.server.HTMLRenderer().Instance(name+"."+ext, objCopy)

	if err := renderer.Render(recorder); err != nil {
		return nil, err
	}

	return recorder.Body.Bytes(), nil
}

func (m Mailer) composeEmail(mail Mail) (*email.Email, error) {
	email := &email.Email{
		From:        mail.From,
		To:          mail.To,
		ReplyTo:     mail.ReplyTo,
		Bcc:         mail.Bcc,
		Cc:          mail.Cc,
		Sender:      mail.Sender,
		Subject:     mail.Subject,
		ReadReceipt: mail.ReadReceipt,
	}

	if mail.Headers == nil {
		email.Headers = textproto.MIMEHeader{}
	}

	tpl := mail.Template
	if regexp.MustCompile(`\.html$`).Match([]byte(tpl)) {
		tpl = strings.TrimSuffix(tpl, ".html")
	} else if regexp.MustCompile(`\.txt$`).Match([]byte(tpl)) {
		tpl = strings.TrimSuffix(tpl, ".txt")
	}

	// TODO: Add locale handling
	html, err := m.Content("html", mail.Template, mail.TemplateData)
	if err != nil {
		return nil, err
	}
	email.HTML = html

	text, err := m.Content("txt", mail.Template, mail.TemplateData)
	if err != nil {
		return nil, err
	}
	email.Text = text

	if mail.Attachments != nil {
		for _, attachment := range mail.Attachments {
			// TODO: Add external/internal file handling.
			email.AttachFile(attachment)
		}
	}

	return email, nil
}
