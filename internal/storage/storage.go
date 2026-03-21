package storage

type MemoryStorage struct {
	memory map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memory: make(map[string]string),
	}
}

func (s *MemoryStorage) Set(key, val string) {
	s.memory[key] = val
}

func (s *MemoryStorage) Get(key string) (res string, ok bool) {
	res, ok = s.memory[key]
	return
}
