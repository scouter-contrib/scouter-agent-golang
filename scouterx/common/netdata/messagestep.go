package netdata

type MessageStep struct {
	SingleStep
	Message string
}

func NewMessageStep(m string, startTime int32) *MessageStep {
	step := new(MessageStep)
	step.Message = m
	step.StartTime = startTime
	return step
}

func (s *MessageStep) GetStepType() byte {
	return MESSAGE
}

func (s *MessageStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	_, err = out.WriteString(s.Message)

	return err
}
