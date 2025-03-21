package http

import (
	"context"
	"reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h handler) test(c *gin.Context) {

	ctx := context.Background()
	connnect, _ := h.uc.Create(ctx, "mongodb+srv://root:123@cluster0.vmsjm.mongodb.net/book_store?retryWrites=true&w=majority&appName=Cluster0")
	response.OK(c, connnect)
}
