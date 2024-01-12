package gopolar

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	filePath string
	tunnels  []Tunnel
}

// read config from $HOME/.config/gopolar/gopolar.toml
// create config file if not exist
func NewConfig() *Config {
	ret := &Config{}

	// (create and) read config file
	cfgDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("fail to read $HOME dir: ", err)
	}
	cfgDir += "/.config/gopolar/"
	cfgPath := cfgDir + "gopolar.toml"
	ret.filePath = cfgPath
	viper.SetConfigName("gopolar")
	viper.AddConfigPath(cfgDir)
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // config file not found, create it and try again
			err = os.MkdirAll(cfgDir, 0700)
			if err != nil {
				log.Fatalln("fail to create config directory: " + cfgDir)
			}
			f, err := os.Create(cfgPath)
			f.Close()
			if err != nil {
				log.Fatalln("fail to create config file:", err)
			}
			return NewConfig()
		} else {
			log.Fatalln("fail to read config file:", err)
		}
	}

	// read tunnels from config file
	ts := viper.Get("tunnels").([]interface{})

	for _, t := range ts {
		ti := t.(map[string]interface{})
		res := Tunnel{}
		if err := mapstructure.Decode(ti, &res); err != nil {
			log.Fatal("fail to parse config file:", err)
		}
		ret.tunnels = append(ret.tunnels, res)
	}

	// TODO: remove this
	tunnels := []Tunnel{
		{
			ID:     1,
			Name:   "first tunnel",
			Enable: false,
			Source: "localhost:2345",
			Dest:   "233.168.10.1:5678",
		},
		{
			ID:     2,
			Name:   "second tunnel",
			Enable: false,
			Source: "localhost:3333",
			Dest:   "192.168.10.1:4567",
		},
		{
			ID:     3,
			Name:   "hahaha this is 3",
			Enable: true,
			Source: "localhost:2789",
			Dest:   "localhost:2333",
		},
	}
	viper.Set("tunnels", tunnels)
	viper.WriteConfig()

	for _, t := range ret.tunnels {
		fmt.Printf("read tunnels: %+v\n", t)
	}
	return ret
}
