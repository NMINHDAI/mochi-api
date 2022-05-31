package handler

import (
	"net/http"
	"strconv"

	"github.com/defipod/mochi/pkg/request"
	"github.com/gin-gonic/gin"
)

func(h *Handler) GiftXpHandler(c *gin.Context) {
	var req request.GiftXpRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//  TODO: validate this admin discord id, check if this id has role admin ?
	// userGuild, err := h.entities.GetUserDiscord(req.GuildId, req.AdminDiscordId)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "cannot get discord admin"})
	// 	return
	// }
	
	_, err := h.entities.GetUser(req.UserDiscordId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	earnedXp, _ := strconv.Atoi(req.XpAmount)
	err = h.entities.CreateGuildUserActivityLog(req.GuildId, req.UserDiscordId, earnedXp, "gifted")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create activity log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}