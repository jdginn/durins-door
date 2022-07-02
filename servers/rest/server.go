package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	//TODO: is explorer the right tool for the job here?
	"github.com/jdginn/durins-door/explorer"
)

func listTypeDefs(c *gin.Context) {}

func getTypeDefByName(c *gin.Context) {}

func listTypeDefChildren(c *gin.Context) {}

func getTypeDefChild(c *gin.Context) {}

func getVariable(c *gin.Context) {}

func getVariableTypeDef(c *gin.Context) {}

func setVariableValue(c *gin.Context) {}

func main() {
	server := gin.Default()

	e := explorer.NewExplorer()
	server.POST("/debugfile", func(c *gin.Context) {
		type filename struct {
			Path string `json:"path"`
		}
		var message filename
		if err := c.BindJSON(&message); err != nil {
			return
		}
    if err := e.CreateReaderFromFile(message.Path); err != nil {
    }
		c.String(http.StatusAccepted, "")
	})
	server.GET("/debugfile", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, e.DwarfFile)
	})
	// server.GET("/typedefs", func(c *gin.Context) {
	//  })
	server.GET("/typedefs/:type", getTypeDefByName)
	server.GET("/typedefs/:type/children", listTypeDefChildren)
	server.GET("/typedefs/:type/children:child", getTypeDefChild)
	server.GET("/variables/:variable", getVariable)
	server.GET("/variables/:variable/typdedef", getVariableTypeDef)
	server.PUT("/variables/:variable", setVariableValue)

	server.Run("localhost:8080")
}
