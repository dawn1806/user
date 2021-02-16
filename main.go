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
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {

	// 创建数据库连接
	db, err := gorm.Open("mysql", common.MysqlConnection)
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
		micro.Registry(etcd.NewRegistry(
			registry.Addrs("127.0.0.1:2379"))),
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
