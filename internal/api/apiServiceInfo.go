package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Service list response struct
type ServiceInfoSuccessResp struct {
	ServiceName string `json:"serviceName" example:"service" format:"string"`
}

type ServiceInfoFailedResp struct {
	Message string `json:"message" example:"msg"`
}

// API to get service list
// Swagger comments:
// @Summary get service info
// @Description get service info
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Accept   json
// @Produce  json
// @Tags Service Information
// @version 1.0
// @Success 200 {object} ServiceInfoSuccessResp "Get service info by GET method with token"
// @Failure 400 {object} ServiceInfoFailedResp
// @Failure 401 {object} AuthFailedResp
// @Failure 500 {object} ServiceInfoFailedResp
// @Router /api/v1/getServiceInfo [get]
func GetServiceInfo(c *gin.Context) {

	// Fetch logger
	logger := c.MustGet("Logger").(*logrus.Entry)

	// Fetch account
	account := c.MustGet("account").(string)

	// Fetch APIServiceName
	apiServiceName := c.MustGet("APIServiceName").(string)

	// Init struct to fetch service list
	var serviceInfoSucceed ServiceInfoSuccessResp

	// Set service name
	serviceInfoSucceed.ServiceName = apiServiceName

	// Succeed and return service list
	c.JSON(http.StatusOK, serviceInfoSucceed)

	logger.Info(account + " fetch service info")
	return
}
