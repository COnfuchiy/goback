package test_suites

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"goback/api/middleware"
	"goback/api/response"
	"goback/domain/entity"
	"goback/tests/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type JwtAuthTestSuite struct {
	suite.Suite
	mockAuthService          *mocks.MockAuthService
	mockUserService          *mocks.MockUserService
	authMiddleware           *middleware.JwtAuthMiddleware
	user                     *entity.User
	ginHandler               *gin.Engine
	validAccessToken         string
	invalidAccessToken       string
	invalidUserIDAccessToken string
}

func (suite *JwtAuthTestSuite) SetupSuite() {
	suite.user = &entity.User{
		ID: uuid.New(),
	}
	suite.validAccessToken = "bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk"
	suite.invalidAccessToken = "bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk1"
	suite.invalidUserIDAccessToken = "bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk2"

	suite.mockAuthService = new(mocks.MockAuthService)
	suite.mockAuthService.On("ExtractIDFromAccessToken", suite.validAccessToken).Return(suite.user.ID.String(), nil)
	suite.mockAuthService.On("ExtractIDFromAccessToken", suite.invalidAccessToken).Return("", errors.New("invalid access token"))
	suite.mockAuthService.On("ExtractIDFromAccessToken", suite.invalidUserIDAccessToken).Return("invalid_uuid", nil)

	suite.mockUserService = new(mocks.MockUserService)
	suite.mockUserService.On("GetUserFromContext", suite.user.ID.String()).Return(suite.user, nil)
	suite.mockUserService.On("GetUserFromContext", "invalid_uuid").Return(nil, errors.New("invalid user uuid"))

	suite.authMiddleware = middleware.NewJwtAuthMiddleware(suite.mockAuthService, suite.mockUserService)
	suite.setupGin()
}

func (suite *JwtAuthTestSuite) TestCorrectGetUserData() {
	responseData, code := suite.fetchTestAuthRequest("/", true, suite.validAccessToken)
	suite.Require().Equal("OK", string(responseData))
	suite.Require().Equal(http.StatusOK, code)
}

func (suite *JwtAuthTestSuite) TestGetUserDataWithoutAuthorizationHeader() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "Not authorized"})
	suite.Require().Nil(err)
	responseData, code := suite.fetchTestAuthRequest("/", false, suite.validAccessToken)
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusUnauthorized, code)
}
func (suite *JwtAuthTestSuite) TestGetUserDataWithInvalidAccessToken() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "invalid access token"})
	suite.Require().Nil(err)
	responseData, code := suite.fetchTestAuthRequest("/", true, suite.invalidAccessToken)
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusUnauthorized, code)
}

func (suite *JwtAuthTestSuite) TestGetUserDataWithInvalidUserID() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "invalid user uuid"})
	suite.Require().Nil(err)
	responseData, code := suite.fetchTestAuthRequest("/", true, suite.invalidUserIDAccessToken)
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusUnauthorized, code)
}

func (suite *JwtAuthTestSuite) handleUserRequest(context *gin.Context) {
	userObject, isUserExist := context.Get("user")
	suite.Require().True(isUserExist)
	user, _ := userObject.(*entity.User)
	suite.Require().False(reflect.ValueOf(user).IsNil())
	suite.Require().Equal(user, suite.user)
	context.Data(200, "text/plain", []byte("OK"))
}

func (suite *JwtAuthTestSuite) fetchTestAuthRequest(url string, headerSet bool, token string) ([]byte, int) {
	request, err := http.NewRequest("GET", url, nil)
	suite.Require().Nil(err)
	if headerSet != false {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	responseRecorder := httptest.NewRecorder()
	suite.ginHandler.ServeHTTP(responseRecorder, request)
	responseData, err := io.ReadAll(responseRecorder.Body)
	suite.Require().Nil(err)
	return responseData, responseRecorder.Code
}

func (suite *JwtAuthTestSuite) setupGin() {
	suite.ginHandler = gin.Default()
	route := suite.ginHandler.Group("/", suite.authMiddleware.Handle())
	route.GET("/", suite.handleUserRequest)
}
