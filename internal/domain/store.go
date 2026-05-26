package domain

type Store interface {
	Save(rec ConfigRecord) (int64, error)
	Find(id int64) (ConfigRecord, error)
	List(offset, limit int) ([]ConfigRecord, error)
}
