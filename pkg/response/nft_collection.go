package response

type NFTTradingVolume struct {
	CollectionAddress string  `json:"collection_address"`
	CollectionName    string  `json:"collection_name"`
	Symbol            string  `json:"collection_symbol"`
	ChainID           float64 `json:"collection_chain_id"`
	TradingVolume     float64 `json:"trading_volume"`
	Token             string  `json:"token"`
}

type NFTTradingVolumeResponse struct {
	Data []NFTTradingVolume `json:"data"`
}