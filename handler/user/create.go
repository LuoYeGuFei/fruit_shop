package user

import (
	"fmt"
	"fruit_shop/pkg/errno"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Create(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json: "password"`
	}

	var err error
	if err := c.Bind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": errno.ErrBind})
		return
	}

	log.Debugf("username is: [%s], password is: [%s]", u.Username, u.Password)
	if u.Username == "" {
		err = errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not be found in the db,")).Add("This is the additional message")
		log.Errorf(err, "Get an error")
	}

	if errno.IsErrUserNotFound(err) {
		log.Debug("err type is ErrUserNotFound")
	}

	if u.Password == "" {
		err = fmt.Errorf("password is empty")
	}

	code, message := errno.DecodeErr(err)
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
}
