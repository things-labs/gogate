package misc

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Usart struct {
	Name     string
	BaudRate int
	DataBit  int
	Parity   string
	StopBit  int
}

type Config struct {
	OrmDbLog bool
	Com0     Usart
}

var APPConfig = Config{
	OrmDbLog: false,
	Com0: Usart{
		Name:     "/dev/ttyS1",
		BaudRate: 115200,
		DataBit:  8,
		Parity:   "N",
		StopBit:  1,
	},
}

// ConfigInit 配置初始化
func ConfigInit() error {
	file, err := os.Open("anytool.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(&APPConfig)
}
