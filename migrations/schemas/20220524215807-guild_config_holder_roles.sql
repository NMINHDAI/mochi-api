
-- +migrate Up
CREATE TABLE IF NOT EXISTS guild_config_nfts (
	id serial PRIMARY KEY,
	guild_id TEXT NOT NULL REFERENCES discord_guilds(id),
	token_name TEXT,
	token_address TEXT NOT NULL,
	chain_id INTEGER NOT NULL,
	erc_format INTEGER NOT NULL,
	token_id INTEGER,
	UNIQUE (guild_id, token_address, chain_id)
);

CREATE TABLE IF NOT EXISTS guild_config_holder_roles (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	config_nft_id INTEGER NOT NULL REFERENCES guild_config_nfts(id),
	number_of_tokens INTEGER NOT NULL,
	role_id TEXT NOT NULL,
	UNIQUE (config_nft_id, number_of_tokens)
);

-- +migrate Down
DROP TABLE IF EXISTS guild_config_holder_roles;
DROP TABLE IF EXISTS guild_config_nfts;