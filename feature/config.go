package feature

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Toggles struct {
		Bhop     bool`json:"bhop"`
		Name     bool`json:"name"`
		Skeleton bool`json:"skeleton"`
	} `json:"toggles"`
	ColorModes struct {
		Name     int`json:"name"`
		Skeleton int`json:"skeleton"`
	} `json:"colorModes"`
	SeeTeammates struct {
		Name    bool `json:"name"`
		Skeleton bool`json:"skeleton"`
	} `json:"seeTeammates"`
}

const (
	ColorModeHealth = iota
	ColorModeTeam = iota
)

func InitConfig() *Config{
	confBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			createDefaultConfig()
			InitConfig()
		}
		log.Println("InitConfig", err)
	}

	var conf Config
	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		log.Println("InitConfig", err)
	}
	return &conf
}

func WatchConfig(mainConfig *Config) {
	confBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			createDefaultConfig()
			go WatchConfig(mainConfig)
			return
		}
		log.Println("watchConfig", err)
	}

	initHash := hash(string(confBytes))
	for {
		currBytes, err := ioutil.ReadFile("config.json")
		if err != nil {
			if os.IsNotExist(err) {
				createDefaultConfig()
				go WatchConfig(mainConfig)
				return
			}
			log.Println("watchConfig", err)
		}
		hash := hash(string(currBytes))
		if hash != initHash {
			initHash = hash
			confBytes, err = ioutil.ReadFile("config.json")
			err = json.Unmarshal(confBytes,&mainConfig)
			if err != nil {
				log.Println("WatchConfig", err)
			}
			fmt.Println("Config changed")
		}

		time.Sleep(time.Millisecond * 15)
	}
}

func createDefaultConfig() {
	fmt.Println("Creating default config")
	c := defaultConfig()
	b, err := json.MarshalIndent(&c,"","	")
	if err != nil {
		log.Println("createDefaultConfig",err)
	}
	err = ioutil.WriteFile("config.json",b,0644)
	if err != nil {
		log.Println("createDefaultConfig",err)
	}
}

func defaultConfig() Config{
	return Config{
		Toggles: struct {
			Bhop     bool `json:"bhop"`
			Name     bool `json:"name"`
			Skeleton bool `json:"skeleton"`
		}{true,true,true},
		ColorModes: struct {
			Name     int `json:"name"`
			Skeleton int `json:"skeleton"`
		}{ColorModeTeam,ColorModeHealth},
		SeeTeammates: struct {
			Name     bool `json:"name"`
			Skeleton bool `json:"skeleton"`
		}{true,true},
	}

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
