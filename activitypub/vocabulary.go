package activitypub

type Actor struct {
	Context           []string  `json:"@context,omitempty"`
	Id                string    `json:"id,omitempty"`
	Type              string    `json:"type,omitempty"`
	Following         string    `json:"following,omitempty"`
	Followers         string    `json:"followers,omitempty"`
	Inbox             string    `json:"inbox,omitempty"`
	Outbox            string    `json:"outbox,omitempty"`
	PreferredUsername string    `json:"preferredUsername,omitempty"`
	Name              string    `json:"name,omitempty"`
	Summary           string    `json:"summary,omitempty"`
	Icon              Icon      `json:"icon"`
	Url               string    `json:"url,omitempty"`
	PublicKey         PublicKey `json:"publicKey,omitempty"`
	Endpoints         Endpoints `json:"endpoints"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox,omitempty"`
}

type Icon struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	Url       string `json:"url"`
}

type PublicKey struct {
	Id           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Outbox struct {
	Context      []string `json:"@context,omitempty"`
	Id           string   `json:"id,omitempty"`
	Type         string   `json:"type,omitempty"`
	TotalItems   int      `json:"totalItems"`
	OrderedItems []string `json:"orderedItems"`
}
