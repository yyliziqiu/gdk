package ext

type Request struct {
	KeyType string `json:"key_type"`
	Key     string `json:"key"`
	Args    Args   `json:"args"`
}

func req(keyType string, key string, args Args) Request {
	if args == nil {
		args = Args{}
	}

	return Request{
		KeyType: keyType,
		Key:     key,
		Args:    args,
	}
}
