package response

import "github.com/defipod/mochi/pkg/model"

type GuildHolderRolesResponse struct {
	model.GuildConfigHolderRole
	RoleName string `json:"role_name"`
	Color    int    `json:"color"`
}
