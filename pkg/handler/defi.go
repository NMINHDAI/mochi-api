package handler

import (
	"fmt"
	"net/http"

	"github.com/defipod/mochi/pkg/request"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHistoricalMarketChart(c *gin.Context) {
	data, err, statusCode := h.entities.GetHistoricalMarketChart(c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) InDiscordWalletTransfer(c *gin.Context) {
	var req request.TransferRequest
	if err := req.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, errs := h.entities.InDiscordWalletTransfer(req)
	if errs != nil {
		for _, err := range errs {
			fmt.Println("error transfer in dcwallet:", err)
		}
	}

	if res == nil || len(res) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errs})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   res,
		"errors": errs,
	})
}

func (h *Handler) InDiscordWalletWithdraw(c *gin.Context) {
	var req request.TransferRequest
	if err := req.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.entities.InDiscordWalletWithdraw(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) InDiscordWalletBalances(c *gin.Context) {
	query := c.Request.URL.Query()
	discordID := query.Get("discord_id")
	username := query.Get("username")

	response, err := h.entities.InDiscordWalletBalances(discordID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": response})
}

func (h *Handler) GetSupportedTokens(c *gin.Context) {
	tokens, err := h.entities.GetSupportedTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
	})
}

func (h *Handler) GetCoin(c *gin.Context) {
	data, err, statusCode := h.entities.GetCoinData(c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}