package global

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFile string = "antares.json"

type Conf struct {
	Root string `json:"root"`
	Port int    `json:"port"`
	Log  bool   `json:"log"`
}

func DefaultConfig() *Conf {
	return &Conf{
		Root: "./",
		Port: 9700,
		Log:  false,
	}
}

func LoadConfig() (*Conf, error) {
	// カレントディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(currentDir, ConfigFile)

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 設定ファイルが存在しない場合はデフォルト設定を返す
			return DefaultConfig(), nil
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var config Conf
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Conf) error {
	// カレントディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	configPath := filepath.Join(currentDir, ConfigFile)

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func GenerateDefaultConfig() error {
	return SaveConfig(DefaultConfig())
}
