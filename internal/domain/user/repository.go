package user

type Repository interface {
	FindAll() ([]User, error)
	FindById(id uint) (User, error)
	Save(User User) error
}
