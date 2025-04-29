package repository

import (
	"errors"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func getQueryFromUserFilter(tx *gorm.DB, filter *user.UserFilter) *gorm.DB {
	if filter == nil {
		return tx
	}
	if filter.FirstName != nil {
		tx = tx.Where("first_name = ?", *filter.FirstName)
	}
	if filter.LastName != nil {
		tx = tx.Where("last_name = ?", *filter.LastName)
	}
	if filter.Email != nil {
		tx = tx.Where("email = ?", *filter.Email)
	}
	if filter.Phone != nil {
		tx = tx.Where("phone = ?", *filter.Phone)
	}
	return tx
}

func (r *UserRepository) FindAll(userFilter *user.UserFilter) ([]user.User, error) {
	var userModels []model.User
	query := r.db.Model(&model.User{})
	query = getQueryFromUserFilter(query, userFilter)
	if err := query.Find(&userModels).Error; err != nil {
		return nil, err
	}
	users := make([]user.User, len(userModels))
	for i, um := range userModels {
		users[i] = um.ToDomain()
	}
	return users, nil
}

func (r *UserRepository) FindById(id uint) (user.User, error) {
	var userModel model.User
	if err := r.db.First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, errors.New("user not found")
		}
		return user.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *UserRepository) Create(user user.User) error {
	userModel := model.CreateUserFromDomain(user)
	return r.db.Create(&userModel).Error
}
func (r *UserRepository) Update(user user.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user.ToNotEmptyValueMap()).Error
}
