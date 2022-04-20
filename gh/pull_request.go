/*
Copyright © 2022 Madhav Jivrajani

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gh

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

const (
	PriorityCriticalUrgent    = "priority/critical-urgent"
	PriorityImportantSoon     = "priority/important-soon"
	PriorityImportantLongterm = "priority/important-longterm"
)

const k8s = "kubernetes"

type RelevantPullRequests struct {
	PRBeingReviewed *github.PullRequest
	PRs             []*github.PullRequest
	NumFilesMatched map[int]int
}

func GetRelevantPullRequests(ctx context.Context, prNum int) (*RelevantPullRequests, error) {
	client, err := newGithubClient(ctx)
	if err != nil {
		return nil, err
	}

	allPRs := []*github.PullRequest{}
	opt := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 30},
	}
	for {
		prs, resp, err := client.PullRequests.List(ctx, k8s, k8s, opt)
		if err != nil {
			return nil, err
		}
		allPRs = append(allPRs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	filteredPRs := filterPRs(allPRs)
	targetPRFiles, err := listFiles(ctx, client, prNum)
	if err != nil {
		return nil, err
	}

	fileMap := map[string]struct{}{}
	for _, file := range targetPRFiles {
		fileMap[*file.Filename] = struct{}{}
		if file.PreviousFilename != nil {
			fileMap[*file.PreviousFilename] = struct{}{}
		}
	}

	result := &RelevantPullRequests{
		PRs:             []*github.PullRequest{},
		NumFilesMatched: make(map[int]int),
	}

	for _, pr := range filteredPRs {
		if *pr.Number == prNum {
			result.PRBeingReviewed = pr
			continue
		}
		files, err := listFiles(ctx, client, *pr.Number)
		if err != nil {
			return nil, err
		}
		atleastOneMatch := false
		for _, file := range files {
			if _, ok := fileMap[*file.Filename]; ok {
				atleastOneMatch = true
				result.NumFilesMatched[*pr.Number]++
			}
		}
		if atleastOneMatch {
			result.PRs = append(result.PRs, pr)
		}
	}

	return result, nil
}

func newGithubClient(ctx context.Context) (*github.Client, error) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return nil, fmt.Errorf("unable to look GITHUB_TOKEN env variable")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func filterPRs(prs []*github.PullRequest) []*github.PullRequest {
	res := []*github.PullRequest{}
	for _, pr := range prs {
		for _, label := range pr.Labels {
			if *label.Name == PriorityCriticalUrgent ||
				*label.Name == PriorityImportantSoon ||
				*label.Name == PriorityImportantLongterm {
				res = append(res, pr)
			}
		}
	}

	return res
}

func listFiles(ctx context.Context, client *github.Client, prNum int) ([]*github.CommitFile, error) {
	opt := &github.ListOptions{
		PerPage: 30,
	}
	allFiles := []*github.CommitFile{}
	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, k8s, k8s, prNum, opt)
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, files...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allFiles, nil
}
