package dynamics

type rawDataBase interface {
	GetValue(key []byte) ([]byte, error)
	SetValue(key []byte, value []byte) error
}

/*
type rawDataBase struct {
	db     *badger.DB
	logger *logrus.Logger
}

func (r *rawDataBase) GetValue(key []byte) ([]byte, error) {
	panic("not implemented")
}

func (r *rawDataBase) SetValue(key []byte, value []byte) error {
	panic("not implemented")
}
*/
