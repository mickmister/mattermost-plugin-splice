package splice

type ProjectsResponse struct {
	Error               string         `json:"error,omitempty"`
	Projects            []Project      `json:"projects"`
	CreationsCount      int            `json:"creations_count"`
	CollaborationsCount int            `json:"collaborations_count"`
	Collaborators       map[string]int `json:"collaborators"`
}

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Daw         string `json:"daw"`
	DawVersion  string `json:"daw_version"`
	ArtworkURL  string `json:"artwork_url"`
	Creator     struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"creator"`
	Collaborators []struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"collaborators"`
	Versions int `json:"versions"`
	LastSave struct {
		UUID     string `json:"uuid"`
		ID       int    `json:"id"`
		Username string `json:"username"`
		Version  int    `json:"version"`
		At       string `json:"at"`
	} `json:"last_save"`
	Bounce *struct {
		Path    string `json:"path"`
		UUID    string `json:"uuid"`
		Version int    `json:"version"`
	} `json:"bounce"`
	IsSplice      bool        `json:"is_splice"`
	IsSolo        bool        `json:"is_solo"`
	UserIds       []int       `json:"user_ids"`
	UserIds2      []int       `json:"user_ids2"`
	Usernames     []string    `json:"usernames"`
	SourceRelease interface{} `json:"source_release"`
	StemsCount    int         `json:"stems_count"`
	Licenses      interface{} `json:"licenses"`
	CreatorID     int         `json:"creator_id"`
	UserSets      interface{} `json:"user_sets"`
}
