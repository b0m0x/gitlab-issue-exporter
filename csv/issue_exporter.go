package csv

import (
	"encoding/csv"
	"io"
	"github.com/b0m0x/gitlab-issue-exporter/gitlab"
)

type CsvIssueWriter struct {
	writer *csv.Writer
}

func (w *CsvIssueWriter) Write(issue *gitlab.GitlabIssue) {
	w.writer.Write([]string{
		"Normal",
		issue.Title,
		issue.Description,
		issue.CreatedAt.Format("02.01.2006"),
		issue.Milestone.DueDate,
		issue.Assignee.Username,

	})
	w.writer.Flush()
}

func NewCsvIssueWriter(w io.Writer) *CsvIssueWriter {
	cw := &CsvIssueWriter{csv.NewWriter(w)}
	cw.writer.Write([]string{
		"priority",
		"subject",
		"description",
		"start_date",
		"due_date",
		"assigned_to",
	})
	return cw
}
