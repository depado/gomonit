package conf

// Repo is a configuration struct to access a repo
type Repo struct {
	Type  string `yaml:"type"`
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
	Path  string `yaml:"path"`
}

// CI is a configuration struct to access a CI
type CI struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
}

// Service is a configuration struct describing a service
type Service struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	URL  string `yaml:"url"`
	Icon string `yaml:"icon"`
	Own  bool   `yaml:"own"`

	CI   *CI   `yaml:"ci"`
	Repo *Repo `yaml:"repo"`
}
