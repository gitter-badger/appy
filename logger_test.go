package appy_test

import (
	"reflect"
	"testing"

	"github.com/appist/appy"
)

type LoggerSuite struct {
	appy.TestSuite
}

func (s *LoggerSuite) TestNewLogger() {
	logger := appy.NewLogger()
	_, ok := reflect.TypeOf(logger).MethodByName("Desugar")
	s.Equal(true, ok)
}

func (s *LoggerSuite) TestNewFakeLogger() {
	logger, buf, writer := appy.NewFakeLogger()
	logger.Info("test")
	writer.Flush()
	s.NotNil(logger)
	s.Contains(buf.String(), "\ttest")

	appy.Build = appy.ReleaseBuild
	defer func() {
		appy.Build = appy.DebugBuild
	}()

	logger, buf, writer = appy.NewFakeLogger()
	logger.Info("test")
	writer.Flush()
	s.NotNil(logger)
	s.Contains(buf.String(), "info\ttest")
}

func (s *LoggerSuite) TestSetDBLogging() {
	logger := appy.NewLogger()
	s.Equal(true, logger.DBLogging())
	logger.SetDBLogging(false)
	s.Equal(false, logger.DBLogging())
}

func TestLoggerSuite(t *testing.T) {
	appy.RunTestSuite(t, new(LoggerSuite))
}
