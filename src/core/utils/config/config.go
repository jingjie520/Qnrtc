package config

import (
	"fmt"
	"reflect"

	"github.com/Unknwon/goconfig"
)

const confFile = "conf.ini"

//Super 主程序配置
type Super struct {
	AccessKey string
	SecretKey string
	AppID     string
}

//Logrus 日志配置
type Logrus struct {
	LogPath string
	LogName string
}

//Config 配置文件
type Config struct {
	Super
	Logrus
}

//Conf 配置文件对象
var Conf *Config

func init() {
	Conf = getConf()
}

func getConf() *Config {
	if Conf != nil {
		return Conf
	}

	cfg, err := goconfig.LoadConfigFile(confFile)
	if err != nil {
		fmt.Println("加载配置文件conf.ini失败。请检查当前目录下是否存在该文件。")
		return nil
	}

	Conf = &Config{}

	initSuper(cfg)

	return Conf
}

func initSuper(cfg *goconfig.ConfigFile) {

	initBySection(cfg, "super", reflect.TypeOf(Conf.Super), reflect.ValueOf(&Conf.Super).Elem())

	initBySection(cfg, "logrus", reflect.TypeOf(Conf.Logrus), reflect.ValueOf(&Conf.Logrus).Elem())
}

func initBySection(cfg *goconfig.ConfigFile, section string, t reflect.Type, v reflect.Value) {
	for k := 0; k < t.NumField(); k++ {
		f := v.Field(k)
		key := strFirstToLow(t.Field(k).Name)

		switch f.Kind() {
		case reflect.Int64:
			f.SetInt(cfg.MustInt64(section, key, 0))
		case reflect.String:
			f.SetString(cfg.MustValue(section, key, ""))
		case reflect.Bool:
			f.SetBool(cfg.MustBool(section, key, false))
		}
	}
}

//strFirstToLow 首字母转xiao'x
func strFirstToLow(str string) string {
	if len(str) == 0 {
		return str
	}

	v := []rune(str)
	if int(v[0]) >= int('A') && int(v[0]) <= int('Z') {
		v[0] += 32
	}
	return string(v)
}
