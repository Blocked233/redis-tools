package common

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppConfig 结构体用于匹配配置文件
type AppConfig struct {
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
}

var Config AppConfig

func init() {
	// 打开配置文件
	file, err := os.Open("common/config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	// 解析配置文件
	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	// 打印读取的配置
	fmt.Println("Redis Host:", Config.Redis.Host)
	fmt.Println("Redis Port:", Config.Redis.Port)
	fmt.Println("Redis Password:", Config.Redis.Password)
	fmt.Println("Redis DB:", Config.Redis.DB)
}
