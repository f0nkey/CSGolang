package feature

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chewxy/math32"
	"hash/fnv"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

const (
	ColorModeHealth = iota
	ColorModeTeam   = iota
)

type Config struct {
	Toggles struct {
		Bhop     bool `json:"bhop"`
		Name     bool `json:"name"`
		Skeleton bool `json:"skeleton"`
	} `json:"toggles"`
	ColorModes struct {
		Name     int `json:"name"`
		Skeleton int `json:"skeleton"`
	} `json:"colorModes"`
	SeeTeammates struct {
		Name     bool `json:"name"`
		Skeleton bool `json:"skeleton"`
	} `json:"seeTeammates"`
}

func defaultConfig() Config {
	return Config{
		Toggles: struct {
			Bhop     bool `json:"bhop"`
			Name     bool `json:"name"`
			Skeleton bool `json:"skeleton"`
		}{true, true, true},
		ColorModes: struct {
			Name     int `json:"name"`
			Skeleton int `json:"skeleton"`
		}{ColorModeTeam, ColorModeHealth},
		SeeTeammates: struct {
			Name     bool `json:"name"`
			Skeleton bool `json:"skeleton"`
		}{true, true},
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func RunConfigEndpoint(masterConfig *Config) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		enableCors(&rw)
		if r.Method == "POST" {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			oldConfig := *masterConfig
			err := json.Unmarshal(bodyBytes, &masterConfig)
			if err != nil {
				rw.Write([]byte("Error json parsing:" + err.Error()))
				log.Println("runConfigEndpoint:", err)
			}
			from, to, err := changedConfigField(&oldConfig, masterConfig)
			if err != nil {
				rw.Write([]byte("Changed nothing"))
				return
			}
			err = writeConfig(masterConfig)
			if err != nil {
				log.Println("RunConfigEndpoint", err)
			}
			changedMsg := fmt.Sprintf("Changed %v = %v to %v", from.Field, from.Value, to.Value)
			rw.Write([]byte(changedMsg))
		} else if r.Method == "GET"{
			byteConfig,_ := json.Marshal(masterConfig)
			rw.Write(byteConfig)
		}
	})
	if err := http.ListenAndServe(":9991", nil); err != nil {
		panic(err)
	}
}

func ServeWebGUI() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("ui/public")))
	http.ListenAndServe(":8085", mux)
}

func writeConfig(c *Config) error {
	b, err := json.MarshalIndent(&c, "", "	")
	if err != nil {
		return errors.New("writeConfig:" + err.Error())
	}
	err = ioutil.WriteFile("config.json", b, 0644)
	if err != nil {
		return errors.New("writeConfig:" + err.Error())
	}
	return nil
}

func changedConfigField(old, new *Config) (FieldValuePair, FieldValuePair, error) {
	oldPairs := fieldValPairs(*old)
	newPairs := fieldValPairs(*new)
	for i := 0; i < len(oldPairs); i++ {
		if oldPairs[i].Value != newPairs[i].Value {
			return oldPairs[i], newPairs[i], nil
		}
	}
	return FieldValuePair{}, FieldValuePair{}, errors.New("Found no changes")
}

func fieldValPairs(x Config) []FieldValuePair {
	s := reflect.ValueOf(&x).Elem()
	res := make([]FieldValuePair, s.NumField())
	typeOfX := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		field := fmt.Sprint(typeOfX.Field(i).Name, f.Type())
		res = append(res, FieldValuePair{field, f.Interface()})
	}
	return res
}

type FieldValuePair struct {
	Field string
	Value interface{}
}

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

func WatchConfig(masterConfig *Config) {
	confBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			createDefaultConfig()
			go WatchConfig(masterConfig)
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
				go WatchConfig(masterConfig)
				return
			}
			log.Println("watchConfig", err)
		}
		hash := hash(string(currBytes))
		if hash != initHash {
			initHash = hash
			confBytes, err = ioutil.ReadFile("config.json")
			err = json.Unmarshal(confBytes, &masterConfig)
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
	b, err := json.MarshalIndent(&c, "", "	")
	if err != nil {
		log.Println("createDefaultConfig", err)
	}
	err = ioutil.WriteFile("config.json", b, 0644)
	if err != nil {
		log.Println("createDefaultConfig", err)
	}
}

func getColor(colorMode int, hp int32, team int32) color.RGBA {
	if colorMode == ColorModeHealth {
		r := uint8(math32.Min((510*(100-float32(hp)))/100, 255))
		g := uint8(math32.Min((510*float32(hp))/100, 255))
		return color.RGBA{r, g, 0, 255}
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
