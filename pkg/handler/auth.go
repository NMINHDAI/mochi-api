package handler

import (
	"net/http"

	"github.com/defipod/mochi/pkg/consts"
	"github.com/defipod/mochi/pkg/request"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {

	var req request.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_token is required"})
		return
	}

	resp, err := h.entities.Login(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(consts.TokenCookieKey, resp.AccessToken, int(resp.ExpiresAt), "/", "", true, true)
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(consts.TokenCookieKey, "", -1, "/", "", false, true)

	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "logged out",
	})
}
