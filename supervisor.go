package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"supervisor/common"
)

type requestJson struct {
	Group  string `json:"group"`
	Name   string `json:"name"`
	Server string `json:"server"`
}

var ErrorCode = 1

var Supervisor = map[string]*common.SupervisorRpc{}

func main() {
	conf()
	for _, supervisor := range common.Config.SupervisorList {
		Supervisor[supervisor.Name] = common.Rpc(supervisor.Url)
	}
	r := gin.Default()
	r.StaticFile("/", "./assets/index.html")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Static("/static", "./assets/static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/get-api-version", func(c *gin.Context) {
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
		version, err := Supervisor[RequestJson.Server].GetAPIVersion()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": ErrorCode,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"version": version,
		})
	})
	r.POST("/stop-process-group", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	})
	r.POST("/start-process-group", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	})
	r.POST("/restart-process-group", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	})
	r.POST("/status-process", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	})
	r.GET("/get-all-process-info", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, list)
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
