package http

import (
	"context"
	"fmt"
	"github.com/aleri-godays/frontend/internal/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
	"regexp"
	"time"
)

type githubAuth struct {
	config           *config.Config
	oauth2Config     *oauth2.Config
	oauthStateString string
	cookieName       string
	sessionTTL       time.Duration
}

func newGithubAuth(conf *config.Config) *githubAuth {
	a := &githubAuth{
		config:           conf,
		cookieName:       "frontend_jwt",
		sessionTTL:       72 * time.Hour,
		oauthStateString: "dummy-insecure-state",
		oauth2Config: &oauth2.Config{
			ClientID:     conf.OAuthClientID,
			ClientSecret: conf.OAuthClientSecret,
			Scopes:       []string{"user:email"},
			Endpoint:     githuboauth.Endpoint,
		},
	}

	return a
}

func (a *githubAuth) Login(c echo.Context) error {
	url := a.oauth2Config.AuthCodeURL(a.oauthStateString, oauth2.AccessTypeOnline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *githubAuth) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = a.cookieName
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (a *githubAuth) Callback(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)

	state := c.FormValue("state")
	if state != a.oauthStateString {
		logger.WithFields(logrus.Fields{
			"expected": a.oauthStateString,
			"got":      state,
		}).Warn("invalid oauth state")
		return EchoError(c, http.StatusUnauthorized, "could not process login")
	}

	code := c.FormValue("code")
	token, err := a.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"code":  code,
			"error": err,
		}).Error("oauth exchange failed")
		return EchoError(c, http.StatusUnauthorized, "could not process login")
	}

	oauthClient := a.oauth2Config.Client(context.Background(), token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("could not fetch user data from github")
		return EchoError(c, http.StatusUnauthorized, "could not process login")
	}

	logger.WithFields(logrus.Fields{
		"github_login": user.GetLogin(),
	}).Debug("logged in as github user")

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user"] = *user.Login
	claims["exp"] = time.Now().Add(a.sessionTTL).Unix()

	type Token struct {
		Token string
	}
	var t Token
	t.Token, err = jwtToken.SignedString([]byte(a.config.JWTSecret))
	if err != nil {
		return EchoError(c, http.StatusInternalServerError, "could not create jwt")
	}

	cookie := new(http.Cookie)
	cookie.Name = a.cookieName
	cookie.Value = t.Token
	cookie.Expires = time.Now().Add(a.sessionTTL - 15*time.Minute)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (a *githubAuth) AuthMiddleware() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		p := c.Path()

		allowed := map[string]bool{
			"/":             true,
			"/health":       true,
			"/login":        true,
			"/logout":       true,
			"/authcallback": true,
		}
		if ok, exists := allowed[p]; ok && exists {
			return true
		}

		matched, err := regexp.MatchString("^/static/", p)
		if err != nil {
			log.WithFields(log.Fields{
				"path":  p,
				"error": err,
			}).Error("regex failed")
			return false
		}

		return matched
	}
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(a.config.JWTSecret),
		Skipper:     skipper,
		TokenLookup: fmt.Sprintf("cookie:%s", a.cookieName),
	})
}
