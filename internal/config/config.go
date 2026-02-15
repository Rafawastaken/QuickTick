package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SupabaseURL string `json:"supabase_url"`
	SupabaseKey string `json:"supabase_key"`
	Token       string `json:"token"`
}

func ConfigPath() (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	var cfg Config
	if os.IsNotExist(err) {
		// File doesn't exist, start with empty config
	} else if err != nil {
		return nil, err
	} else {
		// File exists, unmarshal it
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	// Fallback to Env Vars
	if cfg.SupabaseURL == "" {
		cfg.SupabaseURL = os.Getenv("SUPABASE_URL")
	}
	if cfg.SupabaseKey == "" {
		cfg.SupabaseKey = os.Getenv("SUPABASE_KEY")
		if cfg.SupabaseKey == "" {
			cfg.SupabaseKey = os.Getenv("SUPABASE_ANON_KEY")
		}
		if cfg.SupabaseKey == "" {
			cfg.SupabaseKey = os.Getenv("PUBLISHABLE_KEY") // User mentioned this
		}
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("SUPABASE_TOKEN")
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
