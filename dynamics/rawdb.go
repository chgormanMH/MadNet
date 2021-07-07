package dynamics

type rawDataBase interface {
	GetValue(key []byte) ([]byte, error)
	SetValue(key []byte, value []byte) error
}
