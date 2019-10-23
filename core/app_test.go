package core

import (
	"os"
	"testing"

	"github.com/appist/appy/test"
)

type AppSuite struct {
	test.Suite
	oldConfigPath string
}

func (s *AppSuite) SetupTest() {
}

func (s *AppSuite) TearDownTest() {
	os.Unsetenv("APPY_ENV")
	os.Unsetenv("APPY_MASTER_KEY")
	os.Unsetenv("HTTP_CSRF_SECRET")
	os.Unsetenv("HTTP_SESSION_SECRETS")
}

func (s *AppSuite) TestNewApp() {
	os.Setenv("HTTP_CSRF_SECRET", "58f364f29b568807ab9cffa22c99b538")
	os.Setenv("HTTP_SESSION_SECRETS", "58f364f29b568807ab9cffa22c99b538")
	oldConfigPath := SSRPaths["config"]
	SSRPaths["config"] = "./testdata/.ssr/app/config"

	app, err := NewApp(nil, nil, nil)
	s.Nil(err)
	s.NotNil(app.Config)
	s.NotNil(app.Logger)
	s.NotNil(app.Server)
	SSRPaths["config"] = oldConfigPath
}

func (s *AppSuite) TestNewAppWithMissingRequiredEnvVariables() {
	os.Setenv("APPY_MASTER_KEY", "58f364f29b568807ab9cffa22c99b538")
	os.Args = append(os.Args, "serve")
	_, err := NewApp(nil, nil, nil)
	s.Contains(err.Error(), "required environment variable \"HTTP_SESSION_SECRETS\" is not set. required environment variable \"HTTP_CSRF_SECRET\" is not set")
	os.Args = os.Args[:len(os.Args)-1]
}

func TestApp(t *testing.T) {
	test.Run(t, new(AppSuite))
}
