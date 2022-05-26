package model

type GuildConfigNFT struct {
	ID           int    `json:"-" gorm:"primary_key"`
	GuildID      string `json:"guild_id"`
	TokenName    string `json:"token_name"`
	TokenAddress string `json:"token_address"`
	ChainID      int    `json:"chain_id"`
	ERCFormat    int    `json:"erc_format"`
	TokenID      int    `json:"token_id"`
}
