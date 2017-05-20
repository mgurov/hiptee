package main

import (
	"bufio"
	"flag"
	"github.com/mgurov/hiptee/pkg"
	"github.com/mgurov/hiptee/pkg/hip"
	"github.com/mgurov/hiptee/pkg/std"
	"github.com/spf13/viper"
	"log"
	"os"
)

const mainName = "hiptee"

func main() {

	var hipchatConfig struct {
		token string
		room  string
		third string
	}

	var config string
	flag.StringVar(&config, "config", "", "config file. Defaults to HIPTEE_CONFIG environmental variable. Json of contents {\"token\":<token>,\"room\":<room>} expected.")
	flag.StringVar(&hipchatConfig.token, "token", "", "hipchat token to send notice with. Defaults to HIPCHAT_TOKEN environmental variable.")
	flag.StringVar(&hipchatConfig.room, "room", "", "hipchat room to send notice to. Defaults to HIPCHAT_ROOM environmental variable.")
	flag.Parse()

	viper.BindFlagValue("hipchat.token", stdFlagAdaptorValue{name: "token", value: &hipchatConfig.token})
	viper.BindFlagValue("hipchat.room", stdFlagAdaptorValue{name: "room", value: &hipchatConfig.room})

	viper.BindEnv("hipchat.token", "HIPTEE_TOKEN")
	viper.BindEnv("hipchat.room", "HIPTEE_ROOM")
	readViperConfig(config)

	if hipchatConfig.token = viper.GetString("hipchat.token"); "" == hipchatConfig.token {
		exitUsageIfEmpty(hipchatConfig.token, "Hipchat token missing")
	}

	if hipchatConfig.room = viper.GetString("hipchat.room"); "" == hipchatConfig.room {
		exitUsageIfEmpty(hipchatConfig.token, "Hipchat room missing")
	}

	hc, err := hip.NewHipchatRoomPrinter(hipchatConfig.token, hipchatConfig.room)
	if err != nil {
		log.Fatal(err)
	}

	hipchatAndStdout := pkg.Compose(hc, &std.StdOutPrinter{})

	if len(os.Args) <= 1 {
		//classic tee from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			hipchatAndStdout.Out(scanner.Text())
		}
		hipchatAndStdout.Done(nil)
	} else {
		command := os.Args[1]
		params := os.Args[2:]
		if err := pkg.Execute(command, params, hipchatAndStdout); nil != err {
			os.Exit(1)
		}
	}
}

func readViperConfig(explicitConfigFile string) {
	if "" != explicitConfigFile {
		viper.SetConfigFile(explicitConfigFile)
	} else {
		viper.SetConfigName(mainName)
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err == nil {
		return
	}

	_, configNotFound := err.(viper.ConfigFileNotFoundError)

	if configNotFound && "" == explicitConfigFile {
		//config may be missing unless explicitly pointed to
		return
	}

	println("Could not read config file:", err.Error())
	os.Exit(1)
}

type stdFlagAdaptorValue struct {
	name  string
	value *string
}

func (f stdFlagAdaptorValue) HasChanged() bool    { return "" != *f.value }
func (f stdFlagAdaptorValue) Name() string        { return f.name }
func (f stdFlagAdaptorValue) ValueString() string { return *f.value }
func (f stdFlagAdaptorValue) ValueType() string   { return "string" }

func exitUsageIfEmpty(value, message string) {
	if value == "" {
		println(message)
		flag.Usage()
		os.Exit(1)
	}
}
