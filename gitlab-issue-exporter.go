package main

import (
	"fmt"
	"os"
	"github.com/b0m0x/gitlab-issue-exporter/csv"
	"github.com/b0m0x/gitlab-issue-exporter/gitlab"
	"flag"
)



func main() {
	var project = flag.String("project", "", "the complete project name, including the namespace. e.g: example/example")
	var privateToken = flag.String("token", "", "a gitlab private token with read access to the issues of the requested project")
	var gitlabHost	 = flag.String("host", "", "the host name or ip of the gitlab installation, e.g. git.your-org.com")
	flag.Parse()

	issueReader, err := gitlab.NewGitlabIssueReader(*privateToken, *gitlabHost, *project)
	if err != nil {
		fmt.Printf("GitLab Issue Reader error: %s\n", err.Error())
		return
	}
	f, err := os.Create("export.csv")
	defer f.Close()
	if err != nil {
		fmt.Printf("file error: %s\n", err.Error())
		return
	}
	csvWriter := csv.NewCsvIssueWriter(f)
	for issue, err := issueReader.Next(); issue != nil; issue, err = issueReader.Next() {
		if err != nil {
			fmt.Printf("GitLab Issue Reader error: %s\n", err.Error())
			return
		}
		fmt.Printf("Exporting issue %d: %s\n", issue.Id, issue.Title)
		csvWriter.Write(issue)
	}
	fmt.Println("done")
}
