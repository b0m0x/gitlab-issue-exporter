package gitlab

import (
	"net/http"
	"fmt"
	"encoding/json"
	"errors"
)

func gitlabRequest(privateToken, host, path string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/api/v3%s", host, path), nil)
	req.Header.Add("PRIVATE-TOKEN", privateToken)
	return client.Do(req)
}

func getProjectId(privateToken, host, projectName string) (int, error) {
	resp, err := gitlabRequest(privateToken, host, fmt.Sprintf("/projects/%s", projectName))
	if err != nil {
		return -1, err
	}
	if resp.StatusCode != 200 {
		return -1, errors.New("project fetch error: " + resp.Status)
	}
	var projectInfo struct {
		Id int
	}
	err = json.NewDecoder(resp.Body).Decode(&projectInfo)
	if err != nil {
		return -1, err
	}
	return projectInfo.Id, nil
}
