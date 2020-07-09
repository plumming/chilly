package domain

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/plumming/chilly/pkg/util"

	"github.com/ghodss/yaml"

	"github.com/plumming/chilly/pkg/pr"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/chilly/pkg/cmd"
)

var (
	defaultRepos = []string{"plumming/chilly"}
)

// GetPrs defines get pull request response.
type GetPrs struct {
	cmd.CommonOptions
	ShowDependabot bool
	ShowOnHold     bool
	PullRequests   []pr.PullRequest
}

// Repos defines repos to watch.
type Repos struct {
	Repos []string `json:"repos"`
}

// Data.
type Data struct {
	Search Search `json:"search"`
}

// Search.
type Search struct {
	PullRequests []pr.PullRequest `json:"nodes"`
}

// NewGetPrs.
func NewGetPrs() *GetPrs {
	g := &GetPrs{}
	return g
}

// Validate input.
func (g *GetPrs) Validate() error {
	return nil
}

// Run the cmd.
func (g *GetPrs) Run() error {
	client, err := g.GithubClient()
	if err != nil {
		return err
	}

	query := `{
	search(query: "is:pr is:open %s", type: ISSUE, first: 100) {
      nodes {
      ... on PullRequest {
        number
        title
        url
		createdAt
        closed
        author {
          login
        }
        repository {
          nameWithOwner
        }
        mergeable
        labels(first: 10) {
          nodes {
            name
          }
        }
        commits(last: 1){
          nodes{
            commit{
              status {
                contexts {
                  state
                  targetUrl
                  description
                  context
                }
              }
            }
          }
        }
      }
    }
  }
}`

	data := Data{}
	repos, err := repos()
	if err != nil {
		return err
	}
	err = client.GraphQL(fmt.Sprintf(query, strings.Join(repos, " ")), nil, &data)
	if err != nil {
		return err
	}

	pulls := data.Search.PullRequests
	sort.Sort(pr.ByPullsString(pulls))

	pullsToReturn := []pr.PullRequest{}

	for _, pr := range pulls {
		pullRequest := pr
		if pr.Display(g.ShowDependabot, g.ShowOnHold) {
			pullsToReturn = append(pullsToReturn, pullRequest)
		}
	}

	g.PullRequests = pullsToReturn

	return nil
}

func repos() ([]string, error) {
	configFile := util.ChillyConfigFile()

	var repos []string
	if exists, err := util.FileExists(configFile); err == nil && exists {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		configRepos := Repos{}
		err = yaml.Unmarshal(content, &configRepos)
		if err != nil {
			log.Logger().Infof("no repos configured in %s", configFile)
			os.Exit(1)
		}
		repos = configRepos.Repos
		if len(repos) == 0 {
			log.Logger().Infof("no repos configured in %s", configFile)
			os.Exit(1)
		}
	} else if err != nil {
		return nil, err
	} else {
		repos = defaultRepos
	}

	var repoList []string
	for _, r := range repos {
		repoList = append(repoList, fmt.Sprintf("repo:%s", r))
	}
	return repoList, nil
}

// Retrigger failed prs.
func (g *GetPrs) Retrigger() error {
	client, err := g.GithubClient()
	if err != nil {
		return err
	}

	log.Logger().Infof("Retriggering Failed & Non Conflicting PRs...")

	for _, pr := range g.PullRequests {
		if pr.ContextsString() == "FAILURE" && pr.Mergeable == "MERGEABLE" && pr.HasLabel("updatebot") {
			failedContexts := pr.FailedContexts
			for _, f := range failedContexts() {
				testCommand := fmt.Sprintf("/test %s", f.Context)
				if f.Context == "pr-build" {
					testCommand = "/test this"
				}
				log.Logger().Infof("%s with '%s'", pr.URL, testCommand)

				url := fmt.Sprintf("repos/%s/issues/%d/comments", pr.Repository.NameWithOwner, pr.Number)
				body := fmt.Sprintf("{ \"body\": \"%s\" }", testCommand)

				err := client.REST("POST", url, strings.NewReader(body), nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
