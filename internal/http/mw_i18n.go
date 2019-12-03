package http

import (
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	acceptLanguage   = http.CanonicalHeaderKey("accept-language")
	i18nBundle       *i18n.Bundle
	i18nCtxKey       = ContextKey("i18n")
	i18nLocaleCtxKey = ContextKey("i18nLocale")
)

// I18n is a middleware that provides translations based on `Accept-Language` HTTP header.
func I18n(b *i18n.Bundle) HandlerFunc {
	i18nBundle = b

	return func(ctx *Context) {
		languages := strings.Split(ctx.Request.Header.Get(acceptLanguage), ",")
		localizer := i18n.NewLocalizer(i18nBundle, languages...)
		ctx.Set(i18nCtxKey.String(), localizer)

		if len(languages) > 0 {
			ctx.Set(i18nLocaleCtxKey.String(), languages[0])
		}

		ctx.Next()
	}
}

// I18nLocalizer returns the I18n localizer instance.
func I18nLocalizer(ctx *Context) *i18n.Localizer {
	localizer, exists := ctx.Get(i18nCtxKey.String())

	if !exists {
		return nil
	}

	return localizer.(*i18n.Localizer)
}

// I18nLocale returns the I18n locale.
func I18nLocale(ctx *Context) string {
	locale, exists := ctx.Get(i18nLocaleCtxKey.String())

	if locale == "" || !exists {
		return "en"
	}

	return locale.(string)
}

// I18nLocales returns all the available I18n locales.
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
	localizer := I18nLocalizer(ctx)

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

	if len(args) > 0 {
		locale := args[0]
		localizer = i18n.NewLocalizer(i18nBundle, locale)
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: data})
	if err != nil {
		return ""
	}

	return msg
}
