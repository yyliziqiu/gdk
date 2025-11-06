package xgin

type Config struct {
	Listen           string
	KeyFile          string
	CertFile         string
	DisableAccessLog bool
}

func (c Config) Default() Config {
	if c.Listen == "" {
		if c.TlsEnabled() {
			c.Listen = ":443"
		} else {
			c.Listen = ":80"
		}
	}
	return c
}

func (c Config) TlsEnabled() bool {
	return c.KeyFile != "" && c.CertFile != ""
}
