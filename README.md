# Gitlab Issue Exporter

This is a command line utility to export all open issues from a Gitlab installation to a .csv file.
The resulting csv contains the title, description, start date, due date (if set on the milestone of the ticket) and current assignee.
The description contains a complete history of all comments for the issue and is formatted in Markdown.

## Installation
You need to have go installed and set up.
Install with

```
   go get github.com/b0m0x/gitlab-issue-exporter
```   

   
## Usage
Export:

```
   gitlab-issue-exporter --host git.your-organisation.com --token YOUR_GITLAB_PRIVATE_TOKEN --project your-org/your-project
```

 The export is stored in `export.csv`. The format is compatible with Redmine if Redmine is configured to use markdown.