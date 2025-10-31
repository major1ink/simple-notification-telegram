package yaml

type LoggerConfig struct {
	LogLevel   string `yaml:"logLevel"`
	LogDir     string `yaml:"logDir"`
	LogMode    string `yaml:"logMode"`
	RewriteLog bool   `yaml:"rewriteLog"`
}

func (l *LoggerConfig) GetLevel() string {
	return l.LogLevel
}

func (l *LoggerConfig) GetLogDir() string {
	return l.LogDir
}

func (l *LoggerConfig) GetLogMode() string {
	return l.LogMode
}

func (l *LoggerConfig) GetRewriteLog() bool {
	return l.RewriteLog
}
