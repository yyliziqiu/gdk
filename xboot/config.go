package xboot

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig(path string, c interface{}) error {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config %s failed [%v]", path, err)
	}

	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("unmarshal config %s failed [%v]", path, err)
	}

	return nil
}
