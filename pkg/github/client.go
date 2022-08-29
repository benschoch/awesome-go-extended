package github

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"sorted-awesome-go/pkg/cache"
	"sorted-awesome-go/pkg/hash"
	"strconv"
	"strings"
	"time"
)

const (
	githubBaseURL = "https://github.com/"
)

type ClientInterface interface {
	GetDetails(url string) (stars uint, forks uint, updatedAt time.Time, fromGithub bool)
}

type Client struct {
	client *github.Client
	cache  cache.Handler
}

func NewClient(cache cache.Handler) ClientInterface {
	var githubUserToken string
	if os.Getenv("GITHUB_ACCESS_TOKEN") != "" {
		githubUserToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	}
	var client *github.Client
	if githubUserToken != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubUserToken})
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &Client{client: client, cache: cache}
}

func (c *Client) GetDetails(url string) (stars uint, forks uint, updatedAt time.Time, fromGithub bool) {
	user, repo, err := c.parseGithubURL(url)
	if err != nil {
		return
	}

	fromGithub = true

	starsKey := hash.Sha1(user + "/" + repo + ":stars")
	forksKey := hash.Sha1(user + "/" + repo + ":forks")
	updatedAtKey := hash.Sha1(user + "/" + repo + ":updatedAt")

	if c.cache.Has(starsKey) && c.cache.Has(forksKey) && c.cache.Has(updatedAtKey) {
		cached := c.cache.ReadString(starsKey)
		starsValue, _ := strconv.Atoi(cached)
		stars = uint(starsValue)

		cached = c.cache.ReadString(forksKey)
		forksValue, _ := strconv.Atoi(cached)
		forks = uint(forksValue)

		cached = c.cache.ReadString(updatedAtKey)
		updatedAtValue, _ := strconv.Atoi(cached)
		updatedAt = time.Unix(int64(updatedAtValue), 0)
	} else {
		repository, resp, err := c.client.Repositories.Get(context.TODO(), user, repo)
		if err != nil {
			log.Printf("error occurred while fetching repository details: %s", err)
			return
		}
		defer resp.Body.Close()

		stars = uint(*repository.StargazersCount)
		forks = uint(*repository.ForksCount)
		updatedAt = (*repository.UpdatedAt).Time

		c.cacheResult(starsKey, stars, forksKey, forks, updatedAtKey, updatedAt)
	}

	return
}

func (c *Client) cacheResult(starsKey string, stars uint, forksKey string, forks uint, updatedAtKey string, updatedAt time.Time) {
	err := c.cache.Write(starsKey, []byte(strconv.Itoa(int(stars))))
	if err != nil {
		panic(err)
	}
	err = c.cache.Write(forksKey, []byte(strconv.Itoa(int(forks))))
	if err != nil {
		panic(err)
	}
	err = c.cache.Write(updatedAtKey, []byte(strconv.Itoa(int(updatedAt.Unix()))))
	if err != nil {
		panic(err)
	}
}

func (c *Client) parseGithubURL(url string) (user string, repo string, err error) {
	if !strings.HasPrefix(url, githubBaseURL) {
		return "", "", fmt.Errorf("unsupported URL: %s", url)
	}
	url = strings.TrimPrefix(url, githubBaseURL)
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 2 {
		return "", "", fmt.Errorf("unsupported URL: %s", url)
	}
	user = strings.TrimSpace(urlParts[0])
	repo = strings.TrimSpace(urlParts[1])
	if user == "" || repo == "" {
		return "", "", errors.New("empty user or repo")
	}

	return
}
