package saga_step

import (
	"errors"
	"go.uber.org/mock/gomock"
	mock_saga_step "saga-step/mock"
	"testing"
	"time"
)

func TestNewSaga_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var testCases = []struct {
		Name         string
		Transactions func() []*mock_saga_step.MockTransaction
		Error        error
		RetryCount   int
		RetryDelay   time.Duration
	}{
		{
			Name: "success",
			Transactions: func() []*mock_saga_step.MockTransaction {
				transaction1 := mock_saga_step.NewMockTransaction(ctrl)
				transaction1.EXPECT().Send().Return(nil).Times(1)
				transaction2 := mock_saga_step.NewMockTransaction(ctrl)
				transaction2.EXPECT().Send().Return(nil).Times(1)
				transaction3 := mock_saga_step.NewMockTransaction(ctrl)
				transaction3.EXPECT().Send().Return(nil).Times(1)
				transaction4 := mock_saga_step.NewMockTransaction(ctrl)
				transaction4.EXPECT().Send().Return(nil).Times(1)
				return []*mock_saga_step.MockTransaction{
					transaction1, transaction2, transaction3, transaction4,
				}
			},
			Error: nil,
		},
		{
			Name: "err_steps_not_set",
			Transactions: func() []*mock_saga_step.MockTransaction {
				return []*mock_saga_step.MockTransaction{}
			},
			Error:      ErrStepNotSet,
			RetryCount: 10,
			RetryDelay: time.Millisecond * 1,
		},
		{
			Name: "err_compensate_from_last_step",
			Transactions: func() []*mock_saga_step.MockTransaction {
				transaction1 := mock_saga_step.NewMockTransaction(ctrl)
				transaction1.EXPECT().Send().Return(nil).Times(1)
				transaction1.EXPECT().Compensate().Return(errors.New("error")).Times(10)
				transaction1.EXPECT().OnError(gomock.Any()).Times(10)

				transaction2 := mock_saga_step.NewMockTransaction(ctrl)
				transaction2.EXPECT().Send().Return(nil).Times(1)
				transaction2.EXPECT().Compensate().Return(nil).Times(1)

				transaction3 := mock_saga_step.NewMockTransaction(ctrl)
				transaction3.EXPECT().Send().Return(nil).Times(1)
				transaction3.EXPECT().Compensate().Return(nil).Times(1)

				transaction4 := mock_saga_step.NewMockTransaction(ctrl)
				transaction4.EXPECT().Send().Return(errors.New("error")).Times(10)
				transaction4.EXPECT().OnError(gomock.Any()).Times(10)
				return []*mock_saga_step.MockTransaction{
					transaction1, transaction2, transaction3, transaction4,
				}
			},
			Error:      nil,
			RetryCount: 10,
			RetryDelay: time.Millisecond * 1,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			transactions := testCase.Transactions()
			saga := NewSaga(10, time.Millisecond)
			for _, transaction := range transactions {
				saga.PushStep(&Step{
					transaction: transaction,
				})
			}
			err := saga.Execute()
			if !errors.Is(err, testCase.Error) {
				t.Fatalf("got error %v, want %v", testCase.Error, err)
			}
			for _, transaction := range transactions {
				transaction.EXPECT()
			}
		})
	}
}
