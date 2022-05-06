package entities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/defipod/mochi/pkg/model"
	"github.com/defipod/mochi/pkg/request"
	"github.com/defipod/mochi/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (e *Entity) NewGuildConfigWalletVerificationMessage(req model.GuildConfigWalletVerificationMessage) error {

	_, err := e.repo.DiscordGuilds.GetByID(req.GuildID)
	if err != nil {
		return fmt.Errorf("failed to get discord guild: %v", err.Error())
	}

	_, err = e.repo.GuildConfigWalletVerificationMessage.GetOne(req.GuildID)
	switch err {
	case nil:
		return fmt.Errorf("this guild already have a verification config")
	case gorm.ErrRecordNotFound:
	default:
		return fmt.Errorf("failed to get guild config verification: %v", err.Error())
	}

	if err := e.repo.GuildConfigWalletVerificationMessage.UpsertOne(req); err != nil {
		return fmt.Errorf("failed to upsert guild config verification: %v", err.Error())
	}

	var embeddedMsg discordgo.MessageEmbed

	err = json.Unmarshal([]byte(req.EmbeddedMessage), &embeddedMsg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal embedded message %v: %v", req.EmbeddedMessage, err.Error())
	}

	_, err = e.discord.ChannelMessageSendComplex(req.VerifyChannelID, &discordgo.MessageSend{
		Content: req.Content,
		Embed:   &embeddedMsg,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Verify",
						Style:    discordgo.PrimaryButton,
						CustomID: "mochi_verify",
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err.Error())
	}

	return nil
}

func (e *Entity) GenerateVerification(req request.GenerateVerificationRequest) (data string, statusCode int, err error) {

	_, err = e.repo.GuildConfigWalletVerificationMessage.GetOne(req.GuildID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", http.StatusBadRequest, fmt.Errorf("this guild has not set verification config")
		}
		return "", http.StatusInternalServerError, fmt.Errorf("failed to get guild config verification: %v", err.Error())
	}

	uw, err := e.repo.UserWallet.GetOneByDiscordIDAndGuildID(req.UserDiscordID, req.GuildID)
	switch err {
	case nil:
		if !req.IsReverify {
			return uw.Address, http.StatusBadRequest, fmt.Errorf("already have a verified wallet")
		}
	case gorm.ErrRecordNotFound:
		if req.IsReverify {
			return "", http.StatusBadRequest, fmt.Errorf("unverified user")
		}
	default:
		return "", http.StatusInternalServerError, err
	}

	code := uuid.New().String()

	if err := e.repo.DiscordWalletVerification.UpsertOne(
		model.DiscordWalletVerification{
			Code:          code,
			UserDiscordID: req.UserDiscordID,
			GuildID:       req.GuildID,
			CreatedAt:     time.Now(),
		},
	); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return code, http.StatusOK, nil
}

func (e *Entity) VerifyWalletAddress(req request.VerifyWalletAddressRequest) (int, error) {
	verification, err := e.repo.DiscordWalletVerification.GetByValidCode(req.Code)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid code")
	}

	if err := util.VerifySig(req.WalletAddress, req.Signature, fmt.Sprintf(
		"This will help us connect your discord account to the wallet address.\n\nMochiBotCode=%s", req.Code)); err != nil {
		return http.StatusBadRequest, err
	}

	_, err = e.repo.Users.GetOne(verification.UserDiscordID)
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		u := &model.User{
			ID: verification.UserDiscordID,
			GuildUsers: []*model.GuildUser{
				{
					GuildID: verification.GuildID,
					UserID:  verification.UserDiscordID,
				},
			},
		}
		if err := e.generateInDiscordWallet(u); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to generate in-discord wallet: %v", err.Error())
		}
	default:
		return http.StatusInternalServerError, fmt.Errorf("failed to get user: %v", err.Error())
	}

	uw, err := e.repo.UserWallet.GetOneByGuildIDAndAddress(verification.GuildID, req.WalletAddress)
	switch err {
	case nil:
		if uw.UserDiscordID != verification.UserDiscordID {
			// this address is already used by another user in this guild
			return http.StatusBadRequest, fmt.Errorf("this wallet address already belong to another user")
		}

	case gorm.ErrRecordNotFound:
		if err := e.repo.UserWallet.UpsertOne(model.UserWallet{
			UserDiscordID: verification.UserDiscordID,
			GuildID:       verification.GuildID,
			Address:       req.WalletAddress,
		}); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to upsert user wallet: %v", err.Error())
		}

	default:
		return http.StatusInternalServerError, fmt.Errorf("failed to get user wallet: %v", err.Error())
	}

	if err := e.repo.DiscordWalletVerification.DeleteByCode(verification.Code); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to delete verification: %v", err.Error())
	}

	return http.StatusOK, nil
}