package project

func init() {
	content["/server/conf/config.toml"] = tomlConf()
	content["/server/conf/config.yaml"] = yamlConf()
	content["/server/main.go"] = mainTemplate()
}

func tomlConf() string {
	return `[db]
addr = "root:123123@tcp(127.0.0.1:3306)/xxxx?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s"
max_open_conns = 16
max_idle_conns = 8
conn_max_life_time = 300

[redis]
#地址
addr = "127.0.0.1:6379"
#密码
password = ""
#redis 库
db = 0
#重试次数, 默认不重试
max_retries = 0
#连接池大小
pool_size = 32
#读取超时时间 3秒
read_timeout = 3
#写入超时时间 3秒
write_timeout = 3
#连接空闲时间 300秒
idle_timeout = 300
#检测死连接,并清理 默认60秒
idle_check_frequency = 60
#连接最长时间，300秒
max_conn_age = 300
#如果连接池已满 等待可用连接的时间默认 8秒
pool_timeout = 8

[other]
listen_addr = ":8000"
service_name = "{{.PackageName}}"
repository_request_timeout = 10
prometheus_listen_addr = ":9090"
# "fatal" "error" "warn" "info"  "debug"
logger_level = "debug"
# shutdown_second : Elegant lying off for the longest time
shutdown_second = 3	
`
}

func yamlConf() string {
	return `db:
    addr: root:123123@tcp(127.0.0.1:3306)/xxxx?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s
    max_open_conns: 16
    max_idle_conns: 8
    conn_max_life_time: 300
redis:
    addr: 127.0.0.1:6379
    password:
    db: 0
    max_retries: 0
    pool_size: 32
    read_timeout: 3
    write_timeout: 3
    idle_check_frequency: 60
    max_conn_age: 300
    pool_timeout: 8
other:
    listen_addr: :8000
    service_name: {{.PackageName}}
    repository_request_timeout: 10
    prometheus_listen_addr: :9090
    logger_level: debug
    shutdown_second: 3
`
}

func mainTemplate() string {
	return `
	//Package main generated by 'freedom new-project {{.PackageName}}'
	package main

	import (
		"time"
		"gorm.io/driver/mysql"
		"github.com/8treenet/freedom"
		_ "{{.PackagePath}}/adapter/repository" //Implicit initialization repository
		_ "{{.PackagePath}}/adapter/controller" //Implicit initialization controller
		"{{.PackagePath}}/server/conf"
		"github.com/go-redis/redis"
		"gorm.io/gorm"
		"github.com/8treenet/freedom/middleware"
		"github.com/8treenet/freedom/infra/requests"
	)
	
	func main() {
		app := freedom.NewApplication()
		/*
			installDatabase(app)
			installRedis(app)

			HTTP/2 h2c Runner
			runner := app.NewH2CRunner(conf.Get().App.Other["listen_addr"].(string))
			HTTP/2 AutoTLS Runner
			runner := app.NewAutoTLSRunner(":443", "freedom.com www.freedom.com", "freedom@163.com")
			HTTP/2 TLS Runner
			runner := app.NewTLSRunner(":443", "certFile", "keyFile")
		*/
		installMiddleware(app)
		runner := app.NewRunner(conf.Get().App.Other["listen_addr"].(string))
		//app.InstallParty("/{{.PackageName}}")
		liveness(app)
		app.Run(runner, conf.Get().App)
	}

	func installMiddleware(app freedom.Application) {
		app.InstallMiddleware(middleware.NewRecover())
		app.InstallMiddleware(middleware.NewTrace("x-request-id"))
		//One Loger per request New.
		app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))
		//The middleware output of the log line.
		app.Logger().Handle(middleware.DefaultLogRowHandle)

		//Install the Prometheus middleware.
		middle := middleware.NewClientPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
		requests.InstallMiddleware(middle)
				
		//HTTP request link middleware that controls the header transmission of requests.
		app.InstallBusMiddleware(middleware.NewBusFilter())
	}
	
	func installDatabase(app freedom.Application) {
		app.InstallDB(func() interface{} {
			conf := conf.Get().DB
			db, err := gorm.Open(mysql.Open(conf.Addr), &gorm.Config{})
			if err != nil {
				freedom.Logger().Fatal(err.Error())
			}
	
			sqlDB, err := db.DB()
			if err != nil {
				freedom.Logger().Fatal(err)
			}
			if err = sqlDB.Ping(); err != nil {
				freedom.Logger().Fatal(err)
			}

			sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
			sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifeTime) * time.Second)
			return db
		})
	}
	
	func installRedis(app freedom.Application) {
		app.InstallRedis(func() (client redis.Cmdable) {
			cfg := conf.Get().Redis
			opt := &redis.Options{
				Addr:               cfg.Addr,
				Password:           cfg.Password,
				DB:                 cfg.DB,
				MaxRetries:         cfg.MaxRetries,
				PoolSize:           cfg.PoolSize,
				ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
				WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
				IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
				IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
				MaxConnAge:         time.Duration(cfg.MaxConnAge) * time.Second,
				PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
			}
			redisClient := redis.NewClient(opt)
			if e := redisClient.Ping().Err(); e != nil {
				freedom.Logger().Fatal(e.Error())
			}
			client = redisClient
			return
		})
	}

	func liveness(app freedom.Application) {
		app.Iris().Get("/ping", func(ctx freedom.Context) {
			ctx.WriteString("pong")
		})
	}
	`
}
