package http

import (
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	acceptLanguage   = http.CanonicalHeaderKey("accept-language")
	i18nBundle       *i18n.Bundle
	i18nLocaleCtxKey = ContextKey("i18nLocale")
)

// I18n is a middleware that provides translations based on `Accept-Language` HTTP header.
func I18n(b *i18n.Bundle) HandlerFunc {
	i18nBundle = b

	return func(ctx *Context) {
		languages := strings.Split(ctx.Request.Header.Get(acceptLanguage), ",")

		if len(languages) > 0 {
			ctx.Set(i18nLocaleCtxKey.String(), languages[0])
		}

		ctx.Next()
	}
}

// I18nLocale returns the current context's locale.
func I18nLocale(ctx *Context) string {
	locale, exists := ctx.Get(i18nLocaleCtxKey.String())

	if locale == "" || !exists {
		return "en"
	}

	return locale.(string)
}

// SetI18nLocale sets the current context's locale.
func SetI18nLocale(ctx *Context, locale string) {
	ctx.Set(i18nLocaleCtxKey.String(), locale)
}

// I18nLocales returns all the available locales.
func I18nLocales() []string {
	locales := []string{}

	if i18nBundle != nil {
		for _, tag := range i18nBundle.LanguageTags() {
			locales = append(locales, tag.String())
		}
	}

	return locales
}

// T translates a message based on the given key which count is used to pluralise the translation if needed.
func T(ctx *Context, key string, count int, data map[string]interface{}, args ...string) string {
	if count != -1 {
		switch count {
		case 0:
			key = key + ".Zero"
		case 1:
			key = key + ".One"
		default:
			key = key + ".Other"
		}

		data["Count"] = count
	}

	locale := I18nLocale(ctx)
	if len(args) > 0 {
		locale = args[0]
	}

	localizer := i18n.NewLocalizer(i18nBundle, locale)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: data})
	if err != nil {
		return ""
	}

	return msg
}
