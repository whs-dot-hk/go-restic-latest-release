package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/google/go-github/v28/github"
)

var (
	owner = "restic"
	repo  = "restic"
	expr  = "_linux_amd64\\.bz2"

	NotFound = errors.New("Restic linux url not found")
)

type GithubRepositories interface {
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
}

type Github struct {
	repositories GithubRepositories
}

func NewGithub() *Github {
	client := github.NewClient(nil)
	return &Github{repositories: client.Repositories}
}

func (g *Github) getLatestRelease() (*github.RepositoryRelease, error) {
	ctx := context.Background()
	r, _, err := g.repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Github) getLatestLinuxUrl() (string, error) {
	r, err := g.getLatestRelease()
	if err != nil {
		return "", err
	}
	re, _ := regexp.Compile(expr)
	for _, asset := range r.Assets {
		if url := asset.GetBrowserDownloadURL(); re.MatchString(url) {
			return url, nil
		}
	}
	return "", NotFound
}

func main() {
	g := NewGithub()
	url, err := g.getLatestLinuxUrl()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}
