package godis

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type GodisConfig struct {
	OutputToTerminal bool
	LogDir           string
	*viper.Viper
}

var GodisConf *GodisConfig

func ParseConf() {
	v := viper.New()

	confFile := "./godis_conf.toml"

	v.SetConfigFile(confFile)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config %s err:%v", confFile, err)
		os.Exit(-1)
	}

	v.SetDefault("OutputToTerminal", true)
	v.SetDefault("LogDir", "./log/")

	GodisConf = &GodisConfig{}
	if err := v.Unmarshal(GodisConf); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal config :%v", err)
		os.Exit(-1)
	}

	GodisConf.Viper = v
}
