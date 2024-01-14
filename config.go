package gopolar

import (
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
	ts, ok := viper.Get("tunnels").([]interface{})
	if ok {
		for _, t := range ts {
			ti := t.(map[string]interface{})
			res := Tunnel{}
			if err := mapstructure.Decode(ti, &res); err != nil {
				log.Fatal("fail to parse config file:", err)
			}
			ret.tunnels = append(ret.tunnels, res)
		}
	} // else ts is nil

	return ret
}
