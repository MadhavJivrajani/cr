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
	"log"
	"strings"

	"github.com/google/go-github/v43/github"
)

func HowManyIssuesWithLabels(ctx context.Context, labels ...string) (map[string]int, error) {
	client, err := NewGithubClient(ctx)
	if err != nil {
		return nil, err
	}

	opt := &github.IssueListByRepoOptions{
		Labels:      labels,
		ListOptions: github.ListOptions{PerPage: 200},
	}
	var totalIssues int
	var issues []*github.Issue
	for {
		issuePage, resp, err := client.Issues.ListByRepo(ctx, k8s, k8s, opt)
		if err != nil {
			log.Println(resp)
			return nil, err
		}
		issues = append(issues, issuePage...)
		totalIssues += len(issues)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return groupBySIG(issues), nil
}

func groupBySIG(issues []*github.Issue) map[string]int {
	result := make(map[string]int)
	for _, issue := range issues {
		for _, label := range issue.Labels {
			if !strings.HasPrefix(*label.Name, "sig/") {
				continue
			}
			result[*label.Name]++
		}
	}

	return result
}
