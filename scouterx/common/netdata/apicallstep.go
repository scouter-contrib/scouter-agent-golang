package netdata

type ApiCallStep struct {
	SingleStep
	Txid int64
	Hash int32
	Elapsed int32
	CpuTime int32
	Error int32
	Opt uint8
	Address string
}

func NewApiCallStep() *ApiCallStep {
	step := new(ApiCallStep)
	return step
}

func (s *ApiCallStep) GetStepType() byte {
	return APICALL
}

func (s *ApiCallStep) Write(out *DataOutputX) error {
	err := s.SingleStep.Write(out)
	if err != nil {
		return err
	}
	_, err = out.WriteDecimal(s.Txid)
	_, err = out.WriteDecimal32(s.Hash)
	_, err = out.WriteDecimal32(s.Elapsed)
	_, err = out.WriteDecimal32(s.CpuTime)
	_, err = out.WriteDecimal32(s.Error)
	_, err = out.WriteUInt8(s.Opt)
	if s.Opt == 1 {
		_, err = out.WriteString(s.Address)
	}

	return err
}
