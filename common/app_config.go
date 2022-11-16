package common

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//==========全局变量区 ==========//
var Logger *log.Logger
var DB *gorm.DB
var Redis *redis.Client

//========== 对象定义区 ==========//
type AppConfig struct {
	Env        string `json:"env"`
	AppName    string `yaml:"app_name" json:"app_name"`
	ServerPort int    `yaml:"server_port" json:"server_port"`

	JwtConfig        *JwtConfig        `yaml:"jwt" json:"jwt"`
	LogConfig        *LogConfig        `yaml:"log" json:"log"`
	DatasourceConfig *DatasourceConfig `yaml:"datasource" json:"datasource"`
	RedisConfig      *RedisConfig      `yaml:"redis" json:"redis"`

	initFlag bool
}

type JwtConfig struct {
	SecretKey string        `yaml:"secret_key" json:"secret_key"`
	TTL       time.Duration `yaml:"ttl" json:"ttl"`
	secret    *[]byte
}

type LogConfig struct {
	File  string `yaml:"file,omitempty" json:"file"`
	Level string `yaml:"level,omitempty" json:"level"`
}

type DatasourceConfig struct {
	URL      string `yaml:"url,omitempty" json:"url"`
	Username string `yaml:"username,omitempty" json:"username"`
	Password string `yaml:"password,omitempty" json:"password"`
}

type RedisConfig struct {
	Address  string `yaml:"address,omitempty" json:"address"`
	Password string `yaml:"password,omitempty" json:"password"`
	DB       int    `yaml:"db,omitempty" json:"db"`
}

//========== init ==========//
func (c *AppConfig) Init() error {
	if c.initFlag {
		return nil
	}

	err := c.initJwt()
	if err != nil {
		return err
	}
	err = c.initLog()
	if err != nil {
		return err
	}
	err = c.initDatasource()
	if err != nil {
		return err
	}
	err = c.initRedis()
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) initJwt() error {
	if c.JwtConfig == nil || len(c.JwtConfig.SecretKey) == 0 {
		return nil
	}

	bytes, err := base64.StdEncoding.DecodeString(c.JwtConfig.SecretKey)
	if err != nil {
		return fmt.Errorf("can't decode jwt secret_key: %s", c.JwtConfig.SecretKey)
	}
	c.JwtConfig.secret = &bytes
	return nil
}

func (c *AppConfig) initLog() error {
	var logFile = c.AppName + ".log"
	var logLevel = "debug"
	if c.LogConfig != nil {
		if len(c.LogConfig.File) != 0 {
			logFile = c.LogConfig.File
		}
		if len(c.LogConfig.Level) != 0 {
			logLevel = c.LogConfig.Level
		}
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can't open log file: %s", logFile)
	}

	Logger = &log.Logger{Out: file, Formatter: &log.JSONFormatter{}, ReportCaller: true, Level: level}
	return nil
}

func (c *AppConfig) initDatasource() error {
	if c.DatasourceConfig == nil {
		return nil
	}
	var err error
	DB, err = gorm.Open(mysql.Open(c.DatasourceConfig.Username+":"+c.DatasourceConfig.Password+"@"+c.DatasourceConfig.URL), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "haul_", // table name prefix, table for `User` would be `t_users`
		SingularTable: true,    // use singular table name, table for `User` would be `user` with this option enabled
		//	NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
	}})
	if err != nil {
		return err
	}
	pool, err := DB.DB()
	if err != nil {
		return err
	}

	DB.Config.PrepareStmt = true
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)
	return nil
}

func (c *AppConfig) initRedis() error {
	if c.RedisConfig == nil {
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisConfig.Address,
		Password: c.RedisConfig.Password,
		DB:       c.RedisConfig.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("can't connect to redis, redis config=%s", c.RedisConfig.String())
	}
	Redis = rdb
	return nil
}

//========== 对象String()定义区 ==========//
func (c AppConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c JwtConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c JwtConfig) GetSecret() *[]byte {
	return c.secret
}

func (c LogConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c DatasourceConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c RedisConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}
