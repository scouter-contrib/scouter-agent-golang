package netdata

type AsyncServiceStep struct {
	SingleStep
	Txid int64
	Hash int32
	Elapsed int32
}

func NewAsyncServiceStep() *AsyncServiceStep {
	step := new(AsyncServiceStep)
	return step
}

func (s *AsyncServiceStep) GetStepType() byte {
	return THREAD_CALL_POSSIBLE
}

func (s *AsyncServiceStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteDecimal(s.Txid)
	_, err = out.WriteDecimal32(s.Hash)
	_, err = out.WriteDecimal32(s.Elapsed)
	_, err = out.WriteUInt8(1)
	return err
}
