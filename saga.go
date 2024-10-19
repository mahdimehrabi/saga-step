package saga_step

import (
	"errors"
	"time"
)

var (
	ErrStepNotSet = errors.New("step not set")
)

type Saga struct {
	Head       *Step
	Tail       *Step
	retryCount uint8
	duration   time.Duration
}

func NewSaga(retryCount uint8, duration time.Duration) *Saga {
	return &Saga{
		Head:       nil,
		retryCount: retryCount,
		duration:   duration,
	}
}

func (s *Saga) PushStep(step *Step) {
	if s.Head == nil {
		s.Head = step
		s.Tail = step
		return
	}
	s.Tail.SetNext(step)
	step.SetPrev(s.Tail)
	s.Tail = step
}

func (s *Saga) RemoveStep() {
	s.Tail.prev.next = nil
}

func (s *Saga) Execute() error {
	step := s.Head
	if step == nil {
		return ErrStepNotSet
	}
StepLoop:
	for step != nil {
		for i := uint8(0); i < s.retryCount; i++ {
			if err := step.GetTransaction().Send(); err != nil {
				step.GetTransaction().OnError(err)
				time.Sleep(s.duration)
				continue
			}
			step = step.GetNext()
			continue StepLoop
		}
		step = step.GetPrev()
	CompensateLoop:
		for step != nil {
			for i := uint8(0); i < s.retryCount; i++ {
				if err := step.GetTransaction().Compensate(); err != nil {
					step.GetTransaction().OnError(err)
					time.Sleep(s.duration)
					continue
				}
				step = step.GetPrev()
				continue CompensateLoop
			}
			step = step.GetPrev()
			continue CompensateLoop
		}
	}

	return nil
}
