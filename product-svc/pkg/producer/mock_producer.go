package producer

import (
	"github.com/golang/mock/gomock"
	"reflect"
)

type MockKafkaProducer struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaProducerMockRecorder
}

type MockKafkaProducerMockRecorder struct {
	mock *MockKafkaProducer
}

func NewMockKafkaProducer(ctrl *gomock.Controller) *MockKafkaProducer {
	mock := &MockKafkaProducer{ctrl: ctrl}
	mock.recorder = &MockKafkaProducerMockRecorder{mock}
	return mock
}

func (m *MockKafkaProducer) EXPECT() *MockKafkaProducerMockRecorder {
	return m.recorder
}

func (m *MockKafkaProducer) SendMessage(key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockKafkaProducerMockRecorder) SendMessage(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockKafkaProducer)(nil).SendMessage), key, value)
}
