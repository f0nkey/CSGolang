package feature

import (
	"encoding/json"
	"fmt"
	"github.com/chewxy/math32"
	"hash/fnv"
	"image/color"
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

func InitConfig() *Config {
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

func getColor(colorMode int, hp int32, team int32) color.RGBA{
	if colorMode == ColorModeHealth {
		r := uint8(math32.Min((510*(100-float32(hp)))/100, 255))
		g := uint8(math32.Min((510*float32(hp))/100, 255))
		return color.RGBA{r,g,0,255}
	}
	if colorMode == ColorModeTeam {
		if team == 2 {
			return color.RGBA{255, 131, 8, 255} //ct
		}
		return color.RGBA{7, 21, 210, 255} //t
	}
	return color.RGBA{255, 255, 255, 255}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
