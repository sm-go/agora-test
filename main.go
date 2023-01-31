package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

var (
	appId, appCertificate string
)

func main() {
	os.Setenv("APP_ID", "2327452401534f96b695e0d05a9c924b")
	os.Setenv("APP_CERTIFICATE", "06a51d8a48b947058628d1c33b5b4e5d")
	appIdEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not property configured, check APP_ID and APP_CERTIFICATE")
	} else {
		appId = appIdEnv
		appCertificate = appCertEnv
	}
	api := gin.Default()
	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ping",
		})
	})

	api.GET("rtc/:channelName/:role/:tokenType/:uid/", getRtcToken)

	api.GET("rtm/:uid/", getRtmToken)

	api.GET("rte/:channelName/:role/:tokenType/:uid", getBothTokens)

	api.Run(":8080")
}

func getRtcToken(ctx *gin.Context) {
	// get param values
	channelName, tokeType, uidStr, role, expireTimestamp, err := parseRtcParams(ctx)
	if err != nil {
		ctx.Error(err)
		errMsg := "error generating RTC token" + err.Error()
		ctx.AbortWithStatusJSON(400, gin.H{
			"status":  400,
			"message": errMsg,
		})
		return
	}
	// generate the token
	rtcToken, tokenErr := generateRtcToken(channelName, uidStr, tokeType, role, expireTimestamp)
	// return the token in JSON response
	if tokenErr != nil {
		log.Println(tokenErr)
		ctx.Error(err)
		errMsg := "error generating RTC token : " + tokenErr.Error()
		ctx.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		ctx.JSON(200, gin.H{
			"rtcToken": rtcToken,
		})
	}
}

func getRtmToken(ctx *gin.Context) {
	// get param values
	uidStr, expireTimestamp, err := parseRtmParams(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(400, gin.H{
			"status":  "400",
			"message": "error generating TRM token",
		})
		return
	}
	//build rtm token
	rtmToken, tokenErr := rtmtokenbuilder.BuildToken(appId, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	//return rtm token
	if tokenErr != nil {
		log.Println(err)
		ctx.Error(err)
		errMsg := "error generating RTM token" + tokenErr.Error()
		ctx.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		ctx.JSON(200, gin.H{
			"rtmToken": rtmToken,
		})
	}
}

func getBothTokens(ctx *gin.Context) {
	// get the params
	channelName, tokenType, uidStr, role, expireTimestamp, rtcParamErr := parseRtcParams(ctx)
	if rtcParamErr != nil {
		ctx.Error(rtcParamErr)
		ctx.AbortWithStatusJSON(400, gin.H{
			"status":  400,
			"message": "error generating tokens" + rtcParamErr.Error(),
		})
	}
	// generate rtc token
	rtcToken, rtcTokenErr := generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)
	// generate rtm token
	rtmToken, rtmTokenErr := rtmtokenbuilder.BuildToken(appId, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	// return both tokens
	if rtcTokenErr != nil {
		ctx.Error(rtcTokenErr)
		errMsg := "error generating RTC token" + rtcTokenErr.Error()
		ctx.AbortWithStatusJSON(400, gin.H{
			"status":  400,
			"message": errMsg,
		})
	} else if rtmTokenErr != nil {
		ctx.Error(rtmTokenErr)
		errMsg := "error generating RTM token : " + rtmTokenErr.Error()
		ctx.AbortWithStatusJSON(400, gin.H{
			"status":  400,
			"message": errMsg,
		})
	} else {
		ctx.JSON(200, gin.H{
			"rtcToken": rtcToken,
			"rtmToken": rtmToken,
		})
	}
}

func parseRtcParams(ctx *gin.Context) (channelName, tokenType, uidStr string, role rtctokenbuilder.Role, expireTimestamp uint32, err error) {
	channelName = ctx.Param("channelName")
	roleStr := ctx.Param("role")
	tokenType = ctx.Param("tokenType")
	uidStr = ctx.Param("uid")
	expireTime := ctx.DefaultQuery("expiry", "3600")
	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
	}
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		err = fmt.Errorf("failed to parse expireTime : %s causing error : %s", expireTime, parseErr)
	}
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds
	return channelName, tokenType, uidStr, role, expireTimestamp, err
}

func parseRtmParams(ctx *gin.Context) (uidStr string, expireTimestamp uint32, err error) {
	// get params
	uidStr = ctx.Param("uid")
	expireTime := ctx.DefaultQuery("expiry", "3600")
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		err = fmt.Errorf("failed to parse expireTime : %s, causing error : %s", expireTime, parseErr)
	}
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds
	return uidStr, expireTimestamp, err
}

func generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder.Role, expireTimestamp uint32) (rtcToken string, err error) {
	// check token type
	if tokenType == "userAccount" {
		rtcToken, err = rtctokenbuilder.BuildTokenWithUserAccount(appId, appCertificate, channelName, uidStr, role, expireTimestamp)
		return rtcToken, err
	} else if tokenType == "uid" {
		uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse uidStr :%s to uint causing error: %s", uidStr, parseErr)
			return "", err
		}
		uid := uint32(uid64)
		rtcToken, err = rtctokenbuilder.BuildTokenWithUID(appId, appCertificate, channelName, uid, role, expireTimestamp)
		return rtcToken, err
	} else {
		err = fmt.Errorf("failed to generate RTC token for unknown tokenType: %s", tokenType)
		log.Println(err)
		return "", err
	}
}
