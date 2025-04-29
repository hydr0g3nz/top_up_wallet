package user

type Repository interface {
	FindAll() ([]User, error)
	FindById(id uint) (User, error)
	Create(User User) error
}
