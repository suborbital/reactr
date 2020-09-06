package hive

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/grav/grav"
	"github.com/suborbital/grav/testutil"
)

const msgTypeTester = "hive.test"

type msgRunner struct{}

func (m *msgRunner) Run(job Job, do DoFunc) (interface{}, error) {
	msg := job.Msg()
	if msg == nil {
		return nil, errors.New("not a message")
	}

	name := string(msg.Data())

	reply := grav.NewMsg(msgTypeTester, []byte(fmt.Sprintf("hello, %s", name)))

	return reply, nil
}

func (m *msgRunner) OnStart() error { return nil }

func TestHandleMessage(t *testing.T) {
	hive := New()
	g := grav.New()

	hive.HandleMsg(g.Connect(), msgTypeTester, &msgRunner{})

	counter := testutil.NewAsyncCounter(10)

	sender := g.Connect()

	sender.OnType(func(msg grav.Message) error {
		counter.Count()
		return nil
	}, msgTypeTester)

	sender.Send(grav.NewMsg(msgTypeTester, []byte("charlie brown")))

	if err := counter.Wait(1, 1); err != nil {
		t.Error(errors.Wrap(err, "failed to counter.Wait"))
	}
}

func TestHandleMessagePt2(t *testing.T) {
	hive := New()
	g := grav.New()

	hive.HandleMsg(g.Connect(), msgTypeTester, &msgRunner{})

	counter := testutil.NewAsyncCounter(10000)

	sender := g.Connect()

	sender.OnType(func(msg grav.Message) error {
		counter.Count()
		return nil
	}, msgTypeTester)

	for i := 0; i < 9876; i++ {
		sender.Send(grav.NewMsg(msgTypeTester, []byte("charlie brown")))
	}

	if err := counter.Wait(9876, 1); err != nil {
		t.Error(errors.Wrap(err, "failed to counter.Wait"))
	}
}
