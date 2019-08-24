package user

import (
	"fmt"
	"fruit_shop/handler"
	"fruit_shop/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Create(c *gin.Context) {
	var u CreateRequest

	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	name := c.Param("username")
	log.Infof("URL username: %s", name)

	desc := c.Query("desc")
	log.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type: %s", contentType)

	log.Debugf("username is: [%s], password is: [%s]", u.Username, u.Password)
	if u.Username == "" {
		handler.SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not be found in the db,")), nil)
		return
	}

	if u.Password == "" {
		handler.SendResponse(c, fmt.Errorf("password is empty"), nil)
	}

	resp := CreateResponse{
		Username: u.Username,
	}

	handler.SendResponse(c, nil, resp)
}
