package entities

import (
	"fmt"

	"github.com/defipod/mochi/pkg/model"
	"github.com/defipod/mochi/pkg/request"
	"github.com/defipod/mochi/pkg/response"
	"github.com/google/uuid"
)

func (e *Entity) GetGmConfig(guildID string) (*model.GuildConfigGmGn, error) {
	config, err := e.repo.GuildConfigGmGn.GetByGuildID(guildID)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (e *Entity) UpsertGmConfig(req request.UpsertGmConfigRequest) error {
	if err := e.repo.GuildConfigGmGn.UpsertOne(&model.GuildConfigGmGn{
		GuildID:   req.GuildID,
		ChannelID: req.ChannelID,
	}); err != nil {
		return err
	}

	return nil
}

func (e *Entity) GetGuildTokens(guildID string) ([]model.GuildConfigToken, error) {
	guildTokens, err := e.repo.GuildConfigToken.GetByGuildID(guildID)
	if err != nil {
		return nil, err
	}

	return guildTokens, nil
}

func (e *Entity) UpsertGuildTokenConfig(req request.UpsertGuildTokenConfigRequest) error {
	token, err := e.repo.Token.GetBySymbol(req.Symbol, true)
	if err != nil {
		return err
	}

	if err := e.repo.GuildConfigToken.UpsertMany([]model.GuildConfigToken{{
		GuildID: req.GuildID,
		TokenID: token.ID,
		Active:  req.Active,
	}}); err != nil {
		return err
	}

	return nil
}

func (e *Entity) ConfigLevelRole(req request.ConfigLevelRoleRequest) error {
	return e.repo.GuildConfigLevelRole.UpsertOne(model.GuildConfigLevelRole{
		GuildID: req.GuildID,
		RoleID:  req.RoleID,
		Level:   req.Level,
	})
}

func (e *Entity) GetGuildLevelRoleConfigs(guildID string) ([]model.GuildConfigLevelRole, error) {
	return e.repo.GuildConfigLevelRole.GetByGuildID(guildID)
}

func (e *Entity) GetGuildHolderRoleConfigs(guildID string) ([]model.GuildConfigHolderRole, error) {
	return e.repo.GuildConfigHolderRole.ListByGuildID(guildID)
}

func (e *Entity) GetUserRoleByLevel(guildID string, level int) (string, error) {
	config, err := e.repo.GuildConfigLevelRole.GetHighest(guildID, level)
	if err != nil {
		return "", err
	}

	return config.RoleID, nil
}

func (e *Entity) RemoveGuildMemberRoles(guildID string, rolesToRemove map[string]string) error {
	for userID, roleID := range rolesToRemove {
		if err := e.discord.GuildMemberRoleRemove(guildID, userID, roleID); err != nil {
			return err
		}
	}

	return nil
}

func (e *Entity) AddGuildMemberRoles(guildID string, rolesToAdd map[string]string) error {
	for userID, roleID := range rolesToAdd {
		if err := e.discord.GuildMemberRoleAdd(guildID, userID, roleID); err != nil {
			return err
		}
	}

	return nil
}

func (e *Entity) RemoveGuildMemberRole(guildID, userID, roleID string) error {
	return e.discord.GuildMemberRoleRemove(guildID, userID, roleID)
}

func (e *Entity) ListGuildNFTConfigs(guildID string) ([]model.GuildConfigNFT, error) {
	return e.repo.GuildConfigNFT.ListByGuildID(guildID)
}

func (e *Entity) NewGuildHolderRoleConfig(req request.ConfigHolderRoleRequest) (*model.GuildConfigHolderRole, error) {
	err := e.repo.GuildConfigNFT.UpsertOne(&req.GuildConfigNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert guild nft config: %v", err.Error())
	}

	role := model.GuildConfigHolderRole{
		ConfigNFTID:    req.GuildConfigNFT.ID,
		NumberOfTokens: req.NumberOfTokens,
		RoleID:         req.RoleID,
	}

	err = e.repo.GuildConfigHolderRole.UpsertOne(&role)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert guild config holder role: %v", err.Error())
	}

	return &role, nil
}

func (e *Entity) EditGuildHolderRoleConfig(req request.ConfigHolderRoleRequest) error {

	err := e.repo.GuildConfigNFT.UpsertOne(&req.GuildConfigNFT)
	if err != nil {
		return fmt.Errorf("failed to upsert guild nft config: %v", err.Error())
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return fmt.Errorf("failed to parse id: %v", err.Error())
	}

	role := model.GuildConfigHolderRole{
		ID:             uuid.NullUUID{id, true},
		ConfigNFTID:    req.GuildConfigNFT.ID,
		NumberOfTokens: req.NumberOfTokens,
		RoleID:         req.RoleID,
	}

	err = e.repo.GuildConfigHolderRole.Update(&role)
	if err != nil {
		return fmt.Errorf("failed to update guild config holder role: %v", err.Error())
	}

	return nil
}

func (e *Entity) RemoveGuildHolderRoleConfig(id string) error {
	err := e.repo.GuildConfigHolderRole.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to remove guild holder role config")
	}
	return nil
}

func (e *Entity) ListGuildHolderRoles(guildID string) ([]response.GuildHolderRolesResponse, error) {
	roles, err := e.repo.GuildConfigHolderRole.ListByGuildID(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list guild holder roles: %v", err.Error())
	}

	dr, err := e.discord.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list discord guild roles: %v", err.Error())
	}

	res := make([]response.GuildHolderRolesResponse, len(roles))

	for i, role := range roles {
		roleResp := response.GuildHolderRolesResponse{
			GuildConfigHolderRole: role,
		}
		for _, r := range dr {
			if role.RoleID == r.ID {
				roleResp.RoleName = r.Name
				roleResp.Color = r.Color
			}
		}
		res[i] = roleResp
	}

	return res, nil
}
