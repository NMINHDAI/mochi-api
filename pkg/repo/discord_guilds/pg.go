package discord_guilds

import (
	"github.com/defipod/mochi/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pg struct {
	db *gorm.DB
}

func NewPG(db *gorm.DB) Store {
	return &pg{db: db}
}

func (pg *pg) Gets() ([]*model.DiscordGuild, error) {
	var guilds []*model.DiscordGuild
	return guilds, pg.db.Preload("GuildConfigInviteTracker").Find(&guilds).Error
}

func (pg *pg) CreateIfNotExists(guild model.DiscordGuild) error {
	tx := pg.db.Begin()
	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoNothing: true,
	}).Create(&guild).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (pg *pg) GetByID(id string) (*model.DiscordGuild, error) {
	var guild model.DiscordGuild
	return &guild, pg.db.Preload("GuildConfigInviteTracker").First(&guild, "id = ?", id).Error
}