package guild_config_hodler_role

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

func (pg *pg) GetCurrentRole(configNFTID, numberOfTokens int) (*model.GuildConfigHolderRole, error) {
	config := &model.GuildConfigHolderRole{}
	return config, pg.db.Where("config_nft_id = ? AND number_of_tokens <= ?", configNFTID, numberOfTokens).Order("number_of_tokens desc").First(config).Error
}

func (pg *pg) ListByGuildID(guildID string) ([]model.GuildConfigHolderRole, error) {
	var configs []model.GuildConfigHolderRole
	return configs, pg.db.
		Joins("JOIN guild_config_nfts ON guild_config_nfts.id = guild_config_holder_roles.config_nft_id").
		Where("guild_id = ?", guildID).Find(&configs).Error
}

func (pg *pg) UpsertOne(config *model.GuildConfigHolderRole) error {
	tx := pg.db.Begin()

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "config_nft_id"},
			{Name: "number_of_tokens"},
		},
		UpdateAll: true,
	}).Create(config).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (pg *pg) Update(config *model.GuildConfigHolderRole) error {
	return pg.db.Save(config).Error
}

func (pg *pg) Delete(id string) error {
	return pg.db.Delete(&model.GuildConfigHolderRole{}, "id = ?", id).Error
}
