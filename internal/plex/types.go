package plex

type PlexAccount struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumb"`
}

type PlexDevice struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Product string `json:"product"`
}

type PlexPayload struct {
	Event   string      `json:"event"`
	Account PlexAccount `json:"Account"`
	Player  struct {
		Uuid  string `json:"uuid"`
		Title string `json:"title"`
	} `json:"Player"`
}
