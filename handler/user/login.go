package user

import (
	"fruit_shop/handler"
	"fruit_shop/model"
	"fruit_shop/pkg/auth"
	"fruit_shop/pkg/errno"
	"fruit_shop/pkg/token"

	"github.com/gin-gonic/gin"
)

// Login generates the authentication token if the password was matched with the specified account.
func Login(c *gin.Context) {
	// Binding the data with the user struct
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// Get information by username
	user, err := model.GetUser(u.Username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// Compare the login password with the user password
	if err := auth.Compare(user.Password, u.Password); err != nil {
		handler.SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token
	t, err := token.Sign(c, token.Context{ID: user.Id, Username: user.Username}, "")
	if err != nil {
		handler.SendResponse(c, errno.ErrToken, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{Token: t})
}
