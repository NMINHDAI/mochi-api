package request

type CreateDefaultRoleRequest struct {
	RoleID  string `json:"role_id"`
	GuildID string `json:"guild_id"`
}
