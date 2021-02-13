package main

import (
	"fmt"
	"github.com/dawn1806/common"
	"github.com/dawn1806/user/domain/repository"
	service2 "github.com/dawn1806/user/domain/service"
	"github.com/dawn1806/user/handler"
	user "github.com/dawn1806/user/proto/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"strconv"
)

func main() {

	// 配置中心
	conf, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		fmt.Println("[main] config err=", err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 创建数据库连接
	mysqlInfo := common.GetMysqlFromConsul(conf, "mysql")
	db, err := gorm.Open("mysql",
		mysqlInfo.User + ":" +
		mysqlInfo.Password + "@(" +
		mysqlInfo.Host + ":" + strconv.Itoa(int(mysqlInfo.Port)) + ")/" +
		mysqlInfo.Database + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.DB err=", err)
		return
	}
	defer db.Close()
	db.SingularTable(true)

	// 数据表初始化（只执行一次）
	//repository.NewUserRepository(db).InitTable()

	// 创建服务实例
	userService := service2.NewUserService(repository.NewUserRepository(db))

	// 创建服务
	srv := micro.NewService(
		micro.Name("micro-user"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8001"),
		micro.Registry(consulRegistry),
	)

	// 初始化服务
	srv.Init()

	// 注册Handler
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserService: userService})
	if err != nil {
		fmt.Println("user.RegisterUserHandler err=", err)
		return
	}

	// 启动服务
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
