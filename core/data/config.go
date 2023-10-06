package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	noWordPhrase = "んん"
	japaneseOnly = true
	nullPath     = ""
)

type Config struct {
	NoWordPhrase    *string `json:"NoWordPhrase,omitempty"`
	JapaneseOnly    *bool   `json:"JapaneseOnly,omitempty"`
	BouyomiChanPath *string `json:"BouyomiChanPath,omitempty"`
	CoeiroIncPath   *string `json:"CoeiroIncPath,omitempty"`
	VoiceVoxPath    *string `json:"VoiceVoxPath,omitempty"`
}

func initConfig(path string) (*Config, error) {
	config := &Config{
		NoWordPhrase:    &noWordPhrase,
		JapaneseOnly:    &japaneseOnly,
		BouyomiChanPath: &nullPath,
		CoeiroIncPath:   &nullPath,
		VoiceVoxPath:    &nullPath,
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

	if config.NoWordPhrase == nil {
		config.NoWordPhrase = &noWordPhrase
		isReset = true
	}
	if config.JapaneseOnly == nil {
		config.JapaneseOnly = &japaneseOnly
		isReset = true
	}
	if config.BouyomiChanPath == nil {
		config.BouyomiChanPath = &nullPath
		isReset = true
	}
	if config.CoeiroIncPath == nil {
		config.CoeiroIncPath = &nullPath
		isReset = true
	}
	if config.VoiceVoxPath == nil {
		config.VoiceVoxPath = &nullPath
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
