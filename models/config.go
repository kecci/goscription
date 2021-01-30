package models

// AppConfig is application configuration
type (
	Config struct {
		Title          string   `mapstructure:"title"`
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Database       Database `mapstructure:"database"`
		Godaddy        Godaddy  `mapstructure:"godaddy"`
	}

	Server struct {
		Address string `mapstructure:"address"`
	}

	Database struct {
		Driver string `mapstructure:"driver"`
		Host   string `mapstructure:"host"`
		Port   string `mapstructure:"port"`
		User   string `mapstructure:"user"`
		Pass   string `mapstructure:"pass"`
		Name   string `mapstructure:"name"`
	}

	Godaddy struct {
		Host          string `mapstructure:"host"`
		Authorization string `mapstructure:"authorization"`
	}
)
