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

func DeleteUser(id string) error {
	return database.Instance.Unscoped().Model(&database.User{}).Where("id = ?", id).Delete(&database.User{}).Error
}

type TokenQueryBuilder struct {
	Ids      []string `hsource:"query" hname:"ids"`
	UserId   uint     `hsource:"query" hname:"userId"`
	Page     int      `hsource:"query" hname:"page"`
	PageSize int      `hsource:"query" hname:"pageSize"`
	Order    string   `hsource:"query" hname:"order"`
	Preload  []string `hsource:"query" hname:"preload"`
	AppIds   []int64  `hsource:"query" hname:"appIds"`
}

func (b *TokenQueryBuilder) QueryWithCount() ([]*database.AccessToken, int64, error) {
	tokens := make([]*database.AccessToken, 0)
	var count int64
	query := database.Instance.Model(&database.AccessToken{})
	if len(b.Ids) > 0 {
		query = query.Where("id in (?)", b.Ids)
	}
	if b.UserId != 0 {
		query = query.Where("user_id = ?", b.UserId)
	}
	if b.Order != "" {
		query = query.Order(b.Order)
	}
	if len(b.Preload) > 0 {
		for _, key := range b.Preload {
			query = query.Preload(key)
		}
	}
	if len(b.AppIds) > 0 {
		query = query.Where("app_id in (?)", b.AppIds)
	}
	query = query.Where("app_id > ?", 0)
	err := query.Offset((b.Page - 1) * b.PageSize).
		Limit(b.PageSize).
		Find(&tokens).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return tokens, count, nil
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
	// clear all user's tokens
	err = database.Instance.Model(&database.AccessToken{}).Where("user_id = ?", id).Delete(&database.AccessToken{}).Error
	if err != nil {
		return err
	}
	err = database.Instance.Model(&database.RefreshToken{}).Where("user_id = ?", id).Delete(&database.AccessToken{}).Error
	if err != nil {
		return err
	}
	err = database.Instance.Model(&database.AuthorizationCode{}).Where("user_id = ?", id).Delete(&database.AuthorizationCode{}).Error
	if err != nil {
		return err
	}
	return nil
}
