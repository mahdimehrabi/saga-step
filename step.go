package saga_step

type Transaction interface {
	Send() error
	Compensate() error
	OnError(err error)
}

type Step struct {
	transaction Transaction
	next        *Step
	prev        *Step
}

// NewStep creates a new step with a transaction and optional next step.
func NewStep(transaction Transaction, next *Step) *Step {
	return &Step{transaction: transaction, next: next}
}

// SetNext sets the next step.
func (s *Step) SetNext(step *Step) {
	s.next = step
}

// SetPrev sets the previous step.
func (s *Step) SetPrev(step *Step) {
	s.prev = step
}

// SetTransaction sets the transaction.
func (s *Step) SetTransaction(transaction Transaction) {
	s.transaction = transaction
}

// GetNext returns the next step.
func (s *Step) GetNext() *Step {
	return s.next
}

// GetPrev returns the previous step.
func (s *Step) GetPrev() *Step {
	return s.prev
}

// GetTransaction returns the transaction.
func (s *Step) GetTransaction() Transaction {
	return s.transaction
}
