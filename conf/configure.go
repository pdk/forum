package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Configuration contains settings for the current run
type Configuration struct {
	Database      string
	ListenAddress string
	AssetsDir     string
}

// ReadConfiguration reads the named file as JSON and returns the Configuration.
func ReadConfiguration(fileName string) (Configuration, error) {

	newConf := Configuration{}

	log.Printf("reading configuration from %s", fileName)

	confData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return newConf, fmt.Errorf("cannot read configuration: %w", err)
	}

	err = json.Unmarshal(confData, &newConf)
	if err != nil {
		return newConf, fmt.Errorf("cannot parse configuration %s: %w", fileName, err)
	}

	return newConf, nil
}
