package handler

import (
	"net/http"
	"strings"

	"github.com/defipod/mochi/pkg/request"
	"github.com/defipod/mochi/pkg/util"
	"github.com/gin-gonic/gin"
)

type error interface {
	Error() string
}

func (h *Handler) GetNFTDetail(c *gin.Context) {
	symbol := c.Param("symbol")
	tokenID := c.Param("id")

	data, err := h.entities.GetNFTDetail(symbol, tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data.Image = util.StandardizeUri(data.Image)

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// TODO: add test for this api
func (h *Handler) CreateNFTCollection(c *gin.Context) {
	var req request.CreateNFTCollectionRequest

	if err := req.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.entities.CreateNFTCollection(req)
	if err != nil {
		if strings.Contains(err.Error(), "Already added") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) GetSupportedChains(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []string{"eth", "heco", "bsc", "matic", "op", "btt", "okt", "movr", "celo", "metis", "cro", "xdai", "boba", "ftm", "avax", "arb", "aurora"}})
}

func (h *Handler) ListAllNFTCollections(c *gin.Context) {
	nfts, err := h.entities.ListAllNFTCollections()
	if err != nil {
		h.log.Errorf(err, "[handler.ListAllNFTCollections] failed to get all collections, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": nfts})
}

func (h *Handler) GetNFTCollectionTickers(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	data, err := h.entities.GetNFTCollectionTickers(symbol, c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
func (h *Handler) GetNFTTradingVolume(c *gin.Context) {
	nfts, err := h.entities.GetSvc().Indexer.GetNFTTradingVolume()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nfts})
}

func (h *Handler) GetNFTCollections(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")
	if page == "" {
		page = "0"
	}
	if size == "" {
		size = "10"
	}
	data, err := h.entities.GetNFTCollections(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetNFTTokens(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	data, err := h.entities.GetNFTTokens(symbol, c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) CreateNFTSalesTracker(c *gin.Context) {
	var req request.NFTSalesTrackerRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.entities.CreateSalesTracker(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

}

func (h *Handler) GetDetailNftCollection(c *gin.Context) {
	collectionSymbol := c.Param("symbol")
	if collectionSymbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	collection, err := h.entities.GetDetailNftCollection(collectionSymbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": collection})
}
func (h *Handler) GetAllNFTSalesTracker(c *gin.Context) {
	data, err := h.entities.GetAllNFTSalesTracker()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) GetNewListedNFTCollection(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")
	interval := c.Query("interval")
	if interval == "" {
		interval = "7"
	}
	if page == "" {
		page = "0"
	}
	if size == "" {
		size = "10"
	}

	data, err := h.entities.GetNewListedNFTCollection(interval, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetNftMetadataAttrIcon(c *gin.Context) {
	data, err := h.entities.GetNftMetadataAttrIcon()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get icons"})
		return
	}

	c.JSON(http.StatusOK, data)
}
