package user

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// IService lists the methods the user service implements.
type IService interface {
	CreateUser(u *User) (User, error)
	ReadUser(u uuid.UUID) (t User, err error)
	UpdateUser(t *User, u uuid.UUID) error
	DeleteUser(u uuid.UUID) error
}

// IDatabase list the methods a database needs to
// implement, in order to work with user's package.
type IDatabase interface {
	InsertUser(t User) error
	FindUser(u uuid.UUID) (User, error)
	UpdateUser(t User, u uuid.UUID) error
	DeleteUser(u uuid.UUID) error
}

type service struct {
	db IDatabase
}

// NewService returns a pointer to new instance of
// the user service.
func NewService(db IDatabase) IService {
	return &service{
		db,
	}
}

// CreateUser will add a uuid and times to user and pass it along
// to insert DB method.
func (s *service) CreateUser(u *User) (User, error) {
	u.UUID = uuid.NewV4()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err := s.db.InsertUser(*u)
	if err != nil {
		return *u, err
	}

	return *u, nil
}

// ReadUser will pass the uuid to the find DB method.
func (s *service) ReadUser(u uuid.UUID) (user User, err error) {
	user, err = s.db.FindUser(u)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser will update updatedAt time and pass it
// it to DB update method.
func (s *service) UpdateUser(t *User, u uuid.UUID) error {
	t.UpdatedAt = time.Now()
	err := s.db.UpdateUser(*t, u)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser will call DB delete method with the
// uuid that was passed.
func (s *service) DeleteUser(u uuid.UUID) error {
	return s.db.DeleteUser(u)
}
