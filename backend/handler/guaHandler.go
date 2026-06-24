package handler

import (
	"guadb/model"
	"guadb/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GuaHandler interface {
	AddGuaHandler(c *gin.Context)
	GetGuaHandler(c *gin.Context)
}

type guaHandler struct {
	guaRepo repository.GuaRepo
}

type guaReq struct {
	guaData *model.Gua `json:"data"`
}

func NewGuaHandler(guaRepo repository.GuaRepo) GuaHandler {
	return &guaHandler{guaRepo}
}

func (g *guaHandler) AddGuaHandler(c *gin.Context) {
	var req guaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	guaid, err := g.guaRepo.AddGua(req.guaData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": guaid})
}
func (g *guaHandler) GetGuaHandler(c *gin.Context) {
	var req guaReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
