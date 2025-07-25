package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/domain"
)

type MetaController struct {
	MetaRepository domain.MetaRepository
}

func (t MetaController) AddMeta(c *gin.Context) {
	body := &AddMetaRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	model := domain.MetaData{
		UUID:          uuid.New(),
		NumberOfPages: body.NumberOfPages,
		Height:        body.Height,
		Width:         body.Width,
		Images:        body.Images,
	}

	if err := t.MetaRepository.AddMeta(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metaUUID": model.UUID})
}

func (t MetaController) UpdateMeta(c *gin.Context) {
	body := &UpdateMetaRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	model := domain.MetaData{
		UUID:          uuid.New(),
		NumberOfPages: &body.NumberOfPages,
		Height:        &body.Height,
		Width:         &body.Width,
		Images:        &body.Images,
	}

	if err := t.MetaRepository.UpdateMeta(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (t MetaController) DeleteMeta(c *gin.Context) {
	body := &DeleteMetaRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	model := domain.MetaData{
		UUID: body.UUID,
	}

	if err := t.MetaRepository.DeleteMeta(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (t MetaController) GetMeta(c *gin.Context) {
	if id, isPresent := c.GetQuery("id"); isPresent {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			return
		}

		data, err := t.MetaRepository.GetMeta(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}

		c.JSON(http.StatusOK, data)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

func (t MetaController) SetupRouter(c *gin.RouterGroup) {
	c.GET("/", t.GetMeta)
	c.POST("/", t.AddMeta)
	c.PUT("/", t.UpdateMeta)
	c.DELETE("/", t.DeleteMeta)
}
