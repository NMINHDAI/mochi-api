package guild_user_xp

import (
	"github.com/defipod/mochi/pkg/model"
	"gorm.io/gorm"
)

type pg struct {
	db *gorm.DB
}

func NewPG(db *gorm.DB) Store {
	return &pg{db: db}
}

func (pg *pg) GetOne(guildID, userID string) (*model.GuildUserXP, error) {
	userXP := &model.GuildUserXP{}
	return userXP, pg.db.Where("guild_id = ? AND user_id = ?", guildID, userID).First(userXP).Error
}

func (pg *pg) GetTopUsers(guildID string, limit, offset int) ([]model.GuildUserXP, error) {
	var userXPs []model.GuildUserXP
	return userXPs, pg.db.Where("guild_id = ?", guildID).Preload("User").Offset(offset).Limit(limit).Order("guild_rank").Find(&userXPs).Error
}