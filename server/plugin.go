package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-plugin-splice/splice"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

var cmd = &model.Command{
	Trigger:      "splice",
	AutoComplete: true,
	AutocompleteData: &model.AutocompleteData{
		Trigger:  "splice",
		HelpText: "Interact with your Splice account",
		SubCommands: []*model.AutocompleteData{
			{
				Trigger:  "projects",
				HelpText: "Get the projects",
			},
			{
				Trigger:  "connect",
				HelpText: "Connect with token",
				Arguments: []*model.AutocompleteArg{
					{
						Type:     model.AutocompleteArgTypeText,
						Required: true,
						Data:     &model.AutocompleteTextArg{},
					},
				},
			},
		},
	},
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	parts := strings.Fields(args.Command)

	if len(parts) < 2 {
		return &model.CommandResponse{
			Text: "Please provide a subcommand",
		}, nil
	}

	switch parts[1] {
	case "connect":
		return p.handleCommandConnect(parts[2:])
	case "projects":
		return p.handleCommandProjects(parts[2:])
	}

	return &model.CommandResponse{
		Text: "Please provide a valid subcommand",
	}, nil
}

func (p *Plugin) OnActivate() error {
	return p.API.RegisterCommand(cmd)
}

func (p *Plugin) handleCommandConnect(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) == 0 {
		return &model.CommandResponse{
			Text: "Please provide a token",
		}, nil
	}

	token := strings.Join(args, " ")
	client := &splice.Client{
		Token: token,
	}

	_, err := client.GetProjects()
	if err != nil {
		return &model.CommandResponse{
			Text: "Error using token: " + err.Error(),
		}, nil
	}

	appErr := p.API.KVSet("token", []byte(token))
	if appErr != nil {
		return &model.CommandResponse{
			Text: "Failed to save token: " + appErr.Error(),
		}, nil
	}

	return &model.CommandResponse{
		Text: "Token is valid. Saved the token.",
	}, nil
}

func (p *Plugin) handleCommandProjects(args []string) (*model.CommandResponse, *model.AppError) {
	projectsResponse, err := p.getProjectsWithStoredToken()
	if err != nil {
		return &model.CommandResponse{
			Text: "Error getting splice projects: " + err.Error(),
		}, nil
	}

	projectData := []string{}
	for _, project := range projectsResponse.Projects {
		u := "https://splice.com/studio/" + project.UUID
		nameLink := fmt.Sprintf("[%s](%s)", project.Name, u)
		s := "* " + nameLink

		if project.Bounce != nil {
			bounceURL := "https://api.splice.com" + project.Bounce.Path
			bounceLink := fmt.Sprintf("[Version %v Bounce](%s)", project.Bounce.Version, bounceURL)
			s += " - " + bounceLink
		}

		projectData = append(projectData, s)
	}

	text := strings.Join(projectData, "\n")
	return &model.CommandResponse{
		Text: text,
	}, nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	projects, err := p.getProjectsWithStoredToken()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (p *Plugin) getProjectsWithStoredToken() (*splice.ProjectsResponse, error) {
	tokenBytes, appErr := p.API.KVGet("token")
	if appErr != nil {
		return nil, errors.Wrap(appErr, "Error fetching token")
	}

	if len(tokenBytes) == 0 {
		return nil, errors.New("No token stored. Run `/splice connect`")
	}

	client := &splice.Client{Token: string(tokenBytes)}
	projects, err := client.GetProjects()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting splice projects")
	}

	return projects, nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
