package main

import (
	"log"
	"login/kitex_gen/user/registerservice"
	"login/serve/model"
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
		return
	}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
		return
	}
	dsn := "root:123456@tcp(localhost:3306)/kitex?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal(err)
		return
	}
	svr := registerservice.NewServer(&RegisterServiceImpl{DB: db}, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "userservice"}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{ //限流
			MaxConnections: 1000,
			MaxQPS:         2000,
		}))
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
