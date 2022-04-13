package handler

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/defipod/mochi/pkg/config"
	"github.com/defipod/mochi/pkg/discordwallet"
	"github.com/defipod/mochi/pkg/entities"
	"github.com/defipod/mochi/pkg/logger"
	processor "github.com/defipod/mochi/pkg/processor"
	"github.com/defipod/mochi/pkg/repo"
	"github.com/defipod/mochi/pkg/repo/pg"
	responseConverter "github.com/defipod/mochi/pkg/util/response_converter"

	"github.com/gin-gonic/gin"
)

// Handler for app
type Handler struct {
	cfg       config.Config
	repo      *repo.Repo
	dcwallet  discordwallet.IDiscordWallet
	processor processor.Processor
	entities  *entities.Entity
	discord   *discordgo.Session
}

// New will return an instance of Auth struct
func New(cfg config.Config, l logger.Log, s repo.Store, dcwallet *discordwallet.DiscordWallet) (*Handler, error) {
	r := pg.NewRepo(s.DB())
	processor := processor.NewProcessor(cfg)

	discord, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		cfg:       cfg,
		repo:      r,
		dcwallet:  dcwallet,
		processor: processor,
		entities:  entities.New(cfg, l, s, dcwallet, discord),
		discord:   discord,
	}

	return handler, nil
}

// Healthz handler
// Return "OK"
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "OK")
}

func (h *Handler) Guilds(c *gin.Context) {
	guilds, err := h.repo.DiscordGuilds.Gets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseConverter.ConvertGetGuildsResponse(guilds))
}
