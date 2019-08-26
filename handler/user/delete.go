package user

import (
	"fruit_shop/handler"
	"fruit_shop/model"
	"fruit_shop/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Delete delete user from database by id
func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUser(uint64(userId)); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}
