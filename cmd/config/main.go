package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marioharper/commuter/directions"
)

type Config struct {
	Locations []directions.Location
}

func GetConfig(configFile string) Config {
	var config Config
	f := getConfigFile(configFile)

	jsonParser := json.NewDecoder(f)

	if fi, err := f.Stat(); err != nil {
		fmt.Printf("getting file stat: %s \n", err.Error())
	} else {
		if fi.Size() == 0 {
			fmt.Println("config file empty")
			return config
		}
	}

	if err := jsonParser.Decode(&config); err != nil {
		fmt.Printf("parsing config file: %s \n", err.Error())
		os.Exit(-1)
	}

	return config
}

func SaveConfig(configFile string, config Config) {
	fmt.Println("saving config")
	f := getConfigFile(configFile)

	// get file size
	if data, err := f.Stat(); err != nil {
		fmt.Printf("error gettings specs")
	} else {
		fmt.Printf("file size: %d \n", data.Size())
	}

	// convert config to json
	configJSON, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("marshalling config file: %s \n", err.Error())
		os.Exit(-1)
	}
	fmt.Println(string(configJSON))

	// write json to config file
	n4, err := f.WriteString(string(configJSON))
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("Wrote %d bytes \n", n4)

	f.Sync()

	// Get file size
	if data, err := f.Stat(); err != nil {
		fmt.Printf("error gettings specs")
	} else {
		fmt.Printf("file size: %d \n", data.Size())
	}

}

/////////////////////////////////

func getConfigFile(configFile string) *os.File {
	var f *os.File
	var err error

	// if no config file, create it
	if f, err = os.OpenFile(configFile, os.O_RDWR, 0666); os.IsNotExist(err) {
		f, err = os.Create(configFile)
		if err != nil {
			fmt.Printf("creating config file: %s \n", err.Error())
			os.Exit(-1)
		}
	}

	return f
}
