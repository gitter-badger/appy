package appy

import (
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type ConfigSuite struct {
	TestSuite
	oldSSRPaths map[string]string
}

func (s *ConfigSuite) SetupTest() {
	s.oldSSRPaths = _ssrPaths
	_ssrPaths = map[string]string{
		"root":   "testdata/.ssr",
		"config": "testdata/pkg/config",
		"locale": "testdata/pkg/locales",
		"view":   "testdata/pkg/views",
	}
}

func (s *ConfigSuite) TearDownTest() {
	_ssrPaths = s.oldSSRPaths
	os.Clearenv()
}

func (s *ConfigSuite) TestNewConfigDefaultValue() {
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	tt := map[string]interface{}{
		"AppyEnv":                         "development",
		"GQLPlaygroundEnabled":            false,
		"GQLPlaygroundPath":               "/docs/graphql",
		"GQLCacheSize":                    1000,
		"GQLComplexityLimit":              200,
		"GQLUploadMaxMemory":              int64(100000000),
		"GQLUploadMaxSize":                int64(100000000),
		"GQLWebsocketKeepAliveDuration":   25 * time.Second,
		"HTTPDebugEnabled":                false,
		"HTTPLogFilterParameters":         []string{"password"},
		"HTTPHealthCheckURL":              "/health_check",
		"HTTPHost":                        "localhost",
		"HTTPPort":                        "3000",
		"HTTPGracefulTimeout":             30 * time.Second,
		"HTTPIdleTimeout":                 75 * time.Second,
		"HTTPMaxHeaderBytes":              0,
		"HTTPReadTimeout":                 60 * time.Second,
		"HTTPReadHeaderTimeout":           60 * time.Second,
		"HTTPWriteTimeout":                60 * time.Second,
		"HTTPSSLEnabled":                  false,
		"HTTPSSLCertPath":                 "./tmp/ssl",
		"HTTPSessionCookieDomain":         "localhost",
		"HTTPSessionCookieHTTPOnly":       true,
		"HTTPSessionCookieMaxAge":         0,
		"HTTPSessionCookiePath":           "/",
		"HTTPSessionCookieSecure":         false,
		"HTTPSessionRedisAddr":            "localhost:6379",
		"HTTPSessionRedisAuth":            "",
		"HTTPSessionRedisDb":              "0",
		"HTTPSessionRedisMaxActive":       0,
		"HTTPSessionRedisMaxIdle":         32,
		"HTTPSessionRedisIdleTimeout":     30 * time.Second,
		"HTTPSessionRedisMaxConnLifetime": 0 * time.Second,
		"HTTPSessionRedisWait":            true,
		"HTTPSessionName":                 "_session",
		"HTTPSessionProvider":             "cookie",
		"HTTPSessionSecrets":              [][]byte{},
		"HTTPAllowedHosts":                []string{},
		"HTTPCSRFCookieDomain":            "localhost",
		"HTTPCSRFCookieHTTPOnly":          true,
		"HTTPCSRFCookieMaxAge":            0,
		"HTTPCSRFCookieName":              "_csrf_token",
		"HTTPCSRFCookiePath":              "/",
		"HTTPCSRFCookieSecure":            false,
		"HTTPCSRFFieldName":               "authenticity_token",
		"HTTPCSRFRequestHeader":           "X-CSRF-Token",
		"HTTPCSRFSecret":                  []byte{},
		"HTTPSSLRedirect":                 false,
		"HTTPSSLTemporaryRedirect":        false,
		"HTTPSSLHost":                     "localhost:3443",
		"HTTPSTSSeconds":                  int64(0),
		"HTTPSTSIncludeSubdomains":        false,
		"HTTPFrameDeny":                   true,
		"HTTPCustomFrameOptionsValue":     "",
		"HTTPContentTypeNosniff":          false,
		"HTTPBrowserXSSFilter":            false,
		"HTTPContentSecurityPolicy":       "",
		"HTTPReferrerPolicy":              "",
		"HTTPIENoOpen":                    false,
		"HTTPSSLProxyHeaders":             map[string]string{},
	}

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	cv := reflect.ValueOf(*config)
	for key, defaultVal := range tt {
		fv := cv.FieldByName(key)

		// An exception case to handle a different host in test for Github actions.
		if key == "HTTPSessionRedisAddr" && os.Getenv("HTTP_SESSION_REDIS_ADDR") != "" {
			s.Equal(fv.Interface(), os.Getenv("HTTP_SESSION_REDIS_ADDR"))
			continue
		}

		switch fv.Kind() {
		case reflect.Map:
			switch fv.Type().String() {
			case "map[string]string":
				for key, val := range fv.Interface().(map[string]string) {
					s.Equal(val, defaultVal.(map[string]string)[key])
				}
			}
		case reflect.Slice, reflect.Array:
			switch fv.Type().String() {
			case "[]string":
				s.Equal(len(fv.Interface().([]string)), len(defaultVal.([]string)))

				for idx, val := range fv.Interface().([]string) {
					s.Equal(val, defaultVal.([]string)[idx])
				}
			case "[]uint8":
				s.Equal(len(fv.Interface().([]uint8)), len(defaultVal.([]uint8)))

				for idx, val := range fv.Interface().([]uint8) {
					s.Equal(val, defaultVal.([]uint8)[idx])
				}
			case "[][]uint8":
				s.Equal(len(fv.Interface().([][]uint8)), len(defaultVal.([][]uint8)))

				for idx, val := range fv.Interface().([][]uint8) {
					s.Equal(val, defaultVal.([][]uint8)[idx])
				}
			default:
				s.Equal(fv.Interface(), defaultVal)
			}
		default:
			s.Equal(fv.Interface(), defaultVal)
		}
	}
}

func (s *ConfigSuite) TestNewConfigWithoutSettingRequiredConfig() {
	os.Setenv("APPY_ENV", "development")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.NotNil(config.Errors())
	s.EqualError(config.Errors()[0], `required environment variable "HTTP_SESSION_SECRETS" is not set. required environment variable "HTTP_CSRF_SECRET" is not set`)
}

func (s *ConfigSuite) TestNewConfigWithSettingRequiredConfig() {
	os.Setenv("APPY_ENV", "development")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_CSRF_SECRET", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_SESSION_SECRETS", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Equal([]byte("481e5d98a31585148b8b1dfb6a3c0465"), config.MasterKey())
	s.Nil(config.Errors())
}

func (s *ConfigSuite) TestNewConfigWithUnparsableEnvVariable() {
	os.Setenv("APPY_ENV", "development")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_DEBUG_ENABLED", "nil")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Contains(config.Errors()[0].Error(), `strconv.ParseBool: parsing "nil": invalid syntax.`)
}

func (s *ConfigSuite) TestNewConfigWithUndecryptableConfig() {
	os.Setenv("APPY_ENV", "undecryptable")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Contains(config.Errors()[0].Error(), "unable to decrypt 'HTTP_PORT' value in 'testdata/pkg/config/.env.undecryptable'")
}

func (s *ConfigSuite) TestNewConfigWithInvalidAssetsPath() {
	os.Setenv("APPY_ENV", "development")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	build := ReleaseBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, http.Dir("testdata"))
	s.Contains(config.Errors()[0].Error(), "open testdata/testdata/.ssr/testdata/pkg/config/.env.development: no such file or directory")
}

func (s *ConfigSuite) TestNewConfigWithMissingConfigInAssets() {
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	build := ReleaseBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.EqualError(config.Errors()[0], ErrNoConfigInAssets.Error())
}

func (s *ConfigSuite) TestNewConfigWithUnparsableConfig() {
	os.Setenv("APPY_ENV", "unparsable")
	os.Setenv("APPY_MASTER_KEY", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Contains(config.Errors()[0].Error(), "Can't separate key from value")
}

func (s *ConfigSuite) TestNewConfigWithInvalidDatabaseConfig() {
	os.Setenv("APPY_ENV", "development")
	os.Setenv("APPY_MASTER_KEY", "58f364f29b568807ab9cffa22c99b538")
	os.Setenv("DB_ADDR_PRIMARY", "dummy")
	os.Setenv("DB_APP_NAME_PRIMARY", "dummy")
	os.Setenv("DB_DIAL_TIMEOUT_PRIMARY", "dummy")
	os.Setenv("DB_IDLE_CHECK_FREQUENCY_PRIMARY", "dummy")
	os.Setenv("DB_IDLE_TIMEOUT_PRIMARY", "dummy")
	os.Setenv("DB_MAX_CONN_AGE_PRIMARY", "dummy")
	os.Setenv("DB_MAX_RETRIES_PRIMARY", "dummy")
	os.Setenv("DB_MIN_IDLE_CONNS_PRIMARY", "dummy")
	os.Setenv("DB_NAME_PRIMARY", "dummy")
	os.Setenv("DB_PASSWORD_PRIMARY", "dummy")
	os.Setenv("DB_POOL_SIZE_PRIMARY", "dummy")
	os.Setenv("DB_POOL_TIMEOUT_PRIMARY", "dummy")
	os.Setenv("DB_READ_TIMEOUT_PRIMARY", "dummy")
	os.Setenv("DB_REPLICA_PRIMARY", "dummy")
	os.Setenv("DB_RETRY_STATEMENT_PRIMARY", "dummy")
	os.Setenv("DB_SCHEMA_SEARCH_PATH_PRIMARY", "dummy")
	os.Setenv("DB_SCHEMA_MIGRATIONS_TABLE_PRIMARY", "true")
	os.Setenv("DB_USER_PRIMARY", "dummy")
	os.Setenv("DB_WRITE_TIMEOUT_PRIMARY", "dummy")
	os.Setenv("HTTP_CSRF_SECRET", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_SESSION_SECRETS", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	dbManager := NewDbManager(logger)
	s.Nil(config.Errors())
	s.NotNil(dbManager.Errors())
}

func (s *ConfigSuite) TestNewConfigWithValidDatabaseConfig() {
	os.Setenv("APPY_ENV", "valid_db")
	os.Setenv("APPY_MASTER_KEY", "58f364f29b568807ab9cffa22c99b538")
	os.Setenv("HTTP_CSRF_SECRET", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_SESSION_SECRETS", "481e5d98a31585148b8b1dfb6a3c0465")

	// A workaround for Github Action
	_, err := net.Dial("tcp", "0.0.0.0:5432")
	if err != nil {
		os.Setenv("DB_ADDR_PRIMARY", "localhost:32768")
	}

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	dbManager := NewDbManager(logger)
	s.Nil(config.Errors())
	s.Nil(dbManager.Errors())
	s.Nil(dbManager.ConnectAll(true))
}

func (s *ConfigSuite) TestIsConfigErrored() {
	os.Setenv("APPY_ENV", "development")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Equal(true, IsConfigErrored(config, logger))

	os.Setenv("APPY_MASTER_KEY", "58f364f29b568807ab9cffa22c99b538")
	os.Setenv("HTTP_CSRF_SECRET", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_SESSION_SECRETS", "481e5d98a31585148b8b1dfb6a3c0465")

	logger = NewLogger(build)
	config = NewConfig(build, logger, nil)
	s.Equal(false, IsConfigErrored(config, logger))
}

func (s *ConfigSuite) TestIsProtectedEnv() {
	os.Setenv("APPY_MASTER_KEY", "58f364f29b568807ab9cffa22c99b538")
	os.Setenv("HTTP_CSRF_SECRET", "481e5d98a31585148b8b1dfb6a3c0465")
	os.Setenv("HTTP_SESSION_SECRETS", "481e5d98a31585148b8b1dfb6a3c0465")

	build := DebugBuild
	logger := NewLogger(build)
	config := NewConfig(build, logger, nil)
	s.Equal(false, IsProtectedEnv(config))

	os.Setenv("APPY_ENV", "production")
	logger = NewLogger(build)
	config = NewConfig(build, logger, nil)
	s.Equal(true, IsProtectedEnv(config))
}

func (s *ConfigSuite) TestMasterKeyWithMissingKeyFile() {
	_, err := parseMasterKey()
	s.EqualError(err, ErrReadMasterKeyFile.Error())
}

func (s *ConfigSuite) TestMasterKeyWithMissingAppyMasterKey() {
	Build = ReleaseBuild
	_, err := parseMasterKey()
	s.EqualError(err, ErrNoMasterKey.Error())
	Build = DebugBuild
}

func (s *ConfigSuite) TestMasterKeyWithZeroLength() {
	Build = ReleaseBuild
	_, err := parseMasterKey()
	s.EqualError(err, ErrNoMasterKey.Error())
	Build = DebugBuild
}

func TestConfigSuite(t *testing.T) {
	RunTestSuite(t, new(ConfigSuite))
}
