package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	GitHub GitHubConfig `toml:"github"`
	Repos  ReposConfig  `toml:"repos"`
}

type GitHubConfig struct {
	Timeout  int             `toml:"timeout"`
	Projects []GitHubProject `toml:"projects"`
}

type GitHubProject struct {
	Name          string   `toml:"name"`
	Organization  string   `toml:"organization"`
	ProjectNumber int      `toml:"project_number"`
	StatusFilters []string `toml:"status_filters"`
}

type ReposConfig struct {
	Directory    string `toml:"directory"`
	AutoDiscover bool   `toml:"auto_discover"`
}

func LoadConfig(atpDir string) (*Config, error) {
	configPath := filepath.Join(atpDir, "config.toml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := getDefaultConfig()
		if err := SaveConfig(atpDir, config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(atpDir string, config *Config) error {
	configPath := filepath.Join(atpDir, "config.toml")
	
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getDefaultConfig() *Config {
	return &Config{
		GitHub: GitHubConfig{
			Timeout:  30,
			Projects: []GitHubProject{},
		},
		Repos: ReposConfig{
			Directory:    "",
			AutoDiscover: true,
		},
	}
}

func (c *Config) GetGitHubProject(name string) (*GitHubProject, error) {
	for _, project := range c.GitHub.Projects {
		if project.Name == name {
			return &project, nil
		}
	}
	return nil, fmt.Errorf("GitHub project '%s' not found in config", name)
}

func (c *Config) GetAllGitHubProjects() []GitHubProject {
	return c.GitHub.Projects
}