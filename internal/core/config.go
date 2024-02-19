package core

import (
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	DoLogs    bool
	ReadSaved bool
}

var DefaultConfig Config = Config{
	DoLogs:    false,
	ReadSaved: true,
}

// read tunnels from $HOME/.gopolar/tunnels.toml
// create it if not exist
func readTunnels() []Tunnel {
	cfgDir := homeDir + "/.gopolar/"
	cfgPath := cfgDir + "tunnels.toml"
	viper.SetConfigName("tunnels")
	viper.AddConfigPath(cfgDir)
	if err := viper.ReadInConfig(); err != nil {
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
			return nil
		} else {
			log.Fatalln("fail to read config file:", err)
		}
	}

	// read tunnels from config file
	ts, ok := viper.Get("tunnels").([]interface{})
	ret := []Tunnel{}
	if ok {
		for _, t := range ts {
			ti := t.(map[string]interface{})
			res := Tunnel{}
			if err := mapstructure.Decode(ti, &res); err != nil {
				log.Fatal("fail to parse config file:", err)
			}
			ret = append(ret, res)
		}
	} // else ts is nil

	return ret

}
