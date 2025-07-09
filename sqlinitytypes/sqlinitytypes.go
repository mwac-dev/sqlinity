package sqlinitytypes

type Migration struct {
	ID      string
	Name    string
	UpSQL   string
	DownSQL string
}

type Config struct {
	SqlFolder    string `json:"sqlFolder"`
	OutputFolder string `json:"outputFolder"`
	Namespace    string `json:"namespace"`
}
