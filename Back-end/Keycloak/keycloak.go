package keycloak

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type ClientKeycloak struct {
	kc *gocloak.GoCloak

	clientID     string
	clientSecret string
	realm        string

	userAdmin  string
	pwdAdmin   string
	realmAdmin string
}

type JWT struct {
	AccessToken string
	RefresToken string
	ExpiredIn   int
}

func NewClientKeycloak() *ClientKeycloak {
	return &ClientKeycloak{
		kc:           gocloak.NewClient(os.Getenv("KC")),
		clientID:     os.Getenv("CLIENTID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
		realm:        os.Getenv("REALM"),

		userAdmin:  os.Getenv("KC_BOOTSTRAP_ADMIN_USERNAME"),
		pwdAdmin:   os.Getenv("KC_BOOTSTRAP_ADMIN_PASSWORD"),
		realmAdmin: os.Getenv("REALM_ADMIN"),
	}
}

func (c *ClientKeycloak) Login(ctx context.Context, email, password string) (*JWT, error) {
	jwt, err := c.kc.Login(ctx, c.clientID, c.clientSecret, c.realm, email, password)
	if err != nil {
		return nil, err
	}
	cJWT := JWT{
		AccessToken: jwt.AccessToken,
		RefresToken: jwt.RefreshToken,
		ExpiredIn:   jwt.ExpiresIn,
	}
	return &cJWT, nil
}

type CreateUserParams struct {
	Username string
	Name     string
	Lastname string
	Email    string
	Password string
}

func (c *ClientKeycloak) CreateUser(ctx context.Context, params CreateUserParams) (userId string, error error) {
	//login in superusuario

	jwt, err := c.kc.LoginAdmin(ctx, c.userAdmin, c.pwdAdmin, c.realmAdmin)

	if err != nil {
		return "", err
	}

	newUser := gocloak.User{
		Username:  gocloak.StringP(params.Email), // usar el email como username
		Email:     gocloak.StringP(params.Email),
		FirstName: gocloak.StringP(params.Name),
		LastName:  gocloak.StringP(params.Lastname),
		Enabled:   gocloak.BoolP(true),
	}

	userID, err := c.kc.CreateUser(ctx, jwt.AccessToken, c.realm, newUser)
	if err != nil {
		return "", err
	}
	err = c.kc.SetPassword(ctx, jwt.AccessToken, userID, c.realm, params.Password, false)
	if err != nil {
		return "", err
	}
	return userID, nil
}

type UserInfo struct {
	Username string
	Email    string
}

func (c *ClientKeycloak) UserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	user, err := c.kc.GetUserInfo(ctx, accessToken, c.realm)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("User info not found")
	}

	userInfo := &UserInfo{}

	if user.Nickname != nil {
		userInfo.Username = *user.Nickname
	}
	if user.Email != nil {
		userInfo.Email = *user.Email
	}
	return userInfo, nil
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Enabled  bool   `json:"enabled"`
}

type CreateGroupParams struct {
	Name string
}

func generateNumbers() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%s", n)
}

func (c *ClientKeycloak) CreateGroup(ctx context.Context, params CreateGroupParams, accessToken string) (groupID string, error error) {
	jwt, err := c.kc.LoginAdmin(ctx, c.userAdmin, c.pwdAdmin, c.realmAdmin)
	if err != nil {
		return "", err
	}
	generateNameGroup := fmt.Sprintf("%s-%s", params.Name, generateNumbers())

	newGroup := gocloak.Group{
		Name: gocloak.StringP(generateNameGroup),
		Attributes: &map[string][]string{
			"displayName": {params.Name},
		},
	}

	groupCreateID, err := c.kc.CreateGroup(ctx, jwt.AccessToken, c.realm, newGroup)

	if err != nil {
		return "", err
	}

	userInfo, err := c.kc.GetUserInfo(ctx, accessToken, c.realm)

	if err != nil {
		return "", err
	}

	if userInfo.Sub == nil {
		return "", fmt.Errorf("userInfo.Sub is nil")
	}

	err = c.kc.AddUserToGroup(ctx, jwt.AccessToken, c.realm, *userInfo.Sub, groupCreateID)

	if err != nil {
		return "", err
	}

	return groupCreateID, nil
}

func (c *ClientKeycloak) DeleteGroup(ctx context.Context, groupID string) error {

	jwt, err := c.kc.LoginAdmin(ctx, c.userAdmin, c.pwdAdmin, c.realmAdmin)
	if err != nil {
		return err
	}

	err = c.kc.DeleteGroup(ctx, jwt.AccessToken, c.realm, groupID)

	if err != nil{
		return err
	}

	return nil
}

func (c *ClientKeycloak) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {

	token, err := c.kc.RefreshToken(
		ctx,
		refreshToken,
		c.clientID,
		c.clientSecret,
		c.realm,
	)

	if err != nil {
		return nil, err
	}
	return token, nil
}
