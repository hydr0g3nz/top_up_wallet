package user

type Repository interface {
	FindAll(*UserFilter) ([]User, error)
	FindById(id uint) (User, error)
	Create(User User) error
}
