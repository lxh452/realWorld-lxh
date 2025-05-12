package config

type Logs struct {
	Level  string `json:"level"`
	Dir    string `json:"dir"`
	Layout string `json:"layout"`
}
