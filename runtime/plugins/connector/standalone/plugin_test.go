package standalone

import (
	"context"
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	connectorPlugin "eventcenter-go/runtime/plugins/connector"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestMessage(t *testing.T) {
	connPlugin := plugins.Get(plugins.TypeConnector, plugins.NameConnectorStandalone).(connectorPlugin.Plugin)
	_ = connPlugin.Init(nil)

	c, _ := connPlugin.Consumer()
	c.RegisterHandler(&connector.EventHandler{
		Handler: func(event *cloudevents.Event) (err error) {
			log.Println(event.ID())
			return nil
		},
	})
	_ = c.Subscribe("test")

	done := make(chan bool)

	p, _ := connPlugin.Producer()
	for i := 0; i < 20; i++ {
		ctx := context.TODO()

		event := cloudevents.NewEvent()
		event.SetID(uuid.NewString())
		event.SetSubject("test")

		_ = p.Publish(ctx, &event)
	}

	<-done
}
