package p2p

//import (
//	"fmt"
//	"sync"
//	"sync/atomic"
//	"testing"
//
//	"github.com/ElrondNetwork/elrond-go-sandbox/p2p"
//	"github.com/ElrondNetwork/elrond-go-sandbox/p2p/mock"
//	"github.com/pkg/errors"
//	"github.com/stretchr/testify/assert"
//)
//
//var mockMarshalizer = mock.MarshalizerMock{}
//
//func TestTopic_AddEventHandler_Nil_ShouldNotAddHandler(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	topic.AddEventHandler(nil)
//
//	assert.Equal(t, len(topic.eventBus), 0)
//}
//
//func TestTopic_AddEventHandler_WithARealFunc_ShouldWork(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	topic.AddEventHandler(func(name string, data interface{}) {
//
//	})
//
//	assert.Equal(t, len(topic.eventBus), 1)
//}
//
//func TestTopic_NewMessageReceived_NilMsg_ShouldErr(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	err := topic.NewMessageReceived(nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestTopic_NewMessageReceived_MarshalizerFails_ShouldErr(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	topic.marsh.(*mock.MarshalizerMock).Fail = true
//	defer func() {
//		topic.marsh.(*mock.MarshalizerMock).Fail = false
//	}()
//
//	err := topic.NewMessageReceived(&p2p.Message{})
//
//	assert.NotNil(t, err)
//}
//
//func TestTopic_NewMessageReceived_OKMsg_ShouldWork(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//
//	cnt := int32(0)
//	//attach event handler
//	topic.AddEventHandler(func(name string, data interface{}) {
//		assert.Equal(t, name, "test")
//
//		switch data.(type) {
//		case *string:
//			atomic.AddInt32(&cnt, 1)
//		default:
//			assert.Fail(t, "The data should have been string!")
//		}
//
//		wg.Done()
//	})
//
//	//create a new Message
//	buff, err := topic.marsh.Marshal("a string")
//	assert.Nil(t, err)
//
//	mes := p2p.NewMessage("aaa", buff, &mockMarshalizer)
//	topic.NewMessageReceived(mes)
//
//	//wait for the go routine to finish
//	wg.Wait()
//
//	assert.Equal(t, atomic.LoadInt32(&cnt), int32(1))
//}
//
//func TestTopic_Broadcast_NilData_ShouldErr(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	err := topic.Broadcast(nil, "", nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestTopic_Broadcast_MarshalizerFails_ShouldErr(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	topic.marsh.(*mock.MarshalizerMock).Fail = true
//	defer func() {
//		topic.marsh.(*mock.MarshalizerMock).Fail = false
//	}()
//
//	err := topic.Broadcast("a string", "", nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestTopic_Broadcast_NoOneToSend_ShouldErr(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	err := topic.Broadcast("a string", "", nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestTopic_Broadcast_SendOK_ShouldWork(t *testing.T) {
//	topic := NewTopic("test", "", &mockMarshalizer)
//
//	topic.OnNeedToSendMessage = func(topic string, buff []byte) error {
//		if topic != "test" {
//			return errors.New("should have been test")
//		}
//
//		if buff == nil {
//			return errors.New("should have not been nil")
//		}
//
//		fmt.Printf("Message: %v\n", buff)
//		return nil
//	}
//
//	err := topic.Broadcast("a string", "", nil)
//	assert.Nil(t, err)
//}
