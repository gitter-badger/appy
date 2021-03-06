//+build integration

package appy_test

import (
	"testing"

	"github.com/appist/appy"
)

type AppSuite struct {
	appy.TestSuite
}

func (s *AppSuite) TestNewApp() {
	asset := appy.NewAsset(nil, map[string]string{
		"docker": "testdata/app/.docker",
		"config": "testdata/app/configs",
		"locale": "testdata/app/pkg/locales",
		"view":   "testdata/app/pkg/views",
		"web":    "testdata/app/web",
	})
	app := appy.NewApp(asset, nil)

	s.NotNil(app.Asset())
	s.NotNil(app.Command())
	s.NotNil(app.Config())
	s.NotNil(app.DBManager())
	s.NotNil(app.I18n())
	s.NotNil(app.Logger())
	s.NotNil(app.Mailer())
	s.NotNil(app.Server())
	s.NotNil(app.Support())
	s.NotNil(app.ViewEngine())
	s.NoError(app.Run())
}

func TestAppSuite(t *testing.T) {
	appy.RunTestSuite(t, new(AppSuite))
}
