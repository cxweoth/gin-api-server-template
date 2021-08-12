package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cxweoth/gin-api-server-template/internal/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Auth failed resp, will be used when 401 in all api path
type AuthFailedResp struct {
	ErrorString string `json:"error" example:"error"`
}

// A function to check whther APIKey met
func ValidateAPIKey(c *gin.Context) {

	// Init auth failed struct
	var authFailedResp = AuthFailedResp{}

	// Fetch logger
	logger := c.MustGet("Logger").(*logrus.Entry)

	// Read apikeys
	apikeyFilePath := c.MustGet("APIkeyFilePath").(string)
	apikeyMap, err := utils.ReadUnstructuredJsonFile(apikeyFilePath)

	if err != nil {
		logger.Warn("read API-Key file failed: " + err.Error())
		authFailedResp.ErrorString = "API-Key process failed, please connect admin."
		c.JSON(http.StatusInternalServerError, authFailedResp)
		c.Abort()
		return
	}

	// Read apikey from request
	APIKey := c.Request.Header.Get("X-API-Key")

	// Check apikey is in memo apikeys
	for key, value := range apikeyMap {
		if APIKey == value {

			// Set to memo which client do the access
			c.Set("client", key)

			// If met, do nexy
			c.Next()

			logger.Info(key + " user API-Key authentication succeed")
			return
		}
	}

	// No apikeys met, return auth failed
	authFailedResp.ErrorString = "no such API-Key, authentication failed"
	logger.Warn("no such API-Key, authentication failed")
	c.JSON(http.StatusUnauthorized, authFailedResp)
	c.Abort()

	return
}

// Claim format in JWT
type Claims struct {
	Account string `json:"account"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

// A function to check whether JWT met
func AuthRequired(c *gin.Context) {

	//Init auth failed struct
	var authFailedResp = AuthFailedResp{}

	// Fetch logger
	logger := c.MustGet("Logger").(*logrus.Entry)

	// Read bearer token from request
	auth := c.GetHeader("Authorization")
	bearerSlice := strings.Split(auth, "Bearer ")
	if len(bearerSlice) != 2 {

		authFailedResp.ErrorString = "bearer format is not correct"
		c.JSON(http.StatusUnauthorized, authFailedResp)

		c.Abort()
		return
	}

	// Fetch token
	token := bearerSlice[1]

	// Fetch jwt secret
	jwtSecret := c.MustGet("JwtSecret").([]byte)

	// parse and validate token for six things:
	// validationErrorMalformed => token is malformed
	// validationErrorUnverifiable => token could not be verified because of signing problems
	// validationErrorSignatureInvalid => signature validation failed
	// validationErrorExpired => exp validation failed
	// validationErrorNotValidYet => nbf validation failed
	// validationErrorIssuedAt => iat validation failed
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})

	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}

		authFailedResp.ErrorString = message
		c.JSON(http.StatusUnauthorized, authFailedResp)

		logger.Warn("API querried auth failed - " + token + ": " + message)

		c.Abort()
		return
	}

	// Check whether token is valid
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {

		c.Set("account", claims.Account)
		c.Set("role", claims.Role)

		c.Next()

		logger.Info("Account " + claims.Account + " API querried succeed ")

		return

	} else {

		authFailedResp.ErrorString = "token is invalid"
		c.JSON(http.StatusUnauthorized, authFailedResp)

		c.Abort()
		return
	}
}

// A function to generate token
func GenerateToken(jwtSecret []byte, account, role string) (string, error) {

	// Set jwt id for token, and include time to id
	now := time.Now()
	jwtId := account + strconv.FormatInt(now.Unix(), 10)

	// Set role
	setRole := role

	// Set claims
	claims := Claims{
		Account: account,
		Role:    setRole,
		StandardClaims: jwt.StandardClaims{
			Audience:  account,
			ExpiresAt: now.Add(20 * time.Minute).Unix(), // expired time: 20 mins later
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "JWT",
			NotBefore: now.Unix(), // workable time
			Subject:   account,
		},
	}

	// sign the claims
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return token, nil
}

// Login receive and response struct

type LoginReceiveBody struct {
	Account  string `json:"Account" example:"account" format:"string"`
	Password string `json:"Password" example:"password" format:"string"`
}

type LoginSucceed struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Ijoib3JnYWRtaW4iLCJyb2xlIjoiTWVtYmVyIiwiYXVkIjoib3JnYWRtaW4iLCJleHAiOjE2MjcwMzAzMjcsImp0aSI6Im9yZ2FkbWluMTYyNzAyOTEyNyIsImlhdCI6MTYyNzAyOTEyNywiaXNzIjoiSldUIiwibmJmIjoxNjI3MDI5MTI3LCJzdWIiOiJvcmdhZG1pbiJ9._HlIMv_2_vok5cLjxPNTI2qLUsIZQTtKbpW8UN-C4Tc" format:"string"`
}

type LoginFailed struct {
	Message string `json:"message" example:"msg" format:"string"`
}

// @Summary Login and return token after authenticate.
// @Description login and get token
// @Param X-API-Key header string true "Insert your api key" default(<Add api key here>)
// @Accept  json
// @Produce  json
// @Param Body body LoginReceiveBody true "Account and Password"
// @Tags AAA
// @version 1.0
// @produce text/plain
// @Success 200 {object} LoginSucceed
// @Failure 400 {object} LoginFailed
// @Failure 401 {object} AuthFailedResp
// @Failure 500 {object} LoginFailed
// @Router /api/v1/login [post]
func Login(c *gin.Context) {

	// Fetch logger
	logger := c.MustGet("Logger").(*logrus.Entry)

	// Fetch client
	client := c.MustGet("client").(string)

	// Fetch body received
	var receiveBody = LoginReceiveBody{}

	err := c.ShouldBindJSON(&receiveBody)
	if err != nil {
		var loginFailed = LoginFailed{}
		loginFailed.Message = "Client " + client + " bad request: " + err.Error()
		c.JSON(http.StatusBadRequest, loginFailed)
		logger.Warn("Client " + client + " bad request: " + err.Error())
		return
	}

	logger.Info("Client " + client + " try to login account " + receiveBody.Account)

	// Auth to blockchain CA
	err = AuthFunction(receiveBody.Account, receiveBody.Password)

	if err != nil {
		var loginFailed = LoginFailed{}
		loginFailed.Message = "Account " + receiveBody.Account + " auth blockchain CA failed."
		c.JSON(http.StatusBadRequest, loginFailed)
		logger.Warn("Client " + client + " try to login account " + receiveBody.Account + ", but auth blockchain CA failed and error: " + err.Error())
		return
	}

	logger.Info("Client " + client + " try to login account " + receiveBody.Account + " auth blockchain CA succeed!")

	// Fetch jwt secret
	jwtSecret := c.MustGet("JwtSecret").([]byte)

	// Generate token
	token, err := GenerateToken(jwtSecret, receiveBody.Account, "Member")

	if err != nil {
		var loginFailed = LoginFailed{}
		loginFailed.Message = "Account " + receiveBody.Account + " generate token failed."
		c.JSON(http.StatusBadRequest, loginFailed)
		logger.Warn("Client " + client + " try to login account " + receiveBody.Account + " generate token failed: " + err.Error())
		return
	}

	logger.Info("Client " + client + " try to login account " + receiveBody.Account + " generate token succeed!")

	// Succeed and return token
	var loginSucceed = LoginSucceed{}
	loginSucceed.Token = token
	c.JSON(http.StatusOK, loginSucceed)

	logger.Info("Client " + client + " try to login account " + receiveBody.Account + " and already return token")
	return
}

func AuthFunction(account, pwd string) error {
	return nil
}
