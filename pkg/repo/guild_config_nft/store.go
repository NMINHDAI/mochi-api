package guild_config_nft

import "github.com/defipod/mochi/pkg/model"

type Store interface {
	ListByGuildID(guildID string) ([]model.GuildConfigNFT, error)
	UpsertOne(config *model.GuildConfigNFT) error
}
