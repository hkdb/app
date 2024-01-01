package db

type Packages struct {
	Packages    string `json:"packages"`
}

type Repos struct {
	Repo	string `json:"repo"`
	Gpg		string `json:"gpg"`
}

type Classics struct {
	Classic bool `json:"classic"`
} 

type Channels struct {
	Channel string `json:"channel"`
}
