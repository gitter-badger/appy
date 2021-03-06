package appy

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

var (
	crawlerCtxKey   = ContextKey("crawler")
	staticExtRegex  = regexp.MustCompile(`\.(bmp|css|csv|eot|exif|gif|html|ico|ini|jpg|jpeg|js|json|mp4|otf|pdf|png|svg|webp|woff|woff2|tiff|ttf|toml|txt|xml|xlsx|yml|yaml)$`)
	userAgentHeader = http.CanonicalHeaderKey("user-agent")
	xPrerender      = http.CanonicalHeaderKey("x-prerender")
)

// Prerender dynamically renders client-side rendered SPA for SEO using Chrome.
func Prerender(config *Config, logger *Logger) HandlerFunc {
	scheme := "http"
	host := config.HTTPHost
	port := config.HTTPPort

	if config.HTTPSSLEnabled {
		scheme = "https"
		port = config.HTTPSSLPort
	}

	return func(c *Context) {
		request := c.Request
		userAgent := request.Header.Get(userAgentHeader)

		if !staticExtRegex.MatchString(request.URL.Path) && isSEOBot(userAgent) {
			url := fmt.Sprintf("%s://%s:%s%s", scheme, host, port, request.URL)
			logger.Infof("[HTTP][PRERENDER] SEO bot \"%s\" crawling \"%s\"...", userAgent, url)

			crawler, exists := c.Get(crawlerCtxKey.String())
			if !exists {
				crawler = &Crawl{}
			}

			data, err := crawler.(Crawler).Perform(url)
			if err != nil {
				logger.Error(err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			c.Writer.Header().Add(xPrerender, "1")
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			c.Abort()
			return
		}

		c.Next()
	}
}

func isSEOBot(ua string) bool {
	bots := []string{
		"googlebot", "yahoou", "bingbot", "baiduspider", "yandex", "yeti", "yodaobot", "gigabot", "ia_archiver",
		"facebookexternalhit", "twitterbot", "developers.google.com",
	}

	for _, bot := range bots {
		if strings.Contains(strings.ToLower(ua), bot) {
			return true
		}
	}

	return false
}

type (
	// Crawler satisfies Crawl type and implements all its functions, mainly used for mocking in unit test.
	Crawler interface {
		Perform(url string) ([]byte, error)
	}

	// Crawl is used to manipulate DOM by using Chromium via chromedp.
	Crawl struct {
	}
)

// Perform manipulates the URL's DOM by using Chromium via chromedp.
func (c *Crawl) Perform(url string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Headless,
		chromedp.Flag("ignore-certificate-errors", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var data string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}

			data, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
