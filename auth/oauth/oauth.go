package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/fetch"
)

const (
	BearerToken TokenType = "Bearer"
	MACToken    TokenType = "MAC"
	BasicToken  TokenType = "Basic"

	Default ScopeFormat = 0 // will use "Space"
	Space   ScopeFormat = ' '
	Comma   ScopeFormat = ','
)

type (
	Token struct {
		AccessToken  string    `json:"access_token"`
		ExpiresIn    int       `json:"expires_in"`
		ExtExpiresIn int       `json:"ext_expires_in,omitempty"` // microsoft specific
		RefreshToken string    `json:"refresh_token,omitempty"`
		Scope        string    `json:"scope"`
		TokenType    TokenType `json:"token_type"`
		IDToken      string    `json:"id_token,omitempty"`
	}
	TokenType string
	Config    struct {
		AuthBaseURL  string
		TokenBaseURL string
		ClientID     string
		ClientSecret string
		RedirectURI  string
		Scopes       []Scope
		ScopeFormat  ScopeFormat
	}
	Scope       string
	ScopeFormat byte
)

func (t Token) HTTPHeader() (string, error) {
	if t.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	switch t.TokenType {
	case BearerToken:
		return string(t.TokenType) + " " + t.AccessToken, nil
	default:
		return "", fmt.Errorf("unsupported token type: %s", t.TokenType)
	}
}

func RedirectURL(conf Config) string {
	if conf.ScopeFormat == Default {
		conf.ScopeFormat = Space
	}
	var scopeBuilder strings.Builder
	for i, s := range conf.Scopes {
		if i > 0 {
			scopeBuilder.WriteByte(byte(conf.ScopeFormat))
		}
		scopeBuilder.WriteString(string(s))
	}

	// TODO: CSRF prevention via state signatures

	/*
		stateJSON, err := json.Marshal(state)
		if err != nil {
			panic("invalid oauth state: " + err.Error())
		}
	*/
	query := url.Values{
		"client_id":     {conf.ClientID},
		"response_type": {"code"},
		"redirect_uri":  {conf.RedirectURI},
		"scope":         {scopeBuilder.String()},
		//"state":         {base64.RawURLEncoding.EncodeToString(stateJSON)},
	}
	return conf.AuthBaseURL + "?" + query.Encode()
}

func ParseCallback(conf Config, r *http.Request) (Token, error) {
	if conf.TokenBaseURL == "" {
		return Token{}, errors.New("missing config: token base URL")
	}
	if conf.ClientSecret == "" {
		return Token{}, errors.New("missing config: client secret")
	}
	q := r.URL.Query()
	if err := parseCallbackError(q); err != nil {
		return Token{}, fmt.Errorf("callback: %w", err)
	}

	// TODO: parse and validate state

	return blob.UnmarshalJSONValue[Token](conf.newTokenRequest(q.Get("code")))
}

func (conf Config) newTokenRequest(code string) *fetch.Request {
	data := url.Values{
		"client_id":     {conf.ClientID},
		"client_secret": {conf.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {conf.RedirectURI},
	}
	return fetch.Post(conf.TokenBaseURL).
		Body(blob.FromString(data.Encode()), "application/x-www-form-urlencoded")
}
