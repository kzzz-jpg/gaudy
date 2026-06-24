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
	GuaData *model.Gua `json:"data"`
}

func NewGuaHandler(guaRepo repository.GuaRepo) GuaHandler {
	return &guaHandler{guaRepo}
}

func (g *guaHandler) AddGuaHandler(c *gin.Context) {
	var req guaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	guaid, err := g.guaRepo.AddGua(req.GuaData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"guaId": guaid}})
}
func (g *guaHandler) GetGuaHandler(c *gin.Context) {
	var req guaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	guas := make([]*model.Gua, 0)
	guas, err := g.guaRepo.GetGua(req.GuaData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": guas})
}
