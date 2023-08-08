package plugin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/extrame/goblet"
	"github.com/extrame/jose/crypto"
	"github.com/extrame/jose/jws"
)

//New create a new LoginAsJwt plugin, secret is the secret key for jwt, idKey is the key for id in jwt
func JWT() *_JwtLoginPlugin {
	return &_JwtLoginPlugin{}
}

type _JwtLoginPlugin struct {
	Secret string
	secret []byte
	method crypto.SigningMethod
	Issuer string
	Alg    string
}

func (j *_JwtLoginPlugin) AddCfgAndInit(server *goblet.Server) error {

	server.AddConfig("jwt", j)

	m := jws.GetSigningMethod(j.Alg)
	if m == nil {
		return errors.New("NOT VALID SIGNING METHOD:" + j.Alg)
	}

	j.method = m
	j.secret = []byte(j.Secret)

	return nil
}

func (l *_JwtLoginPlugin) AddLoginAs(ctx *goblet.Context, name string, id string, timeduration ...time.Duration) {
	var claims = make(jws.Claims)
	claims.Set(name, id)
	j := jws.NewJWT(claims, l.method)
	j.Claims().SetIssuer(l.Issuer)
	var d = 24 * time.Hour
	if len(timeduration) > 0 {
		d = timeduration[0]
	}

	j.Claims().SetExpiration(time.Now().Add(d))
	b, err := j.Serialize(l.secret)
	if err == nil {
		ctx.SetHeader("Authorization", fmt.Sprintf("Bearer %s", string(b)))
	}
}

func (l *_JwtLoginPlugin) GetLoginIdAs(ctx *goblet.Context, key string) (string, error) {
	auth := ctx.ReqHeader().Get("Authorization")
	if auth != "" {
		auth = strings.TrimPrefix(auth, "Bearer ")
		token, err := jws.ParseJWT([]byte(auth))
		if err == nil {
			err = token.Validate(l.secret)
			if err == nil {
				return token.Claims().Get(key).(string), nil
			}
		}
	}
	return "", errors.New("NOT VALID LOGIN INFO:" + auth)
}

func (l *_JwtLoginPlugin) DeleteLoginAs(ctx *goblet.Context, key string) error {
	ctx.SetHeader("Authorization", "")
	return nil
}