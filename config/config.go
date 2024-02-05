package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	Token     string
	Botprefix string
	config    *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	readFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	er := json.Unmarshal(readFile, &config)
	if er != nil {
		fmt.Println(er.Error())
		return er
	}

	Token = config.Token
	Botprefix = config.BotPrefix

	return nil

}
