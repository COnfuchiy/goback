package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goback/api/route"
	"goback/bootstrap"
	docs "goback/docs"
	"goback/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type App struct {
	db  *gorm.DB
	gin *gin.Engine
	env *bootstrap.Env
}

func NewApp() *App {
	app := App{}
	app.env = bootstrap.NewEnv()

	app.setupDatabase()
	app.setupServer()
	app.setupRoutes()

	docs.SwaggerInfo.BasePath = "/api/v1"
	app.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &app
}

func (app *App) StartApp() error {
	return app.gin.Run()
}

func (app *App) setupDatabase() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", app.env.DBUser, app.env.DBPass, app.env.DBHost, app.env.DBPort, app.env.DBName)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&entity.User{}, &entity.Workspace{}, &entity.FileHistory{}, &entity.File{})
	if err != nil {
		log.Fatalln(err)
	}
	app.db = db
	return
}
func (app *App) setupServer() {
	app.gin = gin.Default()
}
func (app *App) setupRoutes() {
	route.Init(app.db, app.gin, app.env)
}
