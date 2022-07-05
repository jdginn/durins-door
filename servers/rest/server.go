package main

import (
	"debug/dwarf"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

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
	type serverContext struct {
		debugFile string
		reader    *dwarf.Reader
	}
	s := serverContext{}

	f, err := os.OpenFile("durinsdoor.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)

	server := gin.Default()

	server.POST("/debugfile", func(c *gin.Context) {
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
			c.IndentedJSON(http.StatusInternalServerError, msg)
			return
		}
		reader, err := parser.GetReader(fh)
		if err != nil {
			msg := fmt.Sprintf("Error getting reader from file %s:\n\t%s", s.debugFile, err)
			c.IndentedJSON(http.StatusInternalServerError, msg)
			return
		}
		s.reader = reader

		c.JSON(http.StatusOK, "looks good to me!")
    return
	})
	server.GET("/debugfile", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, s.debugFile)
    return
	})
	server.GET("/cus", func(c *gin.Context) {
		cus, err := parser.GetCUs(s.reader)
		if err != nil {
      fmt.Println(err)
      gin.DefaultWriter.Write([]byte(err.Error()))
			// c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("ListCUs failed with error: %s", err.Error()))
			c.JSON(http.StatusInternalServerError, "bad list CUs")
		}
    ret := make([]string, len(cus))
    for i, c := range cus {
      ret[i] = c.Val(dwarf.AttrName).(string)
    }
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
