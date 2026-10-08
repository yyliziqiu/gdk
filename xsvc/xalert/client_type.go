package ext

type SendFrom struct {
	KeyType string `json:"key_type"`
	Key     string `json:"key"`
	Args    Args   `json:"args"`
}

func newSendFrom(keyType string, key string, args Args) SendFrom {
	return SendFrom{
		KeyType: keyType,
		Key:     key,
		Args:    args,
	}
}
