package activitypub

type Webfinger struct {
	Subject string   `json:"subject,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
	Links   []Link   `json:"links,omitempty"`
}

type Link struct {
	Rel      string `json:"rel,omitempty"`
	Type     string `json:"type,omitempty"`
	Href     string `json:"href,omitempty"`
	Template string `json:"template,omitempty"`
}

type Nodeinfo struct {
	Version           string   `json:"version,omitempty"`
	Software          Software `json:"software,omitempty"`
	Protocols         []string `json:"protocols,omitempty"`
	Services          Services `json:"services,omitempty"`
	Usage             Usage    `json:"usage,omitempty"`
	OpenRegistrations bool     `json:"openRegistrations,omitempty"`
	Metadata          struct{} `json:"metadata,omitempty"`
}

type Software struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type Services struct {
	Outbound []string `json:"outbound,omitempty"`
	Inbound  []string `json:"inbound,omitempty"`
}

type Users struct {
	Total          int `json:"total,omitempty"`
	ActiveMonth    int `json:"activeMonth,omitempty"`
	ActiveHalfyear int `json:"activeHalfyear,omitempty"`
}
type Usage struct {
	Users      Users `json:"users,omitempty"`
	LocalPosts int   `json:"inbound,omitempty"`
}
