package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mgurov/hiptee/pkg"
	"github.com/mgurov/hiptee/pkg/hip"
	"github.com/mgurov/hiptee/pkg/std"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

const mainName = "hiptee"

var commit string
var version string

func main() {

	var hipchatToken string
	var hipchatRoom string
	var hipchatPollPrefix string

	var config string
	showVersionAndExit := flag.Bool("version", false, "show version and exit")
	flag.StringVar(&config, "config", "", "config file. Defaults to HIPTEE_CONFIG environmental variable. Json of contents {\"token\":<token>,\"room\":<room>} expected.")
	flag.StringVar(&hipchatToken, "token", "", "hipchat token to send notice with. Defaults to HIPCHAT_TOKEN environmental variable.")
	flag.StringVar(&hipchatRoom, "room", "", "hipchat room to send notice to. Defaults to HIPCHAT_ROOM environmental variable.")
	flag.StringVar(&hipchatPollPrefix, "poll", "", "if not empty, the hipchat room will be polled for the messages starting with this prefix and the remainder of the line will be sent to the stdin of the command in the exec mode")
	flag.Parse()

	if *showVersionAndExit {
		fmt.Println(mainName, "version", version, "commit", commit)
		return
	}

	viper.BindFlagValue("hipchat.token", stdFlagAdaptorValue{name: "token", value: &hipchatToken})
	viper.BindFlagValue("hipchat.room", stdFlagAdaptorValue{name: "room", value: &hipchatRoom})
	viper.BindFlagValue("hipchat.poll_prefix", stdFlagAdaptorValue{name: "room", value: &hipchatPollPrefix})

	viper.BindEnv("hipchat.token", "HIPCHAT_TOKEN")
	viper.BindEnv("hipchat.room", "HIPCHAT_ROOM")
	viper.BindEnv("hipchat.poll_prefix", "HIPCHAT_POLL_PREFIX")
	readViperConfig(config)

	if hipchatToken = viper.GetString("hipchat.token"); "" == hipchatToken {
		exitUsageIfEmpty(hipchatToken, "Hipchat token missing")
	}

	if hipchatRoom = viper.GetString("hipchat.room"); "" == hipchatRoom {
		exitUsageIfEmpty(hipchatRoom, "Hipchat room missing")
	}

	hipchatPollPrefix = viper.GetString("hipchat.poll_prefix")

	hc, err := hip.NewHipchatRoomPrinter(hipchatToken, hipchatRoom)
	if err != nil {
		log.Fatal("Error establishing hipchat printing ->"+hipchatToken+"<- ", err)
	}

	hipchatAndStdout := pkg.Compose(hc, &std.StdOutPrinter{})

	if len(os.Args) <= 1 {
		//classic tee from stdin

		if hipchatPollPrefix != "" {
			exitUsageIfEmpty("", "don't know what to do with the hipchat polled messages in the tee mode")
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			hipchatAndStdout.Out(scanner.Text())
		}
		hipchatAndStdout.Done(nil)
	} else {
		var inCommandsReader io.Reader

		if hipchatPollPrefix != "" {
			inCommandsReader, err = hip.NewHipchatRoomReader(hipchatToken, hipchatRoom, hipchatPollPrefix)
			if nil != err {
				log.Fatal("Error establishing hipchat polling ", err)
			}
		}

		command := flag.Arg(0)
		params := flag.Args()[1:]
		if err := pkg.Execute(command, params, hipchatAndStdout, inCommandsReader); nil != err {
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
