package user

import (
	"fruit_shop/handler"
	"fruit_shop/model"
	"fruit_shop/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Get get an user by username
func Get(c *gin.Context) {
	username := c.Param("username")
	user, err := model.GetUser(username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	handler.SendResponse(c, nil, user)
}
