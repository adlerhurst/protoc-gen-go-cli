package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Options struct {
		Cmd     string
		Logging string
		Root    string
		Adler   string
	}
}

func main() {
	//Configure the type of the configuration as JSON
	viper.SetConfigType("json")
	//Set the environment prefix as CONFIG
	viper.SetEnvPrefix("CONFIG")

	//Substitute the _ to .
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	//Get the string that is set in the CONFIG_OPTIONS_VALUES environment variable
	var jsonExample = []byte(viper.GetString("options.values"))
	if err := viper.ReadConfig(bytes.NewBuffer(jsonExample)); err != nil {
		log.Printf("failed to read: %v\n", err)
	}

	c := new(Config)
	if err := viper.Unmarshal(c); err != nil {
		log.Printf("failed to unmarshal: %v\n", err)
	}
	fmt.Printf("%#v\n", c.Options)

	//Convert the sub-json string of the options field in map[string]string
	fmt.Println(viper.GetStringMapString("options"))
}
