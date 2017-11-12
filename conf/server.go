package conf

// Server is the structure that holds the server configuration
type Server struct {
	Host  string `yaml:"host" default:"127.0.0.1"`
	Port  int    `yaml:"port" default:"8080"`
	Debug bool   `yaml:"debug"`
}
