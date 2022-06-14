// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/repo/discord_guild_stat_channels/store.go

// Package mock_discord_guild_stat_channels is a generated GoMock package.
package mock_discord_guild_stat_channels

import (
	reflect "reflect"

	model "github.com/defipod/mochi/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStore) Create(statChannel *model.DiscordGuildStatChannel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", statChannel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockStoreMockRecorder) Create(statChannel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStore)(nil).Create), statChannel)
}

// DeleteStatChannelByChannelID mocks base method.
func (m *MockStore) DeleteStatChannelByChannelID(channelID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStatChannelByChannelID", channelID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStatChannelByChannelID indicates an expected call of DeleteStatChannelByChannelID.
func (mr *MockStoreMockRecorder) DeleteStatChannelByChannelID(channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStatChannelByChannelID", reflect.TypeOf((*MockStore)(nil).DeleteStatChannelByChannelID), channelID)
}

// GetStatChannelsByGuildID mocks base method.
func (m *MockStore) GetStatChannelsByGuildID(guildID string) ([]model.DiscordGuildStatChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatChannelsByGuildID", guildID)
	ret0, _ := ret[0].([]model.DiscordGuildStatChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatChannelsByGuildID indicates an expected call of GetStatChannelsByGuildID.
func (mr *MockStoreMockRecorder) GetStatChannelsByGuildID(guildID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatChannelsByGuildID", reflect.TypeOf((*MockStore)(nil).GetStatChannelsByGuildID), guildID)
}