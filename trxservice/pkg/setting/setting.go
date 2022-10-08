package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

var (
	Tracer opentracing.Tracer
	// 保存所有配置信息
	Conf = new(Config)
)

func Init() (err error) {
	//viper.SetConfigFile("config.yaml")
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml") // 配合远程配置中心使用,告诉viper 获取的配置信息使用什么格式去解析
	// 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")   // 还可以在工作目录中查找配置
	viper.AddConfigPath("../") // 还可以在工作目录中查找配置
	viper.AddConfigPath("./conf/")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("/app/")
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		fmt.Printf("viper.ReadInConfig failed, error : %s \n", err.Error())
		return err
	}
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, error : %s \n", err.Error())
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生了修改")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, error : %s \n", err.Error())
			return
		}
		fmt.Printf("cfg:%v\n", Conf)
	})
	return
}

type Config struct {
	App       `mapstructure:"app"`
	Log       `mapstructure:"log"`
	DB        `mapstructure:"db"`
	Redis     `mapstructure:"redis"`
	TokenList map[string]Token `mapstructure:"tokenList" json:"tokenList"`
	Metrics   `mapstructure:"metrics"`
	Trace     `mapstructure:"trace"`
}

type App struct {
	Name      string   `mapstructure:"name"`
	RunMode   string   `mapstructure:"run_mode"`
	HttpPort  int      `mapstructure:"http_port"`
	GrpcPort  int      `mapstructure:"grpc_port"`
	Node_Addr []string `mapstructure:"node_addr"`
}
type Log struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type DB struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DbName          string `mapstructure:"db_name"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_lifetime"`
}
type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type Token struct {
	Name         string `mapstructure:"name" json:"name"`
	Decimal      uint   `mapstructure:"decimal" json:"decimal"`
	ContractAddr string `mapstructure:"contractAddr" json:"contractAddr"`
}

type Metrics struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"service_name"`
}

type Trace struct {
	Enable      bool   `mapstructure:"enable"`
	Endpoint    string `mapstructure:"endpoint"`
	ServiceName string `mapstructure:"service_name"`
	LogSpans    bool   `mapstructure:"log_spans"`
}
