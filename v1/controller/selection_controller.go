package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/repositories"
)

type SelectionController struct {
	SelectionRepository repositories.SelectionRepository
}

func (t SelectionController) getSelectionFromId(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
	}

	results, err := t.SelectionRepository.GetSelectionBySelectionId(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
	}

	c.JSON(200, results)
}

func (t SelectionController) SetupRouter(c *gin.RouterGroup) {
	c.GET("/", t.getSelectionFromId)
}
