package service

import (
	"github.com/projectxpolaris/youauth/database"
	"golang.org/x/crypto/bcrypt"
)

type UserQueryBuilder struct {
	Ids        []string `hsource:"query" hname:"ids"`
	NameSearch string   `hsource:"query" hname:"search"`
	Name       string   `hsource:"query" hname:"name"`
	Page       int      `hsource:"query" hname:"page"`
	PageSize   int      `hsource:"query" hname:"pageSize"`
	Order      string   `hsource:"query" hname:"order"`
}

func (b *UserQueryBuilder) GetDataAndCount() ([]*database.User, int64, error) {
	users := make([]*database.User, 0)
	var count int64
	query := database.Instance.Model(&database.User{})
	if len(b.Ids) > 0 {
		query = query.Where("id in (?)", b.Ids)
	}
	if b.NameSearch != "" {
		query = query.Where("name like ?", "%"+b.NameSearch+"%")
	}
	if b.Name != "" {
		query = query.Where("name = ?", b.Name)
	}
	if b.Order != "" {
		query = query.Order(b.Order)
	}
	err := query.Offset((b.Page - 1) * b.PageSize).
		Limit(b.PageSize).
		Find(&users).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func GetUserById(id string) (*database.User, error) {
	user := &database.User{}
	err := database.Instance.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*database.User, error) {
	user := &database.User{}
	err := database.Instance.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func DeleteUser(id string) error {
	return database.Instance.Unscoped().Model(&database.User{}).Where("id = ?", id).Delete(&database.User{}).Error
}

func ChangePassword(id uint, oldPassword, password string) error {
	// get user
	var user database.User
	err := database.Instance.Model(&database.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return err
	}
	rawPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(rawPassword)
	err = database.Instance.Save(&user).Error
	if err != nil {
		return err
	}
	err = database.Instance.Model(&database.AuthorizationCode{}).Where("user_id = ?", id).Delete(&database.AuthorizationCode{}).Error
	if err != nil {
		return err
	}
	return nil
}
