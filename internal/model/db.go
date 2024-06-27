package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"test/internal/config"
)

var dsn string

var db *gorm.DB
var err error

func InitDB() {
	dsn = config.Cfg.DB.User + ":" + config.Cfg.DB.Password + "@(" + config.Cfg.DB.Host + ":" + strconv.Itoa(config.Cfg.DB.Port) + ")/" + config.Cfg.DB.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&BlackList{})
	if err != nil {
		panic(err)
	}
}
func InsertUser(user *User) error {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return err
	//}
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func GetUserByEmail(email string) (*User, error) {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return nil, err
	//}
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func UpdateUser(user *User) error {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return err
	//}
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func DeleteUser(user *User) error {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return err
	//}
	if err := db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
func IsUserExist(email string) bool {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return false
	//}
	var count int64
	db.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
func InsertBlackList(blackList *BlackList) error {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return err
	//}
	if err := db.Create(blackList).Error; err != nil {
		return err
	}
	return nil
}

func IsTokenBlackListed(token string) bool {
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return false
	//}
	var count int64
	db.Model(&BlackList{}).Where("token = ?", token).Count(&count)
	return count > 0
}
