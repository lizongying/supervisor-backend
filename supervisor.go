package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"supervisor/app"
)

type requestJson struct {
	Group  string `json:"group"`
	Name   string `json:"name"`
	Server string `json:"server"`
}

var ErrorCode = 1
var SuccessCode = 0

var Supervisor = map[string]*app.SupervisorRpc{}

func main() {
	app.InitConfig()
	server := app.Config.Server
	gin.SetMode(server.Mode)
	for _, supervisor := range app.Config.SupervisorList {
		Supervisor[supervisor.Name] = app.Rpc(supervisor.Url)
	}
	r := gin.New()
	r.Use(cors.Default())
	r.StaticFile("/", "./dist/index.html")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")
	r.Static("/static", "./dist/static")
	r.GET("/api/supervisor/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": app.Config,
		})
	})
	r.POST("/api/supervisor/stop", func(c *gin.Context) {
		var RequestJson requestJson
		err := c.BindJSON(&RequestJson)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		if _, ok := Supervisor[RequestJson.Server]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		_, err = Supervisor[RequestJson.Server].StopProcessGroup(RequestJson.Group)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res, err := Supervisor[RequestJson.Server].GetProcessInfo(RequestJson.Group, RequestJson.Name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res.Server = RequestJson.Server
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": res,
		})
	})
	r.POST("/api/supervisor/start", func(c *gin.Context) {
		var RequestJson requestJson
		err := c.BindJSON(&RequestJson)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		if _, ok := Supervisor[RequestJson.Server]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		_, err = Supervisor[RequestJson.Server].StartProcessGroup(RequestJson.Group)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res, err := Supervisor[RequestJson.Server].GetProcessInfo(RequestJson.Group, RequestJson.Name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res.Server = RequestJson.Server
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": res,
		})
	})
	r.POST("/api/supervisor/restart", func(c *gin.Context) {
		var RequestJson requestJson
		err := c.BindJSON(&RequestJson)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		if _, ok := Supervisor[RequestJson.Server]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		_, err = Supervisor[RequestJson.Server].StopProcessGroup(RequestJson.Group)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		_, err = Supervisor[RequestJson.Server].StartProcessGroup(RequestJson.Group)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res, err := Supervisor[RequestJson.Server].GetProcessInfo(RequestJson.Group, RequestJson.Name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res.Server = RequestJson.Server
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": res,
		})
	})
	r.POST("/api/supervisor/status", func(c *gin.Context) {
		var RequestJson requestJson
		err := c.BindJSON(&RequestJson)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		if _, ok := Supervisor[RequestJson.Server]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res, err := Supervisor[RequestJson.Server].GetProcessInfo(RequestJson.Group, RequestJson.Name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		res.Server = RequestJson.Server
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": res,
		})
	})
	r.GET("/api/supervisor/list", func(c *gin.Context) {
		list := make([]app.ProcessInfo, 0)
		for server, item := range Supervisor {
			tempMap := map[string]byte{}
			ret, _ := item.GetAllProcessInfo()
			for _, info := range ret {
				l := len(tempMap)
				tempMap[info.Group] = 0
				if len(tempMap) == l {
					continue
				}
				info.Server = server
				list = append(list, info)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": SuccessCode,
			"data": list,
		})
	})
	if err := r.Run(server.Url); err != nil {
		log.Fatalln(err)
	}
}
