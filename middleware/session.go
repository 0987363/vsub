package middleware

import (
	"time"

	"github.com/0987363/vsub/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

const cookieKey = "token"

func ConnectSession(key string) {
	store = sessions.NewCookieStore([]byte(key))
}

func GetToken(c *gin.Context) (*models.Token, error) {
	s, err := store.Get(c.Request, cookieKey)
	if err != nil {
		return nil, err
	}

	if s.IsNew {
		return nil, nil
	}

	token := models.Token{
		Version: s.Values["version"].(string),
		UserID:  s.Values["user_id"].(string),
		Nonce:   s.Values["nonce"].(string),
		Hmac:    s.Values["hmac"].(string),
	}

	exp := s.Values["expiry"].(string)
	token.Expiry, err = time.Parse(time.RFC3339Nano, exp)
	if err != nil {
		return nil, models.Error("Expiry is invalid.")
	}

	return &token, nil
}

func SetToken(c *gin.Context, token *models.Token) error {
	s, err := store.Get(c.Request, cookieKey)
	if err != nil {
		return err
	}

	s.Values["version"] = token.Version
	s.Values["user_id"] = token.UserID
	s.Values["nonce"] = token.Nonce
	s.Values["hmac"] = token.Hmac
	s.Values["expiry"] = token.Expiry.Format(time.RFC3339Nano)

	s.Save(c.Request, c.Writer)

	return nil
}
