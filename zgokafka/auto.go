package zgokafka

import "time"

type KafkaAuto struct {
	Delay time.Duration
	Topic []string
}

func (hb KafkaAuto) Run() {
	go func() {
		for range time.Tick(hb.Delay) { //another way to get clock signal
			hb.Service()
			time.Sleep(hb.Delay)
		}
	}()
}

func (hb KafkaAuto) Service() {
	//asrcValue := `{"status":200,"content":{"type":"normal","speaker":"xiaoi","msg":"asdf 中文"}}`
	//AsyncProducer(hb.Topic[0], asrcValue)
}

func Start() {
	a := KafkaAuto{
		Delay: 5 * time.Second,
		Topic: []string{"----"},
	}
	a.Run()

}
