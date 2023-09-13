package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xkafka"

	pkafka "github.com/NetEase-Media/easy-ngo/app/plugins/plugin_xkafka"
)

const topic = "test"

var (
	consumer    *xkafka.Consumer
	producer    *xkafka.Producer
	messageChan = make(chan xkafka.ConsumerMessage, 1000)
)

func main() {
	app := app.New()
	if err := app.Init(); err != nil {
		panic(err)
	}
	consumer = pkafka.GetConsumer()
	consumer.AddListener(topic, &listener{})
	consumer.Start()

	producer = pkafka.GetProducer()
	for i := 0; i < 98; i++ {
		producer.Send(topic, "hello world!"+strconv.Itoa(i), nil)
	}
	go func() {
		for {
			r := <-messageChan
			fmt.Print(r.Value, " ", r.Partition, " ", r.Offset)
		}
	}()
	time.Sleep(10 * time.Second)

}

type listener struct {
	xkafka.Listener
}

func (l *listener) Listen(message xkafka.ConsumerMessage, ack *xkafka.Acknowledgment) {
	messageChan <- message
	ack.Acknowledge()
}
