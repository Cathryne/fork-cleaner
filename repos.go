// Package forkcleaner provides functions to find and remove unused forks.
package forkcleaner

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// Delete delete the given list of forks
func Delete(
	ctx context.Context,
	client *github.Client,
	deletions []*github.Repository,
) error {
	for _, repo := range deletions {
		_, err := client.Repositories.Delete(ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Find list the forks from a given owner that could be deleted
func Find(
	ctx context.Context,
	client *github.Client,
	owner string,
	blacklist []string,
) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}
	var deletions []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, owner, opt)
		if err != nil {
			return deletions, err
		}
		for _, repo := range repos {
			if shouldDelete(repo, blacklist) {
				deletions = append(deletions, repo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return deletions, nil
}

func shouldDelete(repo *github.Repository, blacklist []string) bool {
	for _, r := range blacklist {
		if r == *repo.Name {
			return false
		}
	}
	return *repo.Fork &&
		*repo.ForksCount == 0 &&
		*repo.StargazersCount == 0 &&
		!*repo.Private &&
		time.Now().AddDate(0, -1, 0).After((*repo.UpdatedAt).Time)
}
