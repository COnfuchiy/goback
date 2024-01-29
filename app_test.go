package main

import (
	"github.com/stretchr/testify/suite"
	"goback/tests/test_suites"
	"testing"
)

func TestJwtAuthSuite(t *testing.T) {
	suite.Run(t, new(test_suites.JwtAuthTestSuite))
}
func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(test_suites.AuthTestSuite))
}
func TestAppSuite(t *testing.T) {
	suite.Run(t, new(test_suites.AppTestSuite))
}
