package hive

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/grav/grav"
	"github.com/suborbital/grav/testutil"
)

const msgTypeTester = "hive.test"
const msgTypeNil = "hive.testnil"

// to test jobs listening to a Grav message
type msgRunner struct{}

func (m *msgRunner) Run(job Job, do DoFunc) (interface{}, error) {
	name := string(job.Bytes())

	reply := grav.NewMsg(msgTypeTester, []byte(fmt.Sprintf("hello, %s", name)))

	return reply, nil
}

func (m *msgRunner) OnStart() error { return nil }

// to test jobs with a nil result
type nilRunner struct{}

func (m *nilRunner) Run(job Job, do DoFunc) (interface{}, error) {
	return nil, nil
}

func (m *nilRunner) OnStart() error { return nil }

func TestHandleMessage(t *testing.T) {
	hive := New()
	g := grav.New()

	hive.HandleMsg(g.Connect(), msgTypeTester, &msgRunner{})

	counter := testutil.NewAsyncCounter(10)

	sender := g.Connect()

	sender.OnType(msgTypeTester, func(msg grav.Message) error {
		counter.Count()
		return nil
	})

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

	sender.OnType(msgTypeTester, func(msg grav.Message) error {
		counter.Count()
		return nil
	})

	for i := 0; i < 9876; i++ {
		sender.Send(grav.NewMsg(msgTypeTester, []byte("charlie brown")))
	}

	if err := counter.Wait(9876, 1); err != nil {
		t.Error(errors.Wrap(err, "failed to counter.Wait"))
	}
}

func TestHandleMessageNilResult(t *testing.T) {
	hive := New()
	g := grav.New()

	hive.HandleMsg(g.Connect(), msgTypeNil, &nilRunner{})

	counter := testutil.NewAsyncCounter(10)

	pod := g.Connect()

	pod.OnType(MsgTypeHiveNilResult, func(msg grav.Message) error {
		counter.Count()
		return nil
	})

	for i := 0; i < 5; i++ {
		pod.Send(grav.NewMsg(msgTypeNil, []byte("hi")))
	}

	if err := counter.Wait(5, 1); err != nil {
		t.Error(errors.Wrap(err, "failed to counter.Wait"))
	}
}
