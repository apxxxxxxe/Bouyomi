package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	japaneseOnly = true
	nilPath      = []string{""}
)

type Config struct {
	JapaneseOnly *bool    `json:"JapaneseOnly,omitempty"`
	EnginesPath  []string `json:"BouyomiChanPath,omitempty"`
}

func initConfig(path string) (*Config, error) {
	config := &Config{
		JapaneseOnly: &japaneseOnly,
		EnginesPath:  nil,
	}

	bJson, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return nil, err
	}

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

	isReset := false

	if config.JapaneseOnly == nil {
		config.JapaneseOnly = &japaneseOnly
		isReset = true
	}
	if config.EnginesPath == nil {
		config.EnginesPath = nilPath
		isReset = true
	}

	if isReset {
		bJson, err := json.MarshalIndent(config, "", "	")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(path, bJson, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
