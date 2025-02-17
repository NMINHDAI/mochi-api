package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/defipod/mochi/pkg/config"
	"github.com/defipod/mochi/pkg/logger"
	"github.com/defipod/mochi/pkg/response"
)

type Discord struct {
	session           *discordgo.Session
	log               logger.Logger
	mochiLogChannelID string
}

const (
	mochiLogColor = 0xFCD3C1
)

func NewService(
	cfg config.Config,
	log logger.Logger,
) (Service, error) {
	// *** discord ***
	discord, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to init discord: %w", err)
	}
	return &Discord{
		session:           discord,
		log:               log,
		mochiLogChannelID: cfg.MochiLogChannelID,
	}, nil
}

func (d *Discord) NotifyNewGuild(guildID string) error {
	// get new guild info
	guild, err := d.session.Guild(guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild info: %w", err)
	}

	msgEmbed := discordgo.MessageEmbed{
		Title:       "Mochi has joined new Guild!",
		Description: fmt.Sprintf("**%s** (%s)", guild.Name, guild.ID),
		Color:       mochiLogColor,
	}

	_, err = d.session.ChannelMessageSendEmbed(d.mochiLogChannelID, &msgEmbed)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (d *Discord) SendGuildActivityLogs(channelID, userID, title, description string) error {
	if channelID == "" {
		return nil
	}

	dcUser, err := d.session.User(userID)
	if err != nil {
		d.log.Errorf(err, "[SendGuildActivityLogs] - get discord user failed %s", userID)
		return err
	}
	msgEmbed := discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       mochiLogColor,
		Timestamp:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: dcUser.AvatarURL(""),
		},
	}

	_, err = d.session.ChannelMessageSendEmbed(channelID, &msgEmbed)
	if err != nil {
		return fmt.Errorf("[SendGuildActivityLogs] - ChannelMessageSendEmbed failed - channel %s: %s", channelID, err.Error())
	}

	return nil
}

func (d *Discord) SendLevelUpMessage(logChannelID, role string, uActivity *response.HandleUserActivityResponse) {
	if !uActivity.LevelUp {
		return
	}
	if uActivity.ChannelID == "" && logChannelID == "" {
		d.log.Info("Action was not performed at any channel and no log channel configured as well")
		return
	}
	channelID := logChannelID
	if channelID == "" {
		channelID = uActivity.ChannelID
	}
	if role == "" {
		role = "N/A"
	}

	dcUser, err := d.session.User(uActivity.UserID)
	if err != nil {
		d.log.Errorf(err, "SendLevelUpMessage - failed to get discord user %s", uActivity.UserID)
		return
	}

	description := fmt.Sprintf("<@%s> has leveled up **(%d - %d)**\n\n**XP: **%d\n**Role: **%s", uActivity.UserID, uActivity.CurrentLevel-1, uActivity.CurrentLevel, uActivity.CurrentXP, role)
	msgEmbed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Level up!",
			IconURL: "https://cdn.discordapp.com/emojis/984824963112513607.png?size=240&quality=lossless",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: dcUser.AvatarURL(""),
		},
		Description: description,
		Color:       mochiLogColor,
		Timestamp:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	_, err = d.session.ChannelMessageSendEmbed(channelID, &msgEmbed)
	if err != nil {
		d.log.Errorf(err, "SendLevelUpMessage - failed to send level up msg")
		return
	}
}
