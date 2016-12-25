package notify

import (
	"github.com/vtrifonov/http-api-mock/amqp"
	"github.com/vtrifonov/http-api-mock/definition"
)

//MockNotifier notifies the needed parties
type MockNotifier struct {
	Sender amqp.Sender
	Caller Caller
}

func NewMockNotifier() MockNotifier {
	return MockNotifier{
		Sender: amqp.MessageSender{},
		Caller: RequestCaller{},
	}
}

//Notify the needed parties
func (notifier MockNotifier) Notify(mock *definition.Mock) bool {
	success := notifier.Sender.Send(mock)
	for _, request := range mock.Notify.Http {
		successfulRequest := notifier.Caller.Call(request)
		success = success && successfulRequest
	}
	return success
}
