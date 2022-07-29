package kruise

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/viper"
)

// Konfig struct used to combine file metadata with unmarshalled Kruise
// configuration
type Konfig struct {
	Files    []string
	Override string
	Manifest latest.KruiseConfig
}

func NewKonfig() *Konfig {
	cfg := new(Konfig)
	cwd, err := os.Getwd()
	Fatal(err)
	cfg.Files = []string{
		cwd + "/kruise.yaml",
		xdg.ConfigHome + "/kruise/kruise.yaml",
		xdg.Home + "/.kruise.yaml",
	}
	cfg.ApplyUserConfig()
	return cfg
}

// ApplyUserConfig reads in a configuration file that is passed to viper and
// unmarshalled
func (k *Konfig) ApplyUserConfig() {
	Logger.Debug("Setting user config")
	k.setUserConfig()
	Logger.Debug("Unmarshalling user config")
	k.unmarshalExactConfig()
	Logger.Trace("Config unmarshalled")
}

func (k Konfig) setUserConfig() {
	if k.Override != "" {
		Logger.Debugf("Overriding default user config: %s", k.Override)
		viper.SetConfigFile(k.Override)
		Logger.Debugf("Default user config overridden: %s", k.Override)
	} else {
		for _, file := range k.Files {
			Logger.Debugf("Searching for config in %s", file)
			if exists(file) {
				Logger.Debugf("Config found in %s", file)
				viper.SetConfigFile(file)
				break
			}
			Logger.Debugf("Config not found in %s", file)
		}
	}
	k.readConfig()
}

func (k Konfig) readConfig() {
	if err := viper.ReadInConfig(); err == nil {
		Logger.Infof("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		if k.Override != "" {
			Logger.Warnf("No user supplied config found in: %v", k.Override)
		} else {
			Logger.Warnf("No user supplied config found in the following paths: %v", k.Files)
		}
	}
}

func (k *Konfig) unmarshalExactConfig() {
	if err := viper.UnmarshalExact(&k.Manifest); err != nil {
		Fatalf(err, "Unable to decode config into struct")
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
