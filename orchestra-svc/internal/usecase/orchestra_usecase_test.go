package usecase

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/cache"
	mockdb "orchestra-svc/internal/repository/mock"
	"orchestra-svc/internal/repository/sqlc"
	producer2 "orchestra-svc/pkg/producer"
	"testing"
)

var producerTest *producer2.KafkaProducer

func init() {
	producerTest, _ = producer2.NewKafkaProducer([]string{"localhost:29092"}, "orchestra-topic-test")
}

func TestOrchestraUsecase_ProcessWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMocks  func()
		input       event.GlobalEvent[any, any]
		expectedErr error
	}{
		{
			name: "Error finding workflow",
			setupMocks: func() {
				store.EXPECT().FindWorkflowByType(ctx, "order_created").Return(sqlc.Workflow{}, fmt.Errorf("workflow not found"))

				// Ensure CreateProcessLog is also mocked if it is used
				store.EXPECT().CreateProcessLog(ctx, gomock.Any()).Return(nil)
			},
			input: event.GlobalEvent[any, any]{
				EventType:  "order_created",
				State:      "ORDER_CREATED",
				InstanceID: "instance-001",
				EventID:    "event-001",
				StatusCode: 500,
				Payload: event.BasePayload[any, any]{
					Response: map[string]any{"key": "value"},
				},
			},
			expectedErr: fmt.Errorf("find workflow: %w", fmt.Errorf("workflow not found")),
		},
		// Additional test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := uc.ProcessWorkflow(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrchestraUsecase_ProcessWorkflow_ErrorFindingWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	eventMsg := event.GlobalEvent[any, any]{
		EventType:  "order_created",
		State:      "ORDER_CREATED",
		InstanceID: "instance-001",
		EventID:    "event-001",
		StatusCode: 500,
		Payload:    event.BasePayload[any, any]{},
	}

	store.EXPECT().FindWorkflowByType(ctx, "order_created").Return(sqlc.Workflow{}, fmt.Errorf("workflow not found"))
	store.EXPECT().CreateProcessLog(ctx, gomock.Any()).Return(nil)

	err := uc.ProcessWorkflow(ctx, eventMsg)
	expectedErr := fmt.Errorf("find workflow: %w", fmt.Errorf("workflow not found"))
	assert.Error(t, err)
	assert.Equal(t, expectedErr.Error(), err.Error())
}

func TestOrchestraUsecase_ProcessWorkflow_ErrorHandlingInstanceStep(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	eventMsg := event.GlobalEvent[any, any]{
		EventType:  "order_created",
		State:      "ORDER_CREATED",
		InstanceID: "instance-001",
		EventID:    "event-001",
		StatusCode: 500,
		Payload:    event.BasePayload[any, any]{},
	}

	store.EXPECT().FindWorkflowInstanceByID(ctx, "instance-001").Return(sqlc.WorkflowInstance{}, fmt.Errorf("error"))
	store.EXPECT().FindWorkflowByType(ctx, "order_created").Return(sqlc.Workflow{}, nil)
	store.EXPECT().CreateProcessLog(ctx, gomock.Any()).Return(nil)
	store.EXPECT().FindInstanceStepByEventID(ctx, "event-001").Return(sqlc.WorkflowInstanceStep{}, fmt.Errorf("error finding instance step"))

	err := uc.ProcessWorkflow(ctx, eventMsg)
	assert.Error(t, err)
}
func TestOrchestraUsecase_getCachePayload_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	instanceID := "instance-001"
	source := "source-1"
	response := "response-1"

	result, _ := uc.getCachePayload(instanceID, source, response)
	expected := map[string]any{"source-1": "response-1"}

	assert.Equal(t, expected, result)
}
func TestOrchestraUsecase_getCachePayload_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	instanceID := "instance-002"
	source := "source-2"
	response := "response-2"
	cacher.Set(instanceID, map[string]any{"source-2": "old-response"})

	result, _ := uc.getCachePayload(instanceID, source, response)
	expected := map[string]any{"source-2": "response-2"}

	assert.Equal(t, expected, result)
}
func TestOrchestraUsecase_handleInstanceStep_ErrorFindingInstanceStep(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	eventMsg := event.GlobalEvent[any, any]{
		EventID: "event-001",
	}

	store.EXPECT().FindInstanceStepByEventID(ctx, "event-001").Return(sqlc.WorkflowInstanceStep{}, fmt.Errorf("error"))

	err := uc.handleInstanceStep(ctx, eventMsg)
	assert.Error(t, err)
}

func TestOrchestraUsecase_processStep_ErrorCreatingWorkflowInstanceStep(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	eventMsg := event.GlobalEvent[any, any]{}
	instance := sqlc.WorkflowInstance{}
	step := sqlc.FindStepsByTypeAndStateRow{}
	cachePayload := map[string]any{}

	store.EXPECT().FindPayloadKeysByStepID(ctx, gomock.Any()).Return([]string{}, nil)
	store.EXPECT().CreateWorkflowInstanceStep(ctx, gomock.Any()).Return(sqlc.WorkflowInstanceStep{}, fmt.Errorf("error"))

	err := uc.processStep(ctx, eventMsg, instance, step, cachePayload)
	assert.Error(t, err)
}

func TestOrchestraUsecase_mergePayloads_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	keys := []string{"key1", "key2"}
	cachePayload := map[string]any{
		"key1": map[string]any{"subkey1": "value1"},
		"key2": map[string]any{"subkey2": "value2"},
	}

	expected := map[string]any{
		"subkey1": "value1",
		"subkey2": "value2",
	}

	result, err := uc.mergePayloads(keys, cachePayload)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestOrchestraUsecase_createWorkflowInstanceStep_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	gevent := event.GlobalEvent[any, any]{}
	step := sqlc.FindStepsByTypeAndStateRow{}
	eventMessage := []byte{}

	store.EXPECT().CreateWorkflowInstanceStep(ctx, gomock.Any()).Return(sqlc.WorkflowInstanceStep{}, fmt.Errorf("error"))

	err := uc.createWorkflowInstanceStep(ctx, gevent, step, eventMessage)
	assert.Error(t, err)
}
func TestOrchestraUsecase_logDB_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	ctx := context.Background()
	eventMsg := event.GlobalEvent[any, any]{}

	store.EXPECT().CreateProcessLog(ctx, gomock.Any()).Return(fmt.Errorf("error"))

	err := uc.logDB(ctx, eventMsg)
	assert.Error(t, err)
}

func TestOrchestraUsecase_getCachePayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	cacher := cache.NewPayloadCache()
	uc := NewOrchestraUsecase(store, producerTest, cacher)

	tests := []struct {
		name       string
		instanceID string
		source     string
		response   any
		setupCache func()
		expected   map[string]any
	}{
		{
			name:       "Cache miss",
			instanceID: "instance-001",
			source:     "source-1",
			response:   "response-1",
			setupCache: func() {
			},
			expected: map[string]any{"source-1": "response-1"},
		},
		{
			name:       "Cache hit",
			instanceID: "instance-002",
			source:     "source-2",
			response:   "response-2",
			setupCache: func() {
			},
			expected: map[string]any{"source-2": "response-2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupCache()

			result, _ := uc.getCachePayload(tt.instanceID, tt.source, tt.response)

			assert.Equal(t, tt.expected, result)
		})
	}
}
