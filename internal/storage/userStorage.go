package storage

type UserStorage struct {
	storage Storer
}

func NewUserStorage(storage Storer) *UserStorage {
	return &UserStorage{storage: storage}
}
