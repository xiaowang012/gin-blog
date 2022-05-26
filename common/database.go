package common

import (
	"encoding/json"
	"fmt"
	"gin-blog/models"
	"io/ioutil"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Parameters 读取json配置数据的结构体
type Parameters struct {
	DriverName string `json:"drivername"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Database   string `json:"database"`
}

var DB *gorm.DB

// InitDB 数据库连接初始化
func InitDB() {
	//读取config.json配置文件
	jsonFile, err := os.Open("./common/config.json")
	if err != nil {
		fmt.Println("打开config.json文件错误: " + err.Error())
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("读取配置文件config.json错误: " + err.Error())
		return
	}
	var parameters Parameters
	json.Unmarshal(jsonData, &parameters)
	driverName := parameters.DriverName
	Username := parameters.Username
	Password := parameters.Password
	Host := parameters.Host
	Port := parameters.Port
	Database := parameters.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", Username, Password, Host, Port, Database)
	//连接数据库mysql
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: driverName,
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql连接错误: " + err.Error())
		return
	}
	//自动创建表结构
	db.AutoMigrate(&models.Users{}, &models.UserGroup{}, &models.Permissions{}, &models.Group{}, &models.Comments{}, &models.Articles{}, &models.UserLikes{}, &models.MessageBoard{})
	DB = db

}

func GetDB() *gorm.DB {
	return DB
}
