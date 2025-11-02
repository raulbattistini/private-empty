package server

type YamlCfg struct {
	Server struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		MaxStoreSize int    `yaml:"max_store_size"`
		Production   bool   `yaml:"production"`
	} `yaml:"server"`
}
