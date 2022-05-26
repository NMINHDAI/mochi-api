package userwallet

import (
	"github.com/defipod/mochi/pkg/model"
	"github.com/defipod/mochi/pkg/request"
)

type Store interface {
	GetOneByDiscordIDAndGuildID(discordID, guildID string) (*model.UserWallet, error)
	GetOneByGuildIDAndAddress(guildID, address string) (*model.UserWallet, error)
	GetList(req request.GetListUserWallet) ([]model.UserWallet, error)
	UpsertOne(uw model.UserWallet) error
}
