package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"text/template"
	"time"
)

type GitlabUser struct {
	Username string
	Name     string
	Email    string
}
type GitlabMilestone struct {
	DueDate string `json:"due_date"`
}
type GitlabIssue struct {
	Id          int
	Title       string
	Description string
	Author      GitlabUser
	Assignee    GitlabUser
	Milestone   GitlabMilestone
	CreatedAt   time.Time `json:"created_at"`
}
type GitlabComment struct {
	Author    GitlabUser
	Body      string
	CreatedAt time.Time `json:"created_at"`
}

type GitlabIssueReader struct {
	privateToken  string
	gitlabHost    string
	projectId     int
	current       int
	issueBuffer   []GitlabIssue
	nextPage      int
	IssueTemplate *template.Template
}

func NewGitlabIssueReader(privateToken, gitlabHost, project string) (*GitlabIssueReader, error) {
	projectId, err := getProjectId(privateToken, gitlabHost, url.QueryEscape(project))
	if err != nil {
		return nil, err
	}
	issueTemplate, err := template.ParseFiles("markdown.template")
	if err != nil {
		return nil, err
	}
	return &GitlabIssueReader{
		privateToken,
		gitlabHost,
		projectId,
		0,
		[]GitlabIssue{},
		1,
		issueTemplate,
	}, nil
}

func (r *GitlabIssueReader) fetchComments(issueId int) ([]GitlabComment, error) {
	resp, err := gitlabRequest(r.privateToken, r.gitlabHost,
		fmt.Sprintf("/projects/%d/issues/%d/notes", r.projectId, issueId))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("could not fetch issue comments: %s", resp.Status))
	}
	var comments []GitlabComment
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil

}

func (r *GitlabIssueReader) formatIssueDescription(issue *GitlabIssue, comments []GitlabComment) (string, error) {
	var templateData struct {
		Issue    *GitlabIssue
		Comments []GitlabComment
	}
	templateData.Issue = issue
	sort.Sort(ByCreationDate(comments))
	templateData.Comments = comments
	var issueText bytes.Buffer
	err := r.IssueTemplate.Execute(&issueText, templateData)
	return issueText.String(), err
}

func (r *GitlabIssueReader) fetchNextIssuesPage() error {
	resp, err := gitlabRequest(r.privateToken, r.gitlabHost,
		fmt.Sprintf("/projects/%d/issues?state=opened&page=%d", r.projectId, r.nextPage))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("could not fetch issue page %d: %s", r.nextPage, resp.Status))
	}
	err = json.NewDecoder(resp.Body).Decode(&r.issueBuffer)
	if err != nil {
		return err
	}
	for index := len(r.issueBuffer) - 1; index >= 0; index-- {
		issue := &r.issueBuffer[index]
		comments, err := r.fetchComments(issue.Id)
		if err != nil {
			return err
		}
		issue.Description, err = r.formatIssueDescription(issue, comments)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GitlabIssueReader) Next() (*GitlabIssue, error) {
	if r.current == len(r.issueBuffer) {
		r.current = 0
		err := r.fetchNextIssuesPage()
		if err != nil {
			return nil, err
		}
		r.nextPage = r.nextPage + 1
		if len(r.issueBuffer) == 0 {
			return nil, nil
		}
	}
	issue := r.issueBuffer[r.current]
	r.current = r.current + 1
	return &issue, nil
}
