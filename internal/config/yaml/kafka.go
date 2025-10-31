package yaml

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

func (k *KafkaConfig) GetBrokers() []string {
	return k.Brokers
}
