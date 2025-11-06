package xboot

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig(path string, c interface{}) (err error) {
	v := viper.New()
	v.SetConfigFile(path)

	err = v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config %s failed [%v]", path, err)
	}

	err = v.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unmarshal config %s failed [%v]", path, err)
	}

	return nil
}
