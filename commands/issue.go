package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type IssueCommand struct {
	Ui ui.Ui
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	searchParser.Usage = "issue [options]"
	searchParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	if _, err := searchParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	conf, err := config.NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := gitlab.GitlabRemote(c.Ui, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := gitlab.GitlabClient(c.Ui, gitlabRemote, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	fmt.Println(searchOptions.AllProject)
	if searchOptions.AllProject {
		issues, err := getIssues(client)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		for _, issue := range issues {
			data := strings.Join([]string{
				fmt.Sprintf("#%d", issue.IID),
				issue.WebURL,
				issue.Title,
			}, "|")
			datas = append(datas, data)
		}

	} else {
		issues, err := getProjectIssues(client, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		for _, issue := range issues {
			data := strings.Join([]string{
				fmt.Sprintf("#%d", issue.IID),
				issue.Title,
			}, "|")
			datas = append(datas, data)
		}

		result := columnize.SimpleFormat(datas)
		c.Ui.Message(result)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Message(result)

	return ExitCodeOK
}

func getIssues(client *gitlabc.Client) ([]*gitlabc.Issue, error) {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: searchOptions.Line,
	}
	listIssuesOptions := &gitlabc.ListIssuesOptions{
		State:       gitlabc.String(searchOptions.State),
		Scope:       gitlabc.String(searchOptions.Scope),
		OrderBy:     gitlabc.String(searchOptions.OrderBy),
		Sort:        gitlabc.String(searchOptions.Sort),
		ListOptions: *listOption,
	}

	issues, _, err := client.Issues.ListIssues(
		listIssuesOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}

	return issues, nil
}

func getProjectIssues(client *gitlabc.Client, repositoryName string) ([]*gitlabc.Issue, error) {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: searchOptions.Line,
	}
	listProjectIssuesOptions := &gitlabc.ListProjectIssuesOptions{
		State:       gitlabc.String(searchOptions.State),
		Scope:       gitlabc.String(searchOptions.Scope),
		OrderBy:     gitlabc.String(searchOptions.OrderBy),
		Sort:        gitlabc.String(searchOptions.Sort),
		ListOptions: *listOption,
	}

	issues, _, err := client.Issues.ListProjectIssues(
		repositoryName,
		listProjectIssuesOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}

	return issues, nil
}
