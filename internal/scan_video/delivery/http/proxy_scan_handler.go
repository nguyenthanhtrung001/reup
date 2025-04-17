package http

import (
	"reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h handler) DoneAllProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.uc.DoneAllProxyScan(ctx)
	if err != nil {
		h.l.Error(ctx, "error done all proxy scan", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	response.OK(c, "done all proxy scan")

}

// DoneProxyScan implements Handler.
func (h handler) DoneProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		ProxyIP string `json:"proxy_ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.l.Error(ctx, "error binding proxy data", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	err := h.uc.DoneProxyScan(ctx, request.ProxyIP)
	if err != nil {
		h.l.Error(ctx, "error done proxy scan", err)
		response.ErrorWithMap(c, errDoneProxyScan, nil)
		return
	}
	response.OK(c, "done proxy scan")

}

// GetAllProxyScan implements Handler.
func (h handler) GetAllProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	proxyScans, err := h.uc.GetAllProxyScan(ctx)
	if err != nil {
		h.l.Error(ctx, "error get all proxy scan", err)
		response.ErrorWithMap(c, errGetAllProxyScan, nil)
		return
	}
	var proxyIPs []string
	for _, scan := range proxyScans {
		proxyIPs = append(proxyIPs, scan.ProxyIP)
	}

	response.OK(c, proxyIPs)

}

func (h handler) GetProxyScanRandom(c *gin.Context) {
	ctx := c.Request.Context()
	proxyScan, err := h.uc.GetProxyScanRandom(ctx)
	if err != nil {
		h.l.Error(ctx, "error get proxy scan random", err)
		response.ErrorWithMap(c, errGetProxyScanRandom, nil)
		return
	}
	response.OK(c, proxyScan.ProxyIP)
}

func (h handler) InsertProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		ProxyIP string `json:"proxy_ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.l.Error(ctx, "error binding proxy data", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	err := h.uc.InsertProxyScan(ctx, request.ProxyIP)
	if err != nil {
		h.l.Error(ctx, "error insert proxy scan", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	response.OK(c, "insert proxy scan")
}
