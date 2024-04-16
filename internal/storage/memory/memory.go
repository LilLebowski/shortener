package memory

type Storage struct {
	URLs map[string]string
}

func Init() *Storage {
	return &Storage{
		URLs: make(map[string]string),
	}
}

func (s *Storage) Set(short string, full string) {
	s.URLs[short] = full
}

func (s *Storage) Get(short string) (string, bool) {
	value, exists := s.URLs[short]
	return value, exists
}
