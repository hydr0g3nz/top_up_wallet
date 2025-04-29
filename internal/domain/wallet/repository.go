package wallet

type Repository interface {
	FindAll() ([]Wallet, error)
	FindById(id uint) (Wallet, error)
	Save(wallet Wallet) error
}
