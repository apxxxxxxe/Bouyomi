package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NoWordPhrase string `json:"NoWordPhrase"`
	JapaneseOnly bool   `json:"JapaneseOnly"`
}

func initConfig(path string) (*Config, error) {
	config := &Config{
		NoWordPhrase: "んん",
		JapaneseOnly: true,
	}

	bJson, err := json.MarshalIndent(config, "", "	")
	err = os.WriteFile(path, bJson, 0644)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadConfig() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(filepath.Dir(exePath), "config.json")

	fp, err := os.Open(path)
	if err != nil {
		return initConfig(path)
	}
	defer fp.Close()

	var config Config
	if err := json.NewDecoder(fp).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
