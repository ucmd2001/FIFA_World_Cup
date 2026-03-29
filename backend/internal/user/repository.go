package user

import (
	"gorm.io/gorm"
)

// Repository 定義 User domain 的資料存取介面
type Repository interface {
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	CountAll() (int64, error)
	CountByUsername(username string) (int64, error)
	Create(u *User) error
	Save(u *User) error
	Update(u *User, updates map[string]interface{}) error
	Delete(id string) error
	FindAll() ([]User, error)
	FindTopN(n int) ([]User, error)
	DeductPoints(userID uint, amount int) (rowsAffected int64, err error)
	AddPoints(userID uint, amount int) error
}

// GormRepository 是 Repository 的 GORM 實作
type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) FindByUsername(username string) (*User, error) {
	var u User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormRepository) CountByUsername(username string) (int64, error) {
	var count int64
	if err := r.db.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *GormRepository) Save(u *User) error {
	return r.db.Save(u).Error
}

func (r *GormRepository) Update(u *User, updates map[string]interface{}) error {
	return r.db.Model(u).Updates(updates).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *GormRepository) FindAll() ([]User, error) {
	var users []User
	if err := r.db.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormRepository) FindTopN(n int) ([]User, error) {
	var users []User
	if err := r.db.Order("points desc").Limit(n).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormRepository) DeductPoints(userID uint, amount int) (int64, error) {
	result := r.db.Model(&User{}).
		Where("id = ? AND points >= ?", userID, amount).
		UpdateColumn("points", gorm.Expr("points - ?", amount))
	return result.RowsAffected, result.Error
}

func (r *GormRepository) AddPoints(userID uint, amount int) error {
	return r.db.Model(&User{}).
		Where("id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", amount)).Error
}
