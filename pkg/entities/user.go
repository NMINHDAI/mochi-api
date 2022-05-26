package entities

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/defipod/mochi/pkg/contracts/erc1155"
	"github.com/defipod/mochi/pkg/contracts/erc721"
	"github.com/defipod/mochi/pkg/model"
	"github.com/defipod/mochi/pkg/request"
	"github.com/defipod/mochi/pkg/response"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func (e *Entity) CreateUser(req request.CreateUserRequest) error {

	user := &model.User{
		ID:       req.ID,
		Username: req.Username,
		GuildUsers: []*model.GuildUser{
			{
				GuildID:   req.GuildID,
				UserID:    req.ID,
				Nickname:  req.Nickname,
				InvitedBy: req.InvitedBy,
			},
		},
	}

	if err := e.repo.Users.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (e *Entity) CreateUserIfNotExists(id, username string) error {
	user := &model.User{
		ID:       id,
		Username: username,
	}

	if err := e.repo.Users.FirstOrCreate(user); err != nil {
		return fmt.Errorf("failed to create if not exists user: %w", err)
	}

	return nil
}

func (e *Entity) GetUser(discordID string) (*response.GetUserResponse, error) {
	user, err := e.repo.Users.GetOne(discordID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	guildUsers := []*response.GetGuildUserResponse{}
	for _, guildUser := range user.GuildUsers {
		guildUsers = append(guildUsers, &response.GetGuildUserResponse{
			GuildID:   guildUser.GuildID,
			UserID:    guildUser.UserID,
			Nickname:  guildUser.Nickname,
			InvitedBy: guildUser.InvitedBy,
		})
	}

	if user.InDiscordWalletAddress.String == "" {
		if err = e.generateInDiscordWallet(user); err != nil {
			err = fmt.Errorf("cannot generate in-discord wallet: %v", err)
			return nil, err
		}
	}

	res := &response.GetUserResponse{
		ID:                     user.ID,
		Username:               user.Username,
		InDiscordWalletAddress: &user.InDiscordWalletAddress.String,
		InDiscordWalletNumber:  &user.InDiscordWalletNumber.Int64,
		GuildUsers:             guildUsers,
	}
	return res, nil
}

func (e *Entity) GetUserCurrentGMStreak(discordID, guildID string) (*model.DiscordUserGMStreak, int, error) {
	streak, err := e.repo.DiscordUserGMStreak.GetByDiscordIDGuildID(discordID, guildID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get user's gm streak: %v", err)
	}

	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, fmt.Errorf("user has no gm streak")
	}

	return streak, http.StatusOK, nil
}

func (e *Entity) HandleUserActivities(req *request.HandleUserActivityRequest) (*response.HandleUserActivityResponse, error) {
	userXP, err := e.repo.GuildUserXP.GetOne(req.GuildID, req.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	gConfigActivity, err := e.repo.GuildConfigActivity.GetOneByActivityName(req.GuildID, req.Action)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err := e.repo.GuildConfigActivity.ForkDefaulActivityConfigs(req.GuildID); err != nil {
			return nil, err
		}
		gConfigActivity, err = e.repo.GuildConfigActivity.GetOneByActivityName(req.GuildID, req.Action)
		if err != nil {
			return nil, fmt.Errorf("failed to get guild config activity: %v", err.Error())
		}
	}

	if err := e.repo.GuildUserActivityLog.CreateOne(model.GuildUserActivityLog{
		GuildID:      req.GuildID,
		UserID:       req.UserID,
		ActivityName: gConfigActivity.Activity.Name,
		EarnedXP:     gConfigActivity.Activity.XP,
		CreatedAt:    req.Timestamp,
	}); err != nil {
		return nil, err
	}

	latestUserXP, err := e.repo.GuildUserXP.GetOne(req.GuildID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &response.HandleUserActivityResponse{
		GuildID:      req.GuildID,
		UserID:       req.UserID,
		Action:       gConfigActivity.Activity.Name,
		AddedXP:      gConfigActivity.Activity.XP,
		CurrentXP:    latestUserXP.TotalXP,
		CurrentLevel: latestUserXP.Level,
		Timestamp:    req.Timestamp,
		LevelUp:      latestUserXP.Level > userXP.Level,
	}, nil
}

func (e *Entity) InitGuildDefaultActivityConfigs(guildID string) error {
	activities, err := e.repo.Activity.GetDefaultActivities()
	if err != nil {
		return err
	}

	var configs []model.GuildConfigActivity
	for _, a := range activities {
		configs = append(configs, model.GuildConfigActivity{
			GuildID:    guildID,
			ActivityID: a.ID,
			Active:     true,
		})
	}

	return e.repo.GuildConfigActivity.UpsertMany(configs)
}

func (e *Entity) GetTopUsers(guildID, userID string, limit, page int) (*response.GetTopUsersResponse, error) {
	offset := page * limit
	leaderboard, err := e.repo.GuildUserXP.GetTopUsers(guildID, limit, offset)
	if err != nil {
		return nil, err
	}

	author, err := e.repo.GuildUserXP.GetOne(guildID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &response.GetTopUsersResponse{
		Author:      author,
		Leaderboard: leaderboard,
	}, nil
}

func (e *Entity) GetGuildUserXPs(guildID string) ([]model.GuildUserXP, error) {
	return e.repo.GuildUserXP.GetByGuildID(guildID)
}

func (e *Entity) GetGuildMember(guildID, userID string) (*discordgo.Member, error) {
	member, err := e.discord.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (e *Entity) GetUserProfile(guildID, userID string) (*response.GetUserProfileResponse, error) {
	gUserXP, err := e.repo.GuildUserXP.GetOne(guildID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	currentLevel, err := e.repo.ConfigXPLevel.GetNextLevel(gUserXP.TotalXP, false)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	nextLevel, err := e.repo.ConfigXPLevel.GetNextLevel(gUserXP.TotalXP, true)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &response.GetUserProfileResponse{
		ID:           userID,
		CurrentLevel: currentLevel,
		NextLevel:    nextLevel,
		GuildXP:      gUserXP.TotalXP,
		NrOfActions:  gUserXP.NrOfActions,
	}, nil
}

func (e *Entity) GetGuildMemberWallet(guildID string) ([]model.UserWallet, error) {
	uws, err := e.repo.UserWallet.GetList(request.GetListUserWallet{
		GuildID:   guildID,
		ChainType: "evm",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallet: %v", err.Error())
	}
	return uws, nil
}

func (e *Entity) ListGuildMemberHolderRole(config model.GuildConfigNFT) (map[string]string, error) {

	var rpcUrl string
	switch config.ChainID {
	case 1:
		rpcUrl = e.cfg.EthereumRPC
	case 56:
		rpcUrl = e.cfg.BscRPC
	case 250:
		rpcUrl = e.cfg.FantomRPC
	default:
		return nil, fmt.Errorf("chain id %d not supported", config.ChainID)
	}

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain client: %v", err.Error())
	}

	var balanceOf func(string) (int, error)
	switch config.ERCFormat {
	case 721:
		contract721, err := erc721.NewErc721(common.HexToAddress(config.TokenAddress), client)
		if err != nil {
			return nil, fmt.Errorf("failed to init erc721 contract: %v", err.Error())
		}

		balanceOf = func(address string) (int, error) {
			b, err := contract721.BalanceOf(nil, common.HexToAddress(address))
			if err != nil {
				return 0, fmt.Errorf("failed to get balance of %s in chain %d: %v", address, config.ChainID, err.Error())
			}
			return int(b.Int64()), nil
		}

	case 1155:
		contract1155, err := erc1155.NewErc1155(common.HexToAddress(config.TokenAddress), client)
		if err != nil {
			return nil, fmt.Errorf("failed to init erc1155 contract: %v", err.Error())
		}

		balanceOf = func(address string) (int, error) {
			b, err := contract1155.BalanceOf(nil, common.HexToAddress(address), big.NewInt(int64(config.TokenID)))
			if err != nil {
				return 0, fmt.Errorf("failed to get balance of %s in chain %d: %v", address, config.ChainID, err.Error())
			}
			return int(b.Int64()), nil
		}

	default:
		return nil, fmt.Errorf("erc format %d not supported", config.ERCFormat)
	}

	uws, err := e.repo.UserWallet.GetList(request.GetListUserWallet{
		GuildID:   config.GuildID,
		ChainType: "evm",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallet: %v", err.Error())
	}

	res := make(map[string]string)
	for _, uw := range uws {
		n, err := balanceOf(uw.Address)
		if err != nil {
			return nil, err
		}

		role, err := e.repo.GuildConfigHolderRole.GetCurrentRole(config.ID, n)
		if err != nil {
			return nil, fmt.Errorf("failed to get current role: %v", err.Error())
		}

		res[uw.UserDiscordID] = role.RoleID
	}

	return res, nil
}
