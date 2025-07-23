package config

import (
	_ "embed"
	"strings"

	"github.com/spf13/viper"

	"github.com/chaitin/MonkeyCode/backend/pkg/logger"
)

//go:embed config.json.tmpl
var ConfigTmpl []byte

type Config struct {
	Debug bool `mapstructure:"debug"`

	Logger *logger.Config `mapstructure:"logger"`

	Server struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`

	Admin struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Limit    int    `mapstructure:"limit"`
	} `mapstructure:"admin"`

	Session struct {
		ExpireDay int `mapstructure:"expire_day"`
	} `mapstructure:"session"`

	Database struct {
		Master          string `mapstructure:"master"`
		Slave           string `mapstructure:"slave"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"database"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Pass     string `mapstructure:"pass"`
		DB       int    `mapstructure:"db"`
		IdleConn int    `mapstructure:"idle_conn"`
	} `mapstructure:"redis"`

	LLMProxy struct {
		Timeout              string `mapstructure:"timeout"`
		KeepAlive            string `mapstructure:"keep_alive"`
		ClientPoolSize       int    `mapstructure:"client_pool_size"`
		StreamClientPoolSize int    `mapstructure:"stream_client_pool_size"`
		RequestLogPath       string `mapstructure:"request_log_path"`
	} `mapstructure:"llm_proxy"`

	InitModel struct {
		Name string `mapstructure:"name"`
		Key  string `mapstructure:"key"`
		URL  string `mapstructure:"url"`
	} `mapstructure:"init_model"`

	Extension struct {
		Baseurl     string `mapstructure:"baseurl"`
		LimitSecond int    `mapstructure:"limit_second"`
		Limit       int    `mapstructure:"limit"`
	} `mapstructure:"extension"`

	DataReport struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"data_report"`
}

func Init() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("MONKEYCODE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("debug", false)
	v.SetDefault("logger.level", "info")
	v.SetDefault("server.addr", ":8888")
	v.SetDefault("admin.user", "admin")
	v.SetDefault("admin.password", "")
	v.SetDefault("admin.limit", 100)
	v.SetDefault("session.expire_day", 30)
	v.SetDefault("database.master", "")
	v.SetDefault("database.slave", "")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 30)
	v.SetDefault("redis.host", "monkeycode-redis")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.pass", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.idle_conn", 20)
	v.SetDefault("llm_proxy.timeout", "30s")
	v.SetDefault("llm_proxy.keep_alive", "60s")
	v.SetDefault("llm_proxy.client_pool_size", 100)
	v.SetDefault("llm_proxy.stream_client_pool_size", 5000)
	v.SetDefault("llm_proxy.request_log_path", "/app/request/logs")
	v.SetDefault("init_model.name", "")
	v.SetDefault("init_model.key", "")
	v.SetDefault("init_model.url", "")
	v.SetDefault("extension.baseurl", "https://release.baizhi.cloud")
	v.SetDefault("extension.limit", 1)
	v.SetDefault("extension.limit_second", 10)
	v.SetDefault("data_report.key", "")

	c := Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
