package guild_config_hodler_role

import "github.com/defipod/mochi/pkg/model"

type Store interface {
	GetCurrentRole(configNFTID, numberOfTokens int) (*model.GuildConfigHolderRole, error)
	ListByGuildID(guildID string) ([]model.GuildConfigHolderRole, error)
	UpsertOne(config *model.GuildConfigHolderRole) error
	Update(config *model.GuildConfigHolderRole) error
	Delete(id string) error
}
