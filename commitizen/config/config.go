package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config loads config file in git repo or home directory
type Config interface {
	Setup(configPath string) error
	Get(key string) interface{}
	GetString(key string) string
}

type basicConfig struct {
	v                *viper.Viper
	customConfigPath string
}

// New returns a empty config struct
func New() Config {
	return &basicConfig{
		v:                viper.New(),
		customConfigPath: "",
	}
}

func (cfg *basicConfig) Get(key string) interface{} {
	return cfg.v.Get(key)
}
func (cfg *basicConfig) GetString(key string) string {
	return cfg.v.GetString(key)
}
func (cfg *basicConfig) Setup(configPath string) error {
	cfg.customConfigPath = configPath
	if err := cfg.readConfig(); err != nil {
		return err
	}
	return nil
}

func getSystemRootPath(dir string) string {
	systemRootPath := filepath.VolumeName(dir) + string(filepath.Separator)
	return systemRootPath
}

func findRepoRootPath() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	// we assume no .git folder exists in "/" for linux and mac, "C:\\" for windows.
	for root := getSystemRootPath(dir); dir != root; dir = filepath.Dir(dir) {
		dotGitPath := filepath.Join(dir, ".git")
		if stat, err := os.Stat(dotGitPath); err == nil && stat.IsDir() {
			// path is a directory
			logrus.Debug("we find repo base path: " + dir)
			return dir, true
		}
	}

	return "", false
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func absPathify(inPath string) string {

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

// return filepath as below order
// 1. repoDir/.gocz.toml (first ancestor folder "repoDir" which contains .git)
// 2. $HOME/.gocz/gocz.toml
// 3. $HOME/.gocz.toml
// 4. /etc/gocz.toml
func defaultConfigPaths() []string {
	configPaths := []string{}
	if p, ok := findRepoRootPath(); ok {
		p = filepath.Join(p, ".gocz.toml")
		configPaths = append(configPaths, p)
	}
	configPaths = append(configPaths, absPathify(filepath.Join("$HOME", ".gocz", "gocz.toml")))
	configPaths = append(configPaths, absPathify(filepath.Join("$HOME", ".gocz.toml")))

	return configPaths
}

func existFile(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err := f.Stat(); err != nil {
		return false
	}
	return true
}

func (cfg *basicConfig) readConfig() (err error) {
	if cfg.customConfigPath != "" {
		// already set Custom Config Path
		cfg.v.SetConfigFile(cfg.customConfigPath)
		err = cfg.v.ReadInConfig()
	} else {
		err = errors.New("no default config file exist")
		// try to read defaultConfigPath one by one until first success
		for _, p := range defaultConfigPaths() {
			if !existFile(p) {
				continue
			}
			cfg.v.SetConfigFile(p)
			err = cfg.v.ReadInConfig()
			if err == nil {
				break
			}
		}
	}
	if err == nil {
		logrus.Debug("Config File Used: " + cfg.v.ConfigFileUsed())
	}
	// Return err if all failed else nil
	return err
}
