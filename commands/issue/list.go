package issue

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeIssueOption(issueListOption *ListOption) *gitlab.ListIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlab.ListIssuesOptions{
		State:       gitlab.String(issueListOption.getState()),
		Scope:       gitlab.String(issueListOption.getScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		Search:      gitlab.String(issueListOption.Search),
		ListOptions: *listOption,
	}
	return listIssuesOptions
}

type listMethod struct {
	internal.Method
	client  lab.Issue
	opt     *ListOption
	project string
}

func (m *listMethod) Process() (string, error) {
	issues, err := m.client.GetProjectIssues(
		makeProjectIssueOption(m.opt),
		m.project,
	)
	if err != nil {
		return "", err
	}

	output := listOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

type listAllMethod struct {
	internal.Method
	client lab.Issue
	opt    *ListOption
}

func (m *listAllMethod) Process() (string, error) {
	issues, err := m.client.GetAllProjectIssues(makeIssueOption(m.opt))
	if err != nil {
		return "", err
	}

	output := listAllOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func makeProjectIssueOption(issueListOption *ListOption) *gitlab.ListProjectIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(issueListOption.getState()),
		Scope:       gitlab.String(issueListOption.getScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func listAllOutput(issues []*gitlab.Issue) []string {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("%s", cyan(lab.ParceRepositoryFullName(issue.WebURL))),
			fmt.Sprintf("%s", yellow(issue.IID)),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func listOutput(issues []*gitlab.Issue) []string {
	yellow := color.New(color.FgYellow).SprintFunc()
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("%s", yellow(issue.IID)),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}
