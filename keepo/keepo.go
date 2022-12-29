package keepo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		ID             string      `json:"id"`
		ShortID        string      `json:"short_id"`
		CreatedAt      time.Time   `json:"created_at"`
		ParentIds      interface{} `json:"parent_ids"`
		Title          string      `json:"title"`
		Message        string      `json:"message"`
		AuthorName     string      `json:"author_name"`
		AuthorEmail    string      `json:"author_email"`
		AuthoredDate   time.Time   `json:"authored_date"`
		CommitterName  string      `json:"committer_name"`
		CommitterEmail string      `json:"committer_email"`
		CommittedDate  time.Time   `json:"committed_date"`
		Trailers       interface{} `json:"trailers"`
		WebURL         string      `json:"web_url"`
	} `json:"commit"`
	Merged             bool   `json:"merged"`
	Protected          bool   `json:"protected"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
	CanPush            bool   `json:"can_push"`
	Default            bool   `json:"default"`
	WebURL             string `json:"web_url"`
}

type Results struct {
	Status       string   `json:"status"`
	TotalResults int      `json:"totalResults"`
	Branches     []Branch `json:"branches"`
}

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

func (c *Client) FetchEverything(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://keepo.mit.id/api/v4/projects/158/repository/branches", url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
