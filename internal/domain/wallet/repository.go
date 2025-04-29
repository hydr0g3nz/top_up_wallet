package wallet

type Repository interface {
	Create(wallet Wallet) error
	Update(wallet Wallet) error
}
