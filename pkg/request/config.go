package request

import (
	"fmt"

	"github.com/defipod/mochi/pkg/model"
)

type UpsertGmConfigRequest struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

type UpsertGuildTokenConfigRequest struct {
	GuildID string `json:"guild_id"`
	Symbol  string `json:"symbol"`
	Active  bool   `json:"active"`
}

type ConfigLevelRoleRequest struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
	Level   int    `json:"level"`
}

type ConfigHolderRoleRequest struct {
	model.GuildConfigNFT
	ID             string `json:"id"`
	NumberOfTokens int    `json:"number_of_tokens"`
	RoleID         string `json:"role_id"`
}

func (cfg ConfigHolderRoleRequest) Validate() error {
	if cfg.GuildID == "" {
		return fmt.Errorf("guild_id is required")
	}
	if cfg.RoleID == "" {
		return fmt.Errorf("role_id is required")
	}
	if cfg.TokenAddress == "" {
		return fmt.Errorf("invalid token address")
	}
	if cfg.ChainID != 1 && cfg.ChainID != 56 && cfg.ChainID != 250 {
		return fmt.Errorf("unsupported chain_id")
	}
	if cfg.ERCFormat != 721 && cfg.ERCFormat != 1155 {
		return fmt.Errorf("unsupported erc_format")
	}
	if cfg.ERCFormat == 1155 && cfg.TokenID == 0 {
		return fmt.Errorf("token id is required for erc_format 1155")
	}
	return nil
}
