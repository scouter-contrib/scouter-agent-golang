package netdata

type HashedMessageStep struct {
	SingleStep
	Hash int32
	Time int32
	Value int32
}

func NewHashedMessageStep(hash int32, startTime int32) *HashedMessageStep {
	step := new(HashedMessageStep)
	step.Hash = hash
	step.StartTime = startTime
	return step
}

func (s *HashedMessageStep) GetStepType() byte {
	return HASHED_MESSAGE
}

func (s *HashedMessageStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteDecimal32(s.Hash)
	_, err = out.WriteDecimal32(s.Time)
	_, err = out.WriteDecimal32(s.Value)
	return err
}
