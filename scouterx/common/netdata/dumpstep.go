package netdata

type DumpStep struct {
	SingleStep
	Stacks []int32
	_threadId int64
	_threadName string
	_threadState string
	_lockOwnerId int64
	_lockName string
	_lockOwnerName string
}

func NewDumpStep() *DumpStep {
	step := new(DumpStep)
	return step
}

func (s *DumpStep) GetStepType() byte {
	return METHOD
}

func (s *DumpStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteInt32Array(s.Stacks)
	_, err = out.WriteInt64(s._threadId)
	_, err = out.WriteString(s._threadName)
	_, err = out.WriteString(s._threadState)
	_, err = out.WriteInt64(s._lockOwnerId)
	_, err = out.WriteString(s._lockName)
	_, err = out.WriteString(s._lockOwnerName)

	return err
}
