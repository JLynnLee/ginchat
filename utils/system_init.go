package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB
var RDS *redis.Client

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	PublishKey = "websocket"
)

func RandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("====================config mysql:", viper.Get("mysql"))
}

func InitMySql() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
}

func InitRedis() {
	fmt.Println("===================", viper.GetString("redis.addr"))
	ctx := context.Background()
	RDS = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
		DB:           viper.GetInt("redis.DB"),
	})
	pong, err := RDS.Ping(ctx).Result()
	if err != nil {
		fmt.Println("init redis failed ........", err)
	} else {
		fmt.Println("redis connect .....", pong)
	}
}

// 发送消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("-------------Publish")
	err = RDS.Publish(ctx, channel, msg).Err()
	return err
}

// 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	fmt.Println("-------------Subscribe", ctx.Value(""))
	sub := RDS.Subscribe(ctx, channel)
	fmt.Println("redis-Subscribe---------", sub)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("redis-Subscribe====", err)
		return "", err
	}
	fmt.Println("-------------msg===", msg.Payload)

	return msg.Payload, err
}
