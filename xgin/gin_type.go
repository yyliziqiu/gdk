package xgin

type Config struct {
	Listen           string
	KeyFile          string
	CrtFile          string
	DisableAccessLog bool
}

func (c Config) Default() Config {
	if c.Listen == "" {
		if c.CanTls() {
			c.Listen = ":443"
		} else {
			c.Listen = ":80"
		}
	}
	return c
}

func (c Config) CanTls() bool {
	return c.KeyFile != "" && c.CrtFile != ""
}
