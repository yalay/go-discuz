package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ClassIdMap       map[string]int
	EnableEcms       bool
	EnableFidMapping bool
	EnableGenKeyword bool
	EnableFormatImg  bool

	EcmsSql   *SqlConfig
	DiscuzSql *SqlConfig

	DataFile string
	Dict     string
}

type SqlConfig struct {
	UserName string
	Password string
	Database string
	DbPrefix string
	Author   string
	AuthorId int
}

func LoadConfig(filePath string) *Config {
	jsonCfg, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read file err:%v", err)
		return nil
	}

	var cfg = &Config{}
	err = json.Unmarshal(jsonCfg, cfg)
	if err != nil {
		fmt.Printf("json Unmarshal err:%v", err)
		return nil
	}
	return cfg
}

func (c *Config) GetMappingFid(oriClassId int) int {
	if curClassId, ok := c.ClassIdMap[fmt.Sprintf("%d", oriClassId)]; ok {
		return curClassId
	}
	return 0
}
