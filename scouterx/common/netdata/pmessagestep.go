package netdata

import (
	"bytes"
)

const delimETX byte = 3
const PMSG_DEBUG PMessageLevel = 0
const PMSG_INFO PMessageLevel = 1
const PMSG_WARN PMessageLevel = 2
const PMSG_ERROR PMessageLevel = 3
const PMSG_FATAL PMessageLevel = 4

type PMessageLevel byte

type PMessageStep struct {
	SingleStep
	Hash int32
	Elapsed int32
	Level PMessageLevel
	paramString string
	tempMap map[string]string
}

func NewPMessageStep(startTime int32) *PMessageStep {
	step := new(PMessageStep)
	step.tempMap = make(map[string]string)
	step.StartTime = startTime
	return step
}

func (s *PMessageStep) GetStepType() byte {
	return PARAMETERIZED_MESSAGE
}

func (s *PMessageStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteDecimal32(s.Hash)
	_, err = out.WriteDecimal32(s.Elapsed)
	_, err = out.WriteDecimal32(int32(s.Level))
	_, err = out.WriteString(s.paramString)
	return err
}

func (s *PMessageStep) SetMessage(hash int32, params ...string)  {
	s.Hash = hash
	var b bytes.Buffer
	for _, s := range params {
		b.WriteString(s)
		b.WriteByte(delimETX)
	}
	s.paramString = b.String()
}
