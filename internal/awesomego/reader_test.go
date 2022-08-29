package awesomego

import (
	"github.com/stretchr/testify/assert"
	"sorted-awesome-go/pkg/cache"
	"testing"
	"time"
)

func TestReader_parsePackageList(t *testing.T) {
	zeroTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	type githubMockResponse struct {
		stars      uint
		forks      uint
		updatedAt  time.Time
		fromGithub bool
	}
	tests := []struct {
		name                   string
		input                  string
		expectedGitHubPackages []Package
		expectedOtherPackages  []Package
		gitHubClientResponse   githubMockResponse
	}{
		{
			name:  "package with github link",
			input: "- [foo](https://github.com/foo/bar) - baz",
			expectedGitHubPackages: []Package{
				{
					Name:        "foo",
					URL:         "https://github.com/foo/bar",
					Stars:       5,
					Forks:       7,
					UpdatedAt:   zeroTime,
					Description: "baz",
				},
			},
			gitHubClientResponse: githubMockResponse{
				5,
				7,
				zeroTime,
				true,
			},
			expectedOtherPackages: nil,
		},
		{
			name:                   "package with link",
			input:                  "- [foo](https://example.local) - baz",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo",
					URL:         "https://example.local",
					Description: "baz",
				},
			},
		},
		{
			name:                   "package without link",
			input:                  "- [foo] - baz",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo",
					URL:         "",
					Description: "baz",
				},
			},
		},
		{
			name:                   "package without link but with special chars in description",
			input:                  "- [foo-bar-baz] - lorem (ipsum dolor [sit amet])",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo-bar-baz",
					URL:         "",
					Description: "lorem (ipsum dolor [sit amet])",
				},
			},
		},
		{
			name:                   "package with special chars in name",
			input:                  "- [foo (bar)](https://example.local) - lorem (ipsum dolor)",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo (bar)",
					URL:         "https://example.local",
					Description: "lorem (ipsum dolor)",
				},
			},
		},
		{
			name:                   "package with link in description",
			input:                  "- [foo](https://example.local) - lorem [ipsum](https://dolor.local) sit amet",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo",
					URL:         "https://example.local",
					Description: "lorem [ipsum](https://dolor.local) sit amet",
				},
			},
		},
		{
			name:                   "package with link in description and without package link",
			input:                  "- [foo] - lorem [ipsum](https://dolor.local) sit amet",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo",
					URL:         "",
					Description: "lorem [ipsum](https://dolor.local) sit amet",
				},
			},
		},
		{
			name:                   "missing separator between link & description",
			input:                  "- [foo (bar)](https://example.local) lorem (ipsum dolor)",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo (bar)",
					URL:         "https://example.local",
					Description: "lorem (ipsum dolor)",
				},
			},
		},
		{
			name:                   "missing separator between & description",
			input:                  "- [foo] lorem (ipsum dolor)",
			expectedGitHubPackages: nil,
			expectedOtherPackages: []Package{
				{
					Name:        "foo",
					URL:         "",
					Description: "lorem (ipsum dolor)",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewReader(&cache.RuntimeCache{}, NewGitHubMockClient(
				test.gitHubClientResponse.stars,
				test.gitHubClientResponse.forks,
				test.gitHubClientResponse.updatedAt,
				test.gitHubClientResponse.fromGithub,
			))

			githubPackages, otherPackages := reader.parsePackageList(test.input)

			assert.Equal(t, test.expectedGitHubPackages, githubPackages)
			assert.Equal(t, test.expectedOtherPackages, otherPackages)
		})
	}
}

type GitHubMockClient struct {
	stars      uint
	forks      uint
	updatedAt  time.Time
	fromGithub bool
}

func NewGitHubMockClient(stars uint, forks uint, updatedAt time.Time, fromGithub bool) *GitHubMockClient {
	return &GitHubMockClient{stars: stars, forks: forks, updatedAt: updatedAt, fromGithub: fromGithub}
}

func (c *GitHubMockClient) GetDetails(_ string) (stars uint, forks uint, updatedAt time.Time, fromGithub bool) {
	return c.stars, c.forks, c.updatedAt, c.fromGithub
}
