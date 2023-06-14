package reddit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/fetch"

	"golang.org/x/time/rate"
)

const (
	Comment   Kind = "t1"
	Account   Kind = "t2"
	Link      Kind = "t3"
	Message   Kind = "t4"
	Subreddit Kind = "t5"
	Award     Kind = "t6"
	Listing   Kind = "Listing"
)

type Kind string

type Thing struct {
	Kind Kind `json:"kind"`
	Data struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Domain     string  `json:"domain"`
		Subreddit  string  `json:"subreddit"`
		CreatedUTC float64 `json:"created_utc"`
		Over18     bool    `json:"over_18,omitempty"`
		Title      string  `json:"title,omitempty"`
		Ups        int     `json:"ups,omitempty"`
		Downs      int     `json:"downs,omitempty"`
		LinkFlair  string  `json:"link_flair_text,omitempty"`
		URL        string  `json:"url,omitempty"`
		Preview    *struct {
			Images []struct {
				Source struct {
					URL    string `json:"url"`
					Width  int    `json:"width,omitempty"`
					Height int    `json:"height,omitempty"`
				} `json:"source"`
			} `json:"images,omitempty"`
		} `json:"preview,omitempty"`
	}
}

type APIResponse struct {
	Kind Kind `json:"kind"`
	Data struct {
		Children []Thing `json:"children,omitempty"`
		Before   string  `json:"before,omitempty"`
		After    string  `json:"after,omitempty"`
	}
}

type Client struct {
	appID       string
	appSecret   string
	userAgent   string
	rateLimiter *rate.Limiter

	tokenLock sync.Mutex
	token     *fetch.OAuthToken
}

func NewClient(appID, appSecret, userAgent string) *Client {
	return &Client{
		appID:       appID,
		appSecret:   appSecret,
		userAgent:   userAgent,
		rateLimiter: rate.NewLimiter(2, 2),
	}
}

func (c *Client) Subreddit(r, after string) (*APIResponse, error) {
	authHeader, err := c.resolveAuthHeader()
	if err != nil {
		return nil, err
	}
	var (
		url = fmt.Sprintf("https://www.reddit.com/r/%s/%s.json", r, "top") // TODO: parameterize
		res APIResponse
	)
	return &res, blob.Unmarshal(json.Unmarshal, fetch.Get(url).
		Query("limit", 100).
		Query("after", after).
		Query("t", "all"). // TODO: parameterize
		Query("raw_json", 1).
		Authorization(authHeader).
		UserAgent(c.userAgent).
		Limit(c.rateLimiter), &res)
}

func (c *Client) resolveAuthHeader() (string, error) {
	c.tokenLock.Lock()
	defer c.tokenLock.Unlock()

	if c.token == nil {
		if err := blob.Unmarshal(json.Unmarshal,
			fetch.Post("https://www.reddit.com/api/v1/access_token").
				Form(url.Values{
					"grant_type": []string{"client_credentials"},
				}).
				User(url.UserPassword(c.appID, c.appSecret)).
				UserAgent(c.userAgent).
				Limit(c.rateLimiter),
			&c.token,
		); err != nil {
			return "", fmt.Errorf("reddit.resolveAuthHeader: %w", err)
		}
	}
	return c.token.Header()
}
