package xkvs

func (k Kvs) Enabled() bool {
	return k.B("enabled", false)
}

func (k Kvs) Disabled() bool {
	return k.B("disabled", false)
}
