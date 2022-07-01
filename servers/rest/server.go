package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"github.com/jdginn/durins-door/explorer"
)

type RESTClient struct {
  e explorer.Explorer
}

func setDebugFile(c *gin.Context) {}

func getDebugFile(c *gin.Context) {}

func listTypeDefs(c *gin.Context) {}

func getTypeDefByName(c *gin.Context) {}

func listTypeDefChildren(c *gin.Context) {}

func getTypeDefChild(c *gin.Context) {}

func getVariable(c *gin.Context) {}

func getVariableTypeDef(c *gin.Context) {}

func setVariableValue(c *gin.Context) {}

func main() {
  server := gin.Default()
  server.PUT("/debugfile", setDebugFile)
  server.GET("/debugfile", getDebugFile)
  server.GET("/typedefs", listTypeDefs)
  server.GET("/typedefs/:type", getTypeDefByName)
  server.GET("/typedefs/:type/children", listTypeDefChildren)
  server.GET("/typedefs/:type/children:child", getTypeDefChild)
  server.GET("/variables/:variable", getVariable)
  server.GET("/variables/:variable/typdedef", getVariableTypeDef)
  server.PUT("/variables/:variable", setVariableValue)
}
