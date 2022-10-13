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

package printer

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/MadhavJivrajani/cr/gh"
	"github.com/google/go-github/v43/github"
	"github.com/olekukonko/tablewriter"
)

func PrettyDisplayPRs(result *gh.RelevantPullRequests) {
	sort.Slice(result.PRs, func(i, j int) bool {
		var iPriority, jPriority int
		for _, label := range result.PRs[i].Labels {
			switch *label.Name {
			case gh.PriorityCriticalUrgent:
				iPriority = 3
			case gh.PriorityImportantSoon:
				iPriority = 2
			case gh.PriorityImportantLongterm:
				iPriority = 1
			}
		}
		for _, label := range result.PRs[j].Labels {
			switch *label.Name {
			case gh.PriorityCriticalUrgent:
				jPriority = 3
			case gh.PriorityImportantSoon:
				jPriority = 2
			case gh.PriorityImportantLongterm:
				jPriority = 1
			}
		}

		return jPriority < iPriority
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PR", "Title", "Created On", "Time Elapsed", "Priority Labels", "Num File Matches"})
	table.SetCaption(
		true,
		fmt.Sprintf(
			"PR Being Reviewed - %s: https://github.com/kubernetes/kubernetes/pull/%d",
			*result.PRBeingReviewed.Title,
			*result.PRBeingReviewed.Number,
		),
	)
	table.SetRowLine(true)

	for _, pr := range result.PRs {
		row := []string{
			fmt.Sprintf("%d", *pr.Number),
			*pr.Title,
			pr.CreatedAt.String(),
			time.Now().Sub(*pr.CreatedAt).String(),
			getPriorityLabelString(pr.Labels),
			fmt.Sprintf("%d", result.NumFilesMatched[*pr.Number]),
		}
		table.Append(row)
	}

	table.Render()
}

func PrettyDisplayIssues(result map[string]int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SIG", "No. of Triaged Bugs"})
	table.SetCaption(
		true,
		"Issues with triage/accepted + kind/bug broken down by SIGs",
	)
	table.SetRowLine(true)

	for sig, num := range result {
		row := []string{sig, fmt.Sprintf("%d", num)}
		table.Append(row)
	}

	table.Render()
}

func getPriorityLabelString(labels []*github.Label) string {
	for _, label := range labels {
		switch *label.Name {
		case gh.PriorityCriticalUrgent:
			return "critical-urgent"
		case gh.PriorityImportantSoon:
			return "important-soon"
		case gh.PriorityImportantLongterm:
			return "important-longterm"
		}
	}

	// We will never reach here.
	return ""
}
