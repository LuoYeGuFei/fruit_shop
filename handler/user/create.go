package user

import (
	"fruit_shop/handler"
	"fruit_shop/model"
	"fruit_shop/pkg/errno"
	"fruit_shop/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Create create a new user
func Create(c *gin.Context) {
	log.Info("User create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest

	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data
	if err := u.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password
	if err := u.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// Insert the user to the database
	if err := u.Create(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information
	handler.SendResponse(c, nil, rsp)

}
