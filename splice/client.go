package splice

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Client struct {
	Token string
}

func (c *Client) GetProjects() (*ProjectsResponse, error) {
	token := c.Token
	if token == "" {
		return nil, errors.New("no token provided to Splice client")
	}

	u := "https://api.splice.com/studio/projects?collaborator=&f=&page=1&per_page=25&q=&set=&version=2"
	httpClient := &http.Client{}

	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	addCookies(r, token)

	res, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	projectsResponse := &ProjectsResponse{}
	err = json.NewDecoder(res.Body).Decode(&projectsResponse)
	if err != nil {
		return nil, err
	}

	res.Body.Close()

	if projectsResponse.Error != "" {
		return nil, errors.New("Splice returned error: " + projectsResponse.Error)
	}

	return projectsResponse, nil
}

func addCookies(r *http.Request, cookieString string) {
	cookieStrings := strings.Split(cookieString, "; ")
	for _, c := range cookieStrings {
		parsed := strings.Split(c, "=")
		if len(parsed) >= 2 {
			name, value := parsed[0], parsed[1]

			r.AddCookie(&http.Cookie{
				Name:  name,
				Value: value,
			})
		}
	}
}
