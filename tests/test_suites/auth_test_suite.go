package test_suites

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"goback/api/controller"
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
	"goback/tests/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	// "strings"
)

type AuthTestSuite struct {
	suite.Suite
	mockAuthService     *mocks.MockAuthService
	mockUserService     *mocks.MockUserService
	mockUserMapper      *mocks.MockUserMapper
	authController      *controller.AuthController
	ginHandler          *gin.Engine
	user                *entity.User
	accessToken         string
	refreshToken        string
	incorrectEmail      string
	incorrectPassword   string
	invalidRefreshToken string
	newUser             *entity.User
	errorUser           *entity.User
}

func (suite *AuthTestSuite) SetupSuite() {
	userName := "Any"
	email := "Any@test.ru"
	newEmail := "An@test.ru"
	rawPassword := "pass"
	createTokensErrorEmail := "payl@test.ru"
	createTokensErrorPassword := "p"
	suite.user = &entity.User{
		ID:       uuid.New(),
		Username: &userName,
		Email:    &email,
		Password: &rawPassword,
	}
	suite.newUser = &entity.User{
		ID:       uuid.New(),
		Username: &userName,
		Email:    &newEmail,
		Password: &rawPassword,
	}

	suite.errorUser = &entity.User{
		Username: &userName,
		Email:    &createTokensErrorEmail,
		Password: &createTokensErrorPassword,
	}

	suite.incorrectEmail = "brom@test.ru"
	suite.incorrectPassword = "pas"
	suite.refreshToken = "bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk"
	suite.accessToken = "bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk1"
	suite.invalidRefreshToken = "2bfgibhuivdfjklgnjeryuihfovneur8fhoidn2389efuionsdjk2"

	suite.mockUserService = new(mocks.MockUserService)
	suite.mockUserService.On("GetByEmail", email).Return(suite.user, nil)
	suite.mockUserService.On("GetByEmail", newEmail).Return(nil, errors.New(""))
	suite.mockUserService.On("GetByEmail", "empty@empty.com").Return(nil, errors.New(""))
	suite.mockUserService.On("GetByEmail", suite.incorrectEmail).Return(nil, errors.New(""))
	suite.mockUserService.On("GetByEmail", *suite.errorUser.Email).Return(suite.errorUser, nil)
	suite.mockUserService.On("Create", suite.newUser).Return(nil)
	suite.mockUserService.On("Create", suite.errorUser).Return(errors.New("error while creating user"))
	suite.mockUserService.On("GetUserFromContext", suite.user.ID.String()).Return(suite.user, nil)

	suite.mockAuthService = new(mocks.MockAuthService)
	suite.mockAuthService.On("ComparePassword", rawPassword, rawPassword).Return(nil)
	suite.mockAuthService.On("ComparePassword", rawPassword, suite.incorrectPassword).Return(errors.New(""))
	suite.mockAuthService.On("ComparePassword", *suite.errorUser.Password, *suite.errorUser.Password).Return(nil)
	suite.mockAuthService.On("CreateAccessAndRefreshTokens", suite.user).Return(suite.accessToken, suite.refreshToken, nil)
	suite.mockAuthService.On("CreateAccessAndRefreshTokens", suite.errorUser).Return(nil, nil, errors.New("error while creating token"))
	suite.mockAuthService.On("CreateAccessAndRefreshTokens", suite.newUser).Return(suite.accessToken, suite.refreshToken, nil)
	suite.mockAuthService.On("HashPassword", rawPassword).Return(rawPassword, nil)
	suite.mockAuthService.On("HashPassword", "totally_empty").Return("", errors.New("password are empty"))
	suite.mockAuthService.On("HashPassword", "empty").Return("empty", nil)
	suite.mockAuthService.On("ExtractIDFromRefreshToken", suite.refreshToken).Return(suite.user.ID.String(), nil)
	suite.mockAuthService.On("ExtractIDFromRefreshToken", suite.invalidRefreshToken).Return("", errors.New("invalid refresh token"))

	suite.mockUserMapper = new(mocks.MockUserMapper)
	suite.mockUserMapper.On("FromRegisterRequest", &request.RegisterRequest{
		Email:    *suite.newUser.Email,
		Username: *suite.newUser.Username,
		Password: *suite.newUser.Password,
	}).Return(suite.newUser)
	suite.mockUserMapper.On("FromRegisterRequest", &request.RegisterRequest{
		Email:    "empty@empty.com",
		Username: "empty",
		Password: "empty",
	}).Return(suite.errorUser)

	suite.authController = controller.NewAuthController(suite.mockUserService, suite.mockAuthService, suite.mockUserMapper)

	suite.setupGin()
}

func (suite *AuthTestSuite) fetchTestRequest(method, url string, body io.Reader) ([]byte, int) {
	request, err := http.NewRequest(method, url, body)
	suite.Require().Nil(err)
	if method == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	responseRecorder := httptest.NewRecorder()
	suite.ginHandler.ServeHTTP(responseRecorder, request)
	responseData, err := io.ReadAll(responseRecorder.Body)
	suite.Require().Nil(err)
	return responseData, responseRecorder.Code
}

func (suite *AuthTestSuite) TestCorrectLogin() {
	jsonResponse, err := json.Marshal(response.RegisterResponse{
		Success:      true,
		AccessToken:  suite.accessToken,
		RefreshToken: suite.refreshToken,
	})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("", *suite.user.Email, *suite.user.Password)
	responseData, code := suite.fetchTestRequest("POST", "/login", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)
}

func (suite *AuthTestSuite) TestLoginWithIncorrectEmail() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "Invalid credentials"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("", suite.incorrectEmail, *suite.user.Password)
	responseData, code := suite.fetchTestRequest("POST", "/login", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestLoginWithIncorrectPassword() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "Invalid credentials"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("", *suite.user.Email, suite.incorrectPassword)
	responseData, code := suite.fetchTestRequest("POST", "/login", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestLoginWithCreateTokensError() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "error while creating token"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("", *suite.errorUser.Email, *suite.errorUser.Password)
	responseData, code := suite.fetchTestRequest("POST", "/login", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestLoginWithoutEmail() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("", "", *suite.user.Password)
	userForm.Del("email")
	responseData, code := suite.fetchTestRequest("POST", "/login", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestCorrectRegister() {
	jsonResponse, err := json.Marshal(response.RegisterResponse{
		Success:      true,
		AccessToken:  suite.accessToken,
		RefreshToken: suite.refreshToken,
	})
	suite.Require().Nil(err)

	userForm := suite.createUserForm(*suite.newUser.Username, *suite.newUser.Email, *suite.newUser.Password)
	responseData, code := suite.fetchTestRequest("POST", "/register", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)
}

func (suite *AuthTestSuite) TestRegisterWithEmptyPassword() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "password are empty"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm(*suite.newUser.Username, *suite.newUser.Email, "totally_empty")
	responseData, code := suite.fetchTestRequest("POST", "/register", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestRegisterCreateUserError() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "error while creating user"})
	suite.Require().Nil(err)

	userForm := suite.createUserForm("empty", "empty@empty.com", "empty")
	responseData, code := suite.fetchTestRequest("POST", "/register", strings.NewReader(userForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) TestCorrectRefreshToken() {
	jsonResponse, err := json.Marshal(response.RegisterResponse{
		Success:      true,
		AccessToken:  suite.accessToken,
		RefreshToken: suite.refreshToken,
	})
	suite.Require().Nil(err)

	refreshTokenForm := suite.createRefreshTokenForm(suite.refreshToken)
	responseData, code := suite.fetchTestRequest("POST", "/refresh", strings.NewReader(refreshTokenForm.Encode()))
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)
}

func (suite *AuthTestSuite) TestRefreshTokenWithInvalidAccessToken() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "invalid refresh token"})
	suite.Require().Nil(err)

	refreshTokenForm := suite.createRefreshTokenForm(suite.invalidRefreshToken)
	responseData, code := suite.fetchTestRequest("POST", "/refresh", strings.NewReader(refreshTokenForm.Encode()))
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AuthTestSuite) createUserForm(userName, email, password string) url.Values {
	userForm := url.Values{}
	userForm.Add("username", userName)
	userForm.Add("email", email)
	userForm.Add("password", password)
	return userForm
}

func (suite *AuthTestSuite) createRefreshTokenForm(refreshToken string) url.Values {
	refreshTokenForm := url.Values{}
	refreshTokenForm.Add("refresh_token", refreshToken)
	return refreshTokenForm
}

func (suite *AuthTestSuite) setupGin() {
	suite.ginHandler = gin.Default()
	suite.ginHandler.POST("/login", suite.authController.Login)
	suite.ginHandler.POST("/register", suite.authController.Register)
	suite.ginHandler.POST("/refresh", suite.authController.RefreshToken)
}
