package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v28/github"
)

type FakeGithubRepositories struct {
	urls []string
}

func (f *FakeGithubRepositories) GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	return newRepositoryRelease(f.urls), nil, nil
}

func NewFakeGithub(urls []string) *Github {
	return &Github{repositories: &FakeGithubRepositories{urls: urls}}
}

func newRepositoryRelease(urls []string) *github.RepositoryRelease {
	var assets []github.ReleaseAsset
	for _, url := range urls {
		u := url
		asset := github.ReleaseAsset{BrowserDownloadURL: &u}
		assets = append(assets, asset)
	}
	return &github.RepositoryRelease{Assets: assets}
}

func TestNoLinuxUrl(t *testing.T) {
	urls := []string{"url1", "url2", "url3"}
	g := NewFakeGithub(urls)
	if _, got := g.getLatestLinuxUrl(); got != NotFound {
		t.Errorf("Linux url should not be found")
	}
}

func TestResolveLinuxUrl(t *testing.T) {
	linux_url := "url2_linux_amd64.bz2"
	urls := []string{"url1", linux_url}
	g := NewFakeGithub(urls)
	got, err := g.getLatestLinuxUrl()
	if err != nil {
		t.Errorf("Linux url resolved, error should be nil")
	}
	if got != linux_url {
		t.Errorf("Resolve incorrect linux url")
	}
}
