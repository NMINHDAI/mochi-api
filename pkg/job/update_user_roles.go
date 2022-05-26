package job

import (
	"fmt"

	"github.com/defipod/mochi/pkg/entities"
	"github.com/defipod/mochi/pkg/logger"
	"gorm.io/gorm"
)

type updateUserRoles struct {
	entity *entities.Entity
	log    logger.Logger
}

func NewUpdateUserRolesJob(e *entities.Entity, l logger.Logger) Job {
	return &updateUserRoles{
		entity: e,
		log:    l,
	}
}

func (c *updateUserRoles) Run() error {
	guilds, err := c.entity.GetGuilds()
	if err != nil {
		return err
	}

	for _, guild := range guilds.Data {
		err = c.updateLevelRoles(guild.ID)
		if err != nil {
			return err
		}

		err = c.updateHolderRoles(guild.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *updateUserRoles) updateLevelRoles(guildID string) error {
	c.log.Infof("start updating users role - guild %s", guildID)
	lrConfigs, err := c.entity.GetGuildLevelRoleConfigs(guildID)
	if err != nil {
		return err
	}

	if len(lrConfigs) == 0 {
		c.log.Infof("no levelrole configs found - guild %s", guildID)
		return nil
	}

	userXPs, err := c.entity.GetGuildUserXPs(guildID)
	if err != nil {
		return err
	}
	if len(userXPs) == 0 {
		c.log.Infof("no user XP found - guild %s", guildID)
		return nil
	}

	rolesToAdd := make(map[string]string)
	rolesToRemove := make(map[string]string)
	for _, userXP := range userXPs {
		member, err := c.entity.GetGuildMember(guildID, userXP.UserID)
		if err != nil {
			c.log.Errorf(err, "cannot get guild member %s - guild %s", userXP.UserID, guildID)
			return err
		}

		userLevelRole, err := c.entity.GetUserRoleByLevel(guildID, userXP.Level)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				c.log.Errorf(err, "cannot get role by level %d - guild %s", userXP.Level, guildID)
				return err
			}
			c.log.Infof("no config found for level %d - guild %s", userXP.Level, guildID)
			return err
		}

		memberRoles := make(map[string]bool)
		for _, roleID := range member.Roles {
			memberRoles[roleID] = true
		}

		// add role if not assigned yet
		if _, ok := memberRoles[userLevelRole]; !ok {
			rolesToAdd[userXP.UserID] = userLevelRole
		}

		for _, lrConfig := range lrConfigs {
			if _, ok := memberRoles[lrConfig.RoleID]; ok && lrConfig.RoleID != userLevelRole {
				rolesToRemove[userXP.UserID] = lrConfig.RoleID
			}
		}
	}

	if err := c.entity.RemoveGuildMemberRoles(guildID, rolesToRemove); err != nil {
		c.log.Errorf(err, "cannot remove guild member roles - guild %s", guildID)
		return err
	}

	if err := c.entity.AddGuildMemberRoles(guildID, rolesToAdd); err != nil {
		c.log.Errorf(err, "cannot add guild member roles - guild %s", guildID)
		return err
	}

	return nil
}

func (c *updateUserRoles) updateHolderRoles(guildID string) error {
	hrConfigs, err := c.entity.GetGuildHolderRoleConfigs(guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild holder role configs: %v", err.Error())
	}

	if len(hrConfigs) == 0 {
		return nil
	}

	isHolderRoles := make(map[string]bool)
	for _, hrConfig := range hrConfigs {
		isHolderRoles[hrConfig.RoleID] = true
	}

	nftConfigs, err := c.entity.ListGuildNFTConfigs(guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild nft configs: %v", err.Error())
	}

	holders, err := c.entity.GetGuildMemberWallet(guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild member wallets: %v", err.Error())
	}

	rolesToAdd := make([]map[string]string, len(nftConfigs))

	for i, config := range nftConfigs {
		roleToAdd, err := c.entity.ListGuildMemberHolderRole(config)
		if err != nil {
			return err
		}
		rolesToAdd[i] = roleToAdd
	}

	for _, holder := range holders {
		member, err := c.entity.GetGuildMember(guildID, holder.UserDiscordID)
		if err != nil {
			return fmt.Errorf("failed to get guild member: %v", err.Error())
		}

	ROLES:
		for _, roleID := range member.Roles {
			if isHolderRoles[roleID] {
				for _, roleToAdd := range rolesToAdd {
					if roleID == roleToAdd[holder.UserDiscordID] {
						delete(roleToAdd, holder.UserDiscordID)
						continue ROLES
					}
				}

				err = c.entity.RemoveGuildMemberRole(guildID, holder.UserDiscordID, roleID)
				if err != nil {
					return err
				}
			}
		}

		for _, roleToAdd := range rolesToAdd {
			err = c.entity.AddGuildMemberRoles(guildID, roleToAdd)
			if err != nil {
				return fmt.Errorf("failed to add guild member roles: %v", err.Error())
			}
		}
	}

	return nil
}
