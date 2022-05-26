package guild_config_nft

import (
	"github.com/defipod/mochi/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pg struct {
	db *gorm.DB
}

func NewPG(db *gorm.DB) *pg {
	return &pg{
		db: db,
	}
}

func (pg *pg) ListByGuildID(guildID string) ([]model.GuildConfigNFT, error) {
	var configs []model.GuildConfigNFT
	return configs, pg.db.Where("guild_id = ?", guildID).Find(&configs).Error
}

func (pg *pg) UpsertOne(config *model.GuildConfigNFT) error {
	tx := pg.db.Begin()

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "guild_id"},
			{Name: "token_address"},
			{Name: "chain_id"},
		},
		UpdateAll: true,
	}).Create(config).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
