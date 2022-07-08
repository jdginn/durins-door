package main

import (
	"debug/dwarf"
  "flag"
	"fmt"
	// "io"
	"net/http"
	// "os"
  "time"

	"github.com/gin-gonic/gin"
  "github.com/golang/glog"
  "github.com/szuecs/gin-glog"

	"github.com/jdginn/durins-door/explorer/plat"
	"github.com/jdginn/durins-door/parser"
)

func listTypeDefs(c *gin.Context) {}

func getTypeDefByName(c *gin.Context) {}

func listTypeDefChildren(c *gin.Context) {}

func getTypeDefChild(c *gin.Context) {}

func getVariable(c *gin.Context) {}

func getVariableTypeDef(c *gin.Context) {}

func setVariableValue(c *gin.Context) {}

func main() {
  flag.Parse() 

	type serverContext struct {
		debugFile string
		reader    *dwarf.Reader
	}
	s := serverContext{}

	// f, err := os.OpenFile("durinsdoor.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// gin.DisableConsoleColor()
	// gin.DefaultWriter = io.MultiWriter(f)

	server := gin.New()
  server.Use(ginglog.Logger(3 * time.Second))
  server.Use(gin.Recovery())

  glog.Info("Launching durins-door REST server")

  server.GET("/ping", func(c *gin.Context) {
    glog.Infof("Responding to ping")
    c.IndentedJSON(http.StatusOK, "pong")
  })
	server.POST("/debugfile", func(c *gin.Context) {
    glog.Infof("Attempting to set debugfile") 
		type filename struct {
			Path string `json:"path"`
		}
		var message filename
		if err := c.BindJSON(&message); err != nil {
			return
		}
		s.debugFile = message.Path

		fh, err := plat.GetReaderFromFile(s.debugFile)
		if err != nil {
			msg := fmt.Sprintf("Error getting opening file %s:\n\t%s", s.debugFile, err)
      glog.Error(msg)
			c.IndentedJSON(http.StatusInternalServerError, msg)
			return
		}
		reader, err := parser.GetReader(fh)
		if err != nil {
			msg := fmt.Sprintf("Error getting reader from file %s:\n\t%s", s.debugFile, err)
      glog.Error(msg)
			c.IndentedJSON(http.StatusInternalServerError, msg)
			return
		}
		s.reader = reader

		c.JSON(http.StatusOK, "looks good to me!")
    glog.Infof("Successfully set debugfile to %s", message.Path)
    return
	})
	server.GET("/debugfile", func(c *gin.Context) {
    glog.Infof("Returning debugfile path: %s", s.debugFile)
		c.IndentedJSON(http.StatusOK, s.debugFile)
    return
	})
	server.GET("/cus", func(c *gin.Context) {
		cus, err := parser.GetCUs(s.reader)
		if err != nil {
      msg := fmt.Sprintf("/cus failed with error: %s", err.Error())
      glog.Error(msg) 
			c.IndentedJSON(http.StatusInternalServerError, msg)
		}
    ret := make([]string, len(cus))
    for i, c := range cus {
      ret[i] = c.Val(dwarf.AttrName).(string)
    }
    glog.Infof("Successfully found CUs: %v", ret)
    c.IndentedJSON(http.StatusOK, ret)
		return
	})
	server.GET("/typedefs", func(c *gin.Context) {
	})
	server.GET("/typedefs/:type", getTypeDefByName)
	server.GET("/typedefs/:type/children", listTypeDefChildren)
	server.GET("/typedefs/:type/children:child", getTypeDefChild)
	server.GET("/variables/:variable", getVariable)
	server.GET("/variables/:variable/typdedef", getVariableTypeDef)
	server.PUT("/variables/:variable", setVariableValue)

	server.Run("localhost:8080")
}
