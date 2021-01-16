package netdata

type MethodStep struct {
	SingleStep
	Hash int32
	Elapsed int32
	CpuTime int32
}

func NewMethodStep() *MethodStep {
	step := new(MethodStep)
	return step
}

func (s *MethodStep) GetStepType() byte {
	return METHOD
}

func (s *MethodStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteDecimal32(s.Hash)
	_, err = out.WriteDecimal32(s.Elapsed)
	_, err = out.WriteDecimal32(s.CpuTime)

	return err
}
