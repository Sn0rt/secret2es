package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/Sn0rt/secret2es/pkg/converter"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/convert", func(c *gin.Context) {
		var request struct {
			Content        string `json:"content" binding:"required"`
			StoreType      string `json:"storeType" binding:"required"`
			StoreName      string `json:"storeName" binding:"required"`
			CreationPolicy string `json:"creationPolicy" binding:"required"`
			Resolve        bool   `json:"resolve"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := converter.ConvertSecretContent(
			request.Content,
			request.StoreType,
			request.StoreName,
			esv1beta1.ExternalSecretCreationPolicy(request.CreationPolicy),
			request.Resolve,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}