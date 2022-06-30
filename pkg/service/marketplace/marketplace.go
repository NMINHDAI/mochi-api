package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/defipod/mochi/pkg/config"
)

type marketplace struct {
	config *config.Config
}

func NewMarketplace(cfg *config.Config) Service {
	return &marketplace{
		config: cfg,
	}
}

type openseaPrimaryAssetContracts struct {
	Address string `json:"address"`
}
type openseaCollection struct {
	Editors               interface{}                    `json:"editors"`
	PaymentTokens         interface{}                    `json:"payment_tokens"`
	PrimaryAssetContracts []openseaPrimaryAssetContracts `json:"primary_asset_contracts"`
	Traits                interface{}                    `json:"traits"`
}

type openseaGetCollectionResponse struct {
	Collection openseaCollection `json:"collection"`
}

func (e *marketplace) ConvertPaintswapToFtmAddress(paintswapMarketplace string) string {
	splittedPaintswap := strings.Split(paintswapMarketplace, "/")
	return splittedPaintswap[len(splittedPaintswap)-1]
}

func (e *marketplace) ConvertOpenseaToEthAddress(openseaMarketplace string) string {
	splittedOpensea := strings.Split(openseaMarketplace, "/")
	collectionSymbol := splittedOpensea[len(splittedOpensea)-1]
	openseaCollection, _ := e.GetCollectionFromOpensea(collectionSymbol)
	return openseaCollection.Collection.PrimaryAssetContracts[0].Address
}

func (e *marketplace) HandleMarketplaceLink(contractAddress, chain string) string {
	switch strings.Contains(contractAddress, "/") {
	case false:
		return contractAddress
	case true:
		switch chain {
		case "paintswap":
			return e.ConvertPaintswapToFtmAddress(contractAddress)
		case "opensea":
			return e.ConvertOpenseaToEthAddress(contractAddress)
		default:
			return e.ConvertPaintswapToFtmAddress(contractAddress)
		}
	default:
		return contractAddress
	}
}

func (e *marketplace) GetCollectionFromOpensea(collectionSymbol string) (*openseaGetCollectionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/collection/%s", e.config.MarketplaceBaseUrl.Opensea, collectionSymbol)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-API-KEY", e.config.MarketplaceApiKey.Opensea)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errBody := new(bytes.Buffer)
		_, err = errBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("openseaGetCollection - failed to read response: %v", err)
		}

		err = fmt.Errorf("GetNFTCollections - failed to get opensea collections with symbol=%s: %v", collectionSymbol, errBody.String())
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data := &openseaGetCollectionResponse{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}