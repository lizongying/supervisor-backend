package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
	"supervisor/common"
)

type requestJson struct {
	Group  string `json:"group"`
	Name   string `json:"name"`
	Server string `json:"server"`
}

var ErrorCode = 1
var SuccessCode = 0

var Supervisor = map[string]*common.SupervisorRpc{}

func main() {
	conf()
	for _, supervisor := range common.Config.SupervisorList {
		Supervisor[supervisor.Name] = common.Rpc(supervisor.Url)
	}
	r := gin.Default()
	r.Use(Cors())
	r.StaticFile("/", "./assets/index.html")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Static("/static", "./assets/static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
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
		list := make([]common.ProcessInfo, 0)
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
	server := common.Config.Server.Url
	err := r.Run(server)
	if err != nil {
		fmt.Println(err)
	}
}

func conf() {
	configPathDefault, _ := os.Getwd()
	configPathDefault = path.Join(configPathDefault, "conf", "example.yml")
	configPath := flag.String("c", configPathDefault, "config file path")
	flag.Parse()
	err := common.LoadConfig(*configPath)
	if err != nil {
		return
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
