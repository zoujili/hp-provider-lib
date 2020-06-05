package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"strings"
)

const Authorization = "Authorization"

type JwtOperator struct {
	token     string
	TokenType string `json:"type"`
	UserID    string `json:"user_id"`
	MapClaims jwt.MapClaims
}

func (j *JwtOperator) GetUserID() uuid.UUID {
	return uuid.FromStringOrNil(j.UserID)
}

func (j *JwtOperator) Token() string {
	return j.token
}

func NewJwtOperator(ctx context.Context) (*JwtOperator, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("No metadata from context")
	}

	auths := md.Get(Authorization)
	if auths == nil || len(auths) == 0 {
		return nil, errors.New("Authorization Missing")
	}

	operator := &JwtOperator{token: auths[0]}

	if len(operator.Token()) == 0 {
		return nil, fmt.Errorf("No authorization token from context")
	}

	tokenString, err := stripBearerPrefixFromTokenString(operator.Token())
	if err != nil {
		return nil, err
	}

	parser := jwt.Parser{SkipClaimsValidation: true}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Can't convert to JWT")
	}
	operator.MapClaims = mapClaims

	bs, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, operator)
	if err != nil {
		return nil, err
	}
	return operator, nil
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}
