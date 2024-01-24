package test_suites

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"goback/api/controller"
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
	"goback/mapper"
	"goback/tests/mocks"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type AppTestSuite struct {
	suite.Suite
	ginHandler               *gin.Engine
	mockWorkspaceService     *mocks.MockWorkspaceService
	mockUserService          *mocks.MockUserService
	mockFileService          *mocks.MockFileService
	mockFileStorageService   *mocks.MockFileStorageService
	mockFileHistoriesService *mocks.MockFileHistoriesService
	workspaceController      *controller.WorkspaceController
	fileHistoryController    *controller.FileHistoryController
	userController           *controller.UserController
	fileController           *controller.FileController
	userMapper               mapper.IUserMapper
	workspaceMapper          mapper.IWorkspaceMapper
	fileHistoriesMapper      mapper.IFileHistoryMapper
	fileMapper               mapper.IFileMapper
	paginateMapper           mapper.IPaginationMapper

	user                *entity.User
	workspace           *entity.Workspace
	totalWorkspaceCount int
	fileHistory         *entity.FileHistory
	file                *entity.File
}

func (suite *AppTestSuite) SetupSuite() {
	suite.userMapper = mapper.UserMapper{}
	suite.workspaceMapper = mapper.NewWorkspaceMapper(suite.userMapper)
	suite.fileMapper = mapper.NewFileMapper(suite.userMapper)
	suite.fileHistoriesMapper = mapper.NewFileHistoryMapper(suite.fileMapper)
	suite.paginateMapper = mapper.NewPaginationMapper(10)

	username := faker.Name()
	email := faker.Email()
	password := faker.Word()
	suite.user = &entity.User{
		ID:         uuid.New(),
		Username:   &username,
		Email:      &email,
		Password:   &password,
		Workspaces: nil,
	}
	suite.totalWorkspaceCount = 20

	for i := 0; i < suite.totalWorkspaceCount; i++ {
		workspace := &entity.Workspace{
			ID:             uuid.New(),
			Name:           faker.Word(),
			CreatorID:      suite.user.ID,
			Creator:        *suite.user,
			FilesHistories: []entity.FileHistory{},
		}

		for j := 0; j < 5; j++ {
			fileHistory := &entity.FileHistory{
				ID:          uuid.New(),
				WorkspaceID: workspace.ID,
				Workspace:   *workspace,
				Files:       []entity.File{},
			}

			for k := 0; k < 5; k++ {
				file := &entity.File{
					ID:            uuid.New(),
					FileHistoryID: fileHistory.ID,
					FileHistory:   *fileHistory,
					Filename:      faker.Word(),
					Tag:           "",
					Size:          int64(rand.Intn(10000000)),
					DownloadURL:   "",
					UserId:        suite.user.ID,
					User:          *suite.user,
				}

				fileHistory.Files = append(fileHistory.Files, *file)
			}

			workspace.FilesHistories = append(workspace.FilesHistories, *fileHistory)
		}

		suite.user.Workspaces = append(suite.user.Workspaces, *workspace)
	}

	suite.workspace = &suite.user.Workspaces[0]
	suite.fileHistory = &suite.workspace.FilesHistories[0]
	suite.file = &suite.fileHistory.Files[0]

	suite.mockWorkspaceService = new(mocks.MockWorkspaceService)
	suite.mockWorkspaceService.On("GetAllByUserID", suite.user.ID, 1).Return(suite.user.Workspaces[0:10], suite.totalWorkspaceCount, nil)
	suite.mockWorkspaceService.On("GetAllByUserID", suite.user.ID, 2).Return(suite.user.Workspaces[10:20], suite.totalWorkspaceCount, nil)

	suite.mockFileHistoriesService = new(mocks.MockFileHistoriesService)
	suite.mockFileHistoriesService.On("GetAllByWorkspaceID", suite.workspace.ID, 1).Return(suite.workspace.FilesHistories, len(suite.workspace.FilesHistories), nil)
	suite.mockFileHistoriesService.On("GetFromContext", suite.workspace.FilesHistories[0].ID.String()).Return(&suite.workspace.FilesHistories[0], nil)
	suite.mockFileHistoriesService.On("Create", suite.workspace.ID.String()).Return(nil)

	suite.mockFileService = new(mocks.MockFileService)
	suite.mockFileService.On("Create", suite.file.Filename).Return(nil)
	suite.mockFileService.On("GetFromContext", suite.file.ID.String()).Return(suite.file, nil)
	suite.mockFileService.On("CheckExisting", suite.file.Filename).Return(true, nil)

	suite.mockFileStorageService = new(mocks.MockFileStorageService)

	suite.userController = controller.NewUserController(suite.mockUserService, suite.mockWorkspaceService, suite.userMapper, suite.workspaceMapper, suite.paginateMapper)
	suite.workspaceController = controller.NewWorkspaceController(suite.mockWorkspaceService, suite.mockFileHistoriesService, suite.workspaceMapper, suite.fileHistoriesMapper, suite.paginateMapper)
	suite.fileHistoryController = controller.NewFileHistoryController(suite.mockFileHistoriesService, suite.fileHistoriesMapper)
	suite.fileController = controller.NewFileController(suite.mockUserService, suite.mockFileHistoriesService, suite.mockFileService, suite.mockFileStorageService, suite.mockWorkspaceService, suite.fileMapper)

	suite.setupGin()

}

func (suite *AppTestSuite) setupGin() {
	suite.ginHandler = gin.Default()
	gin.SetMode(gin.TestMode)
	suite.ginHandler.Use(suite.middlewareHandle())
	suite.ginHandler.GET("/get-all-user-workspaces", suite.userController.GetAllWorkspaces)
	suite.ginHandler.POST("/create-workspace", suite.workspaceController.CreateWorkspace)
	suite.ginHandler.GET("/get-workspace", suite.workspaceController.GetWorkspace)
	suite.ginHandler.GET("/get-all-workspace-file-histories", suite.workspaceController.GetAllFilesHistories)
	suite.ginHandler.GET("/get-file-history/:file_history_id", suite.fileHistoryController.GetFileHistory)
	suite.ginHandler.POST("/create-file", suite.fileController.Create)
	suite.ginHandler.GET("/get-file-download-link/:file_id", suite.fileController.GetFileDownloadLink)
	suite.ginHandler.GET("/check-filename-existing", suite.fileController.CheckFilenameExisting)
}

func (suite *AppTestSuite) TestGetAllUserWorkspaces() {
	workspacesFirstPage := suite.user.Workspaces[0:10]
	workspacesSecondPage := suite.user.Workspaces[10:20]
	var workspacesFirstPageResponse response.WorkspacesResponse
	var workspacesSecondPageResponse response.WorkspacesResponse
	workspacesFirstPageResponse.Pagination = *suite.paginateMapper.ToPaginationResponse(int64(suite.totalWorkspaceCount), 1)
	workspacesSecondPageResponse.Pagination = *suite.paginateMapper.ToPaginationResponse(int64(suite.totalWorkspaceCount), 2)

	for _, workspace := range workspacesFirstPage {
		workspacesFirstPageResponse.Workspaces = append(workspacesFirstPageResponse.Workspaces, *suite.workspaceMapper.ToWorkspaceResponse(&workspace))
	}

	for _, workspace := range workspacesSecondPage {
		workspacesSecondPageResponse.Workspaces = append(workspacesSecondPageResponse.Workspaces, *suite.workspaceMapper.ToWorkspaceResponse(&workspace))
	}

	jsonFirstPageResponse, err := json.Marshal(workspacesFirstPageResponse)
	suite.Require().Nil(err)
	jsonSecondPageResponse, err := json.Marshal(workspacesSecondPageResponse)
	suite.Require().Nil(err)

	responseData, code := suite.fetchTestRequest("GET", "/get-all-user-workspaces?page=2", nil)
	suite.Require().Equal(jsonSecondPageResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

	responseData, code = suite.fetchTestRequest("GET", "/get-all-user-workspaces", nil)
	suite.Require().Equal(jsonFirstPageResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

}

func (suite *AppTestSuite) TestInvalidUserObjectInContext() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "User is not type of entity.User"})
	suite.Require().Nil(err)

	tmpUser := suite.user
	suite.user = nil

	responseData, code := suite.fetchTestRequest("GET", "/get-all-user-workspaces", nil)
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)

	suite.user = tmpUser
}

func (suite *AppTestSuite) TestInvalidWorkspaceObjectInContext() {
	jsonError, err := json.Marshal(response.ErrorResponse{Message: "Workspace is not type of entity.Workspace"})
	suite.Require().Nil(err)

	tmpWorkspace := suite.workspace
	suite.workspace = nil

	responseData, code := suite.fetchTestRequest("GET", "/get-workspace", nil)
	suite.Require().Equal(jsonError, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)

	suite.workspace = tmpWorkspace
}

func (suite *AppTestSuite) TestCreateWorkspace() {
	correctRequest := &request.CreateWorkspaceRequest{Name: faker.Word()}
	correctWorkspace := suite.workspaceMapper.FromCreateRequest(correctRequest)
	correctWorkspace.Creator = *suite.user
	correctResponse := suite.workspaceMapper.ToWorkspaceResponse(correctWorkspace)
	jsonCorrectResponse, err := json.Marshal(correctResponse)
	suite.Require().Nil(err)

	incorrectRequest := &request.CreateWorkspaceRequest{Name: faker.Word()}
	incorrectWorkspace := suite.workspaceMapper.FromCreateRequest(incorrectRequest)
	incorrectWorkspace.Creator = *suite.user
	jsonIncorrectResponse, err := json.Marshal(response.ErrorResponse{Message: "error with creating workspace"})
	suite.Require().Nil(err)

	suite.mockWorkspaceService.On("Create", correctWorkspace.Name).Return(nil)
	suite.mockWorkspaceService.On("Create", incorrectWorkspace.Name).Return(errors.New("error with creating workspace"))

	correctWorkspaceForm := suite.createWorkspaceForm(correctWorkspace.Name)
	responseData, code := suite.fetchTestRequest("POST", "/create-workspace", strings.NewReader(correctWorkspaceForm.Encode()))
	suite.Require().Equal(jsonCorrectResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

	incorrectWorkspaceForm := suite.createWorkspaceForm(incorrectWorkspace.Name)
	responseData, code = suite.fetchTestRequest("POST", "/create-workspace", strings.NewReader(incorrectWorkspaceForm.Encode()))
	suite.Require().Equal(jsonIncorrectResponse, responseData)
	suite.Require().Equal(http.StatusInternalServerError, code)

}

func (suite *AppTestSuite) TestGetWorkspace() {
	workspaceResponse := suite.workspaceMapper.ToWorkspaceResponse(suite.workspace)
	jsonResponse, err := json.Marshal(workspaceResponse)
	suite.Require().Nil(err)

	responseData, code := suite.fetchTestRequest("GET", "/get-workspace", nil)
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)
}

func (suite *AppTestSuite) TestGetAllWorkspaceFileHistories() {
	var fileHistoriesResponse response.FileHistoriesResponse

	fileHistoriesResponse.Pagination = *suite.paginateMapper.ToPaginationResponse(int64(len(suite.workspace.FilesHistories)), 1)

	for _, history := range suite.workspace.FilesHistories {
		fileHistoriesResponse.FileHistories = append(fileHistoriesResponse.FileHistories, *suite.fileHistoriesMapper.ToFileHistoryResponse(&history))
	}
	jsonfileHistoriesResponse, err := json.Marshal(fileHistoriesResponse)
	suite.Require().Nil(err)

	responseData, code := suite.fetchTestRequest("GET", "/get-all-workspace-file-histories", nil)
	suite.Require().Equal(jsonfileHistoriesResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

}

func (suite *AppTestSuite) TestGetFileHistory() {
	fileHistoryResponse := suite.fileHistoriesMapper.ToFileHistoryResponse(suite.fileHistory)
	jsonResponse, err := json.Marshal(fileHistoryResponse)
	suite.Require().Nil(err)

	uuidErrorJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "File history ID is not uuid"})
	suite.Require().Nil(err)

	emptyUuidJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "File history ID is not specified"})
	suite.Require().Nil(err)

	notFoundJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "file history not found"})
	suite.Require().Nil(err)

	errorFileHistoryID := uuid.New()
	suite.mockFileHistoriesService.On("GetFromContext", errorFileHistoryID.String()).Return(nil, errors.New("file history not found"))

	responseData, code := suite.fetchTestRequest("GET", "/get-file-history/"+suite.fileHistory.ID.String(), nil)
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-history/"+errorFileHistoryID.String(), nil)
	suite.Require().Equal(notFoundJsonResponse, responseData)
	suite.Require().Equal(http.StatusInternalServerError, code)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-history/f", nil)
	suite.Require().Equal(uuidErrorJsonResponse, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)

	suite.ginHandler.GET("/get-file-history", suite.fileHistoryController.GetFileHistory)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-history", nil)
	suite.Require().Equal(emptyUuidJsonResponse, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AppTestSuite) TestCreateFile() {
	fileResponse := suite.fileMapper.ToFileResponse(suite.file)
	jsonResponse, err := json.Marshal(fileResponse)
	suite.Require().Nil(err)

	tempFile, err := os.CreateTemp("", "test.txt")
	suite.Require().Nil(err)

	defer func(name string) {
		err = os.Remove(name)
		suite.Require().Nil(err)
	}(tempFile.Name())
	defer func(tempFile *os.File) {
		err = tempFile.Close()
		suite.Require().Nil(err)
	}(tempFile)

	_, err = tempFile.WriteString("Test file content")

	suite.Require().Nil(err)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("filename", suite.file.Filename)
	_ = writer.WriteField("size", strconv.Itoa(int(suite.file.Size)))
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	suite.Require().Nil(err)

	file, err := os.Open(tempFile.Name())
	suite.Require().Nil(err)

	defer func(file *os.File) {
		err = file.Close()
		suite.Require().Nil(err)
	}(file)

	_, err = io.Copy(fileWriter, file)
	suite.Require().Nil(err)

	err = writer.Close()
	suite.Require().Nil(err)

	suite.mockFileStorageService.On("SaveFileToStorage", uuid.UUID{}, "test.txt").Return(nil)

	req, err := http.NewRequest("POST", "/create-file", &requestBody)
	suite.Require().Nil(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	responseRecorder := httptest.NewRecorder()
	suite.ginHandler.ServeHTTP(responseRecorder, req)
	responseData, err := io.ReadAll(responseRecorder.Body)
	suite.Require().Nil(err)

	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, responseRecorder.Code)

}

func (suite *AppTestSuite) TestGetFileDownloadLink() {
	fileResponse := suite.fileMapper.ToDownloadFileLinkResponse(suite.file)
	jsonResponse, err := json.Marshal(fileResponse)
	suite.Require().Nil(err)

	uuidErrorJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "File ID is not uuid"})
	suite.Require().Nil(err)

	emptyUuidJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "File ID is not specified"})
	suite.Require().Nil(err)

	notFoundJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "file not found"})
	suite.Require().Nil(err)

	errorFileID := uuid.New()
	suite.mockFileService.On("GetFromContext", errorFileID.String()).Return(nil, errors.New("file not found"))

	responseData, code := suite.fetchTestRequest("GET", "/get-file-download-link/"+suite.file.ID.String(), nil)
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-download-link/"+errorFileID.String(), nil)
	suite.Require().Equal(notFoundJsonResponse, responseData)
	suite.Require().Equal(http.StatusInternalServerError, code)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-download-link/f", nil)
	suite.Require().Equal(uuidErrorJsonResponse, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)

	suite.ginHandler.GET("/get-file-download-link", suite.fileController.GetFileDownloadLink)

	responseData, code = suite.fetchTestRequest("GET", "/get-file-download-link", nil)
	suite.Require().Equal(emptyUuidJsonResponse, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)
}

func (suite *AppTestSuite) TestCheckFilenameExisting() {
	jsonResponse, err := json.Marshal(response.CheckFileResponse{IsFileExist: true})
	suite.Require().Nil(err)

	noFilenameProvidesJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "no filename provides"})
	suite.Require().Nil(err)

	errorFilenameJsonResponse, err := json.Marshal(response.ErrorResponse{Message: "server error"})
	suite.Require().Nil(err)

	errorFilename := faker.Word()
	suite.mockFileService.On("CheckExisting", errorFilename).Return(nil, errors.New("server error"))

	responseData, code := suite.fetchTestRequest("GET", "/check-filename-existing?filename="+suite.file.Filename, nil)
	suite.Require().Equal(jsonResponse, responseData)
	suite.Require().Equal(http.StatusOK, code)

	responseData, code = suite.fetchTestRequest("GET", "/check-filename-existing", nil)
	suite.Require().Equal(noFilenameProvidesJsonResponse, responseData)
	suite.Require().Equal(http.StatusBadRequest, code)

	responseData, code = suite.fetchTestRequest("GET", "/check-filename-existing?filename="+errorFilename, nil)
	suite.Require().Equal(errorFilenameJsonResponse, responseData)
	suite.Require().Equal(http.StatusInternalServerError, code)
}

func (suite *AppTestSuite) createWorkspaceForm(name string) url.Values {
	workspaceForm := url.Values{}
	workspaceForm.Add("name", name)
	return workspaceForm
}

func (suite *AppTestSuite) fetchTestRequest(method, url string, body io.Reader) ([]byte, int) {
	newRequest, err := http.NewRequest(method, url, body)

	suite.Require().Nil(err)
	if method == "POST" {
		newRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	responseRecorder := httptest.NewRecorder()
	suite.ginHandler.ServeHTTP(responseRecorder, newRequest)
	responseData, err := io.ReadAll(responseRecorder.Body)
	suite.Require().Nil(err)
	return responseData, responseRecorder.Code
}

func (suite *AppTestSuite) middlewareHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("user", suite.user)
		context.Set("workspace", suite.workspace)
		context.Next()
	}
}
