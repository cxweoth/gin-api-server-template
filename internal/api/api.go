package api

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files

	apiDocs "github.com/cxweoth/gin-api-server-template/api/docs"
	"github.com/cxweoth/gin-api-server-template/internal/conf"
	"github.com/cxweoth/gin-api-server-template/internal/utils"
)

// This function is used to setup cors, and setup allowed methods and headers
func CorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowCredentials = true
	corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
	corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers", "X-API-Key", "Access-Control-Allow-Origin"}
	return corsConf
}

// APIMiddleware will add middleware to the context
func APIMiddleware(cfg conf.IConf, logger *logrus.Entry, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

		apiCfg := cfg.APICfg()
		apiServiceName := apiCfg.APIServiceName
		apiKeyFilePath := apiCfg.APIKeyFilePath

		c.Set("Logger", logger)
		c.Set("JwtSecret", jwtSecret)
		c.Set("APIServiceName", apiServiceName)
		c.Set("APIkeyFilePath", apiKeyFilePath)

		c.Next()
	}
}

func SetupServer(cfg conf.IConf, logger *logrus.Entry) (*gin.Engine, error) {

	// Fetch cfg params
	apiCfg := cfg.APICfg()

	apiMode := apiCfg.APIMode
	apiProtocol := apiCfg.APIProtocol
	apiHost := apiCfg.APIHost
	apiPort := apiCfg.APIPort

	// Set mode of api server
	if apiMode == "Debug" {
		gin.SetMode(gin.DebugMode)
	} else if apiMode == "Release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		return nil, errors.New("no such mode")
	}

	// Init server
	server := gin.Default()
	server.Use(cors.New(CorsConfig()))

	// Generate JWT secret which is used to generate token
	jwtSecret := utils.GenerateRandomBytes(32)

	if jwtSecret == nil {
		logger.Fatal("Generate jwt secret failed")
		return nil, errors.New("Generate jwt secret failed")
	}

	// Setup middleware
	server.Use(APIMiddleware(cfg, logger, jwtSecret))

	// Setup max memory can be used in each request
	server.MaxMultipartMemory = 32 << 20 // 32MiB

	// Set swagger document and swagger GET
	if mode := gin.Mode(); mode == gin.DebugMode {
		apiDocs.SwaggerInfo.Title = "API Service"
		apiDocs.SwaggerInfo.Description = "Supports API access to login and fetch service info."
		apiDocs.SwaggerInfo.Version = "1.0"
		apiDocs.SwaggerInfo.Host = apiHost + ":" + apiPort
		apiDocs.SwaggerInfo.Schemes = []string{apiProtocol}
		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API Key authorized group
	apiKeyAuthorized := server.Group("/")

	// ValludateAPIKey is in aaa.go file
	apiKeyAuthorized.Use(ValidateAPIKey)
	{
		// Login, PulsarLogin are in aaa.go
		// Use them to generate JWT
		apiKeyAuthorized.POST("/api/v1/login", Login)
	}

	// JWT authorized group
	tokenAuthorized := server.Group("/")
	// AuthRequired is in aaa.go file
	tokenAuthorized.Use(AuthRequired)
	{
		// UploadComponent, DownloadComponent are in fileTransmission.go
		tokenAuthorized.GET("/api/v1/getServiceInfo", GetServiceInfo)

	}

	return server, nil
}

// RunServer is used to run API server
func RunServer(cfg conf.IConf, logger *logrus.Entry) error {

	// Fetch cfg params
	apiCfg := cfg.APICfg()

	// Setup server
	server, err := SetupServer(cfg, logger)
	if err != nil {
		return errors.New("API server setup failed")
	}

	// Run server
	server.Run(":" + apiCfg.APIPort)

	return errors.New("API server shutdown")
}
