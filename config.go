package common

import (
	// "fmt"
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DTM服务地址
const (
	VmwareAddr     = "192.168.71.200"
	ConsulIP       = "192.168.71.200"
	ConsulStr      = "http://192.168.71.200:8500"
	ConsulReistStr = "192.168.71.200:8500"
	DTMServer      = "http://192.168.71.200:36789/api/dtmsvr"
	QSIP           = "192.168.71.1"
	QSBusi         = "http://192.168.71.1:6669" //注意本机IP
	productFileKey = "mysql-product"
	tradeFileKey   = "mysql-trade"
	userFileKey    = "mysql-user"
	redisFileKey   = "redis"
	QPS            = 100
)

func GetConsulConfig(url string, fileKey string) (*viper.Viper, error) {
	conf := viper.New()
	conf.AddRemoteProvider("consul", url, fileKey)
	conf.SetConfigType("json")
	err := conf.ReadRemoteConfig()
	if err != nil {
		log.Println("viper conf err:", err)
	}
	return conf, nil

	// conf := viper.New()
	// conf.AddRemoteProvider("consul", url, fileKey)
	// conf.SetConfigType("json")

	// err := conf.ReadRemoteConfig()
	// if err != nil {
	// 	log.Printf("viper配置读取错误: %v", err)
	// 	// 返回错误而不是nil Viper实例
	// 	return nil, fmt.Errorf("无法从Consul读取配置: %w", err)
	// }

	// // 验证配置是否真的读取到了内容
	// if len(conf.AllKeys()) == 0 {
	// 	return nil, fmt.Errorf("配置读取成功但未找到任何配置项")
	// }

	// return conf, nil
}

func GetMysqlFromConsul(vip *viper.Viper) (db *gorm.DB, err error) {
	// 添加配置值检查
	// host := vip.GetString("host")
	// port := vip.GetString("port")
	// user := vip.GetString("user")
	// database := vip.GetString("database")

	// log.Printf("数据库配置 - host:%s, port:%s, user:%s, database:%s", host, port, user, database)

	// // 检查必要配置是否存在
	// if host == "" || port == "" || user == "" || database == "" {
	// 	return nil, fmt.Errorf("数据库配置不完整")
	// }

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	str := vip.GetString("user") + ":" + vip.GetString("pwd") + "@tcp(" + vip.GetString("host") + ":" + vip.GetString("port") + ")/" + vip.GetString("database") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(str), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("db err:", err)
	}
	return db, nil
}

// 获取Redis配置
func GetRedisFromConsul(vip *viper.Viper) (red *redis.Client, err error) {
	red = redis.NewClient(&redis.Options{
		Addr:         vip.GetString("addr"),
		Password:     vip.GetString("password"),
		DB:           vip.GetInt("DB"),
		PoolSize:     vip.GetInt("poolSize"),
		MinIdleConns: vip.GetInt("minIdleConns"),
	})
	//集群的配法
	// clusterClients := redis.NewClusterClient(
	// 	&redis.ClusterOptions{
	// 		Addrs: []string{"vmware1Ip:port", "vmware2Ip:port", "vmware3Ip:port"},
	// 	})
	// fmt.Println(clusterClients)

	return red, nil

	// ⭐ 添加详细调试
	// fmt.Println("=== GetRedisFromConsul调试信息 ===")

	// // 打印所有可用的键
	// fmt.Printf("所有配置键: %v\n", vip.AllKeys())

	// // 逐个检查每个字段
	// addr := vip.GetString("addr")
	// password := vip.GetString("password")
	// db := vip.GetInt("DB")
	// poolSize := vip.GetInt("poolSize")
	// minIdleConn := vip.GetInt("minIdleConn")

	// fmt.Printf("addr: %s (类型: %T)\n", addr, addr)
	// fmt.Printf("password: %s\n", password)
	// fmt.Printf("DB: %d (成功读取: %v)\n", db, db != 0)
	// fmt.Printf("poolSize: %d\n", poolSize)
	// fmt.Printf("minIdleConn: %d\n", minIdleConn)

	// // 尝试不同的字段名变体
	// fmt.Println("--- 尝试不同的字段名 ---")
	// fmt.Printf("GetInt(\"db\"): %d\n", vip.GetInt("db"))                     // 小写
	// fmt.Printf("GetInt(\"Db\"): %d\n", vip.GetInt("Db"))                     // 混合
	// fmt.Printf("GetInt(\"minIdleConns\"): %d\n", vip.GetInt("minIdleConns")) // 复数

	// // 获取原始值
	// if vip.IsSet("DB") {
	// 	rawDB := vip.Get("DB")
	// 	fmt.Printf("原始DB值: %v (类型: %T)\n", rawDB, rawDB)
	// }

	// // 创建Redis客户端
	// red = redis.NewClient(&redis.Options{
	// 	Addr:         addr,
	// 	Password:     password,
	// 	DB:           db,
	// 	PoolSize:     poolSize,
	// 	MinIdleConns: minIdleConn,
	// })

	// // 测试连接
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// _, err = red.Ping(ctx).Result()
	// if err != nil {
	// 	fmt.Printf("❌ Redis连接测试失败: %v\n", err)
	// 	// 尝试使用硬编码地址
	// 	fmt.Println("尝试使用硬编码地址...")
	// 	red = redis.NewClient(&redis.Options{
	// 		Addr:         "192.168.71.200:6379",
	// 		Password:     "",
	// 		DB:           0,
	// 		PoolSize:     30,
	// 		MinIdleConns: 30,
	// 	})

	// 	_, err = red.Ping(ctx).Result()
	// 	if err != nil {
	// 		fmt.Printf("❌ 硬编码连接也失败: %v\n", err)
	// 	} else {
	// 		fmt.Println("✅ 硬编码连接成功")
	// 	}
	// } else {
	// 	fmt.Println("✅ Redis连接成功")
	// }

	// return red, err

}

// 设置用户登录信息
func SetUserToken(red *redis.Client, key string, val []byte, timeTTL time.Duration) {
	red.Set(context.Background(), key, val, timeTTL)
}

// 获取用户登录信息
func GetUserToken(red *redis.Client, key string) string {
	res, err := red.Get(context.Background(), key).Result()
	if err != nil {
		log.Print("GetUserToken err", err)
	}
	return res
}
