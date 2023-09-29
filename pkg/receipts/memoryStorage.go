package receipts

import "errors"

var (
	ErrNotFound = errors.New("receipt id not found")
)

type MemStore struct {
	list map[string]ReceiptProperties
}

func NewMemStore() *MemStore {
	list := make(map[string]ReceiptProperties)
	return &MemStore{
		list,
	}
}

func (m MemStore) PostReceipt(name string, receipt Receipt, points int) error {
	m.list[name] = ReceiptProperties{receipt, points}
	return nil
}

func (m MemStore) GetPoints(name string) (int, error) {

	if val, ok := m.list[name]; ok {
		return val.points, nil
	}

	return 0, ErrNotFound
}
