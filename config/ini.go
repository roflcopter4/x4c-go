package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/ini.v1"

	"github.com/roflcopter4/x4c/util"
)

var Ini_Data map[string]string = nil

func init() {
	var usekey string
	machine := runtime.GOOS
	switch machine {
	case "windows":
		usekey = "windows"
	default:
		usekey = "unix"
	}

	if Ini_Data == nil {
		Ini_Data = parse_ini(usekey)
	}
}

func parse_ini(usekey string) map[string]string {
	var (
		cfg     *ini.File
		err     error
		iniFile string
	)

	if iniFile, err = filepath.Abs(os.Args[0]); err != nil {
		panic(err)
	}
	iniFile = filepath.Join(filepath.Dir(iniFile), "x4c.ini")

	if cfg, err = ini.Load(iniFile); err != nil {
		cfg = ini.Empty()
		cfg.Section("paths").Key("windows").SetValue("")
		cfg.Section("paths").Key("unix").SetValue("")
		cfg.SaveTo(iniFile)
		util.Die(1, "Uninitialized ini file.")
	}

	data := map[string]string{
		usekey:      cfg.Section("paths").Key(usekey).String(),
		"libraries": cfg.Section("extra").Key("libraries").String(),
		"aiscripts": cfg.Section("extra").Key("aiscripts").String(),
		"md":        cfg.Section("extra").Key("md").String(),
		"t":         cfg.Section("extra").Key("t").String(),
	}

	if data[usekey] == "" {
		util.Die(1, "Uninitialized ini file.")
	}

	for k, v := range data {
		if v == "" {
			data[k] = filepath.Join(data[usekey], k)
		}
	}

	return data
}
