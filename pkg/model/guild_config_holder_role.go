package model

import "github.com/google/uuid"

type GuildConfigHolderRole struct {
	ID             uuid.NullUUID `json:"id" gorm:"default:uuid_generate_v4()"`
	ConfigNFTID    int           `json:"config_nft_id"`
	NumberOfTokens int           `json:"number_of_tokens"`
	RoleID         string        `json:"role_id"`
}
