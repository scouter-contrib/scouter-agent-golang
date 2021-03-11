package netdata

const METHOD byte = 1
const METHOD2 byte = 10;
const SPAN byte = 51;
const SQL byte = 2;
const SQL2 byte = 8;
const SQL3 byte = 16;
const MESSAGE byte = 3;
const SOCKET byte = 5;
const APICALL byte = 6;
const APICALL2 byte = 15;
const SPANCALL byte = 52;
const THREAD_SUBMIT byte = 7;
const HASHED_MESSAGE byte = 9;
const PARAMETERIZED_MESSAGE byte = 17;
const DUMP byte = 12;
const DISPATCH byte = 13;
const THREAD_CALL_POSSIBLE byte = 14;

const METHOD_SUM byte = 11;
const SQL_SUM byte = 21;
const MESSAGE_SUM byte = 31;
const SOCKET_SUM byte = 42;
const APICALL_SUM byte = 43;
const CONTROL byte = 99;

type Step interface {
	GetOrder() int32
	GetStepType() byte
	Write(out *DataOutputX) error
	Read(in *DataInputX) (Step, error)

	SetIndex(index int32)
	SetParent(parent int32)
	GetParent() int32
}

type SingleStep struct {
	Parent int32
	Index int32
	StartTime int32
}

func (ss *SingleStep) GetOrder() int32 {
	return ss.Index
}

func (ss *SingleStep) GetStepType() byte {
	//TODO SingleStep.GetStepType()
	return 0
}

func (ss *SingleStep) SetIndex(index int32) {
	ss.Index = index
}

func (ss *SingleStep) SetParent(parent int32) {
	ss.Parent = parent
}

func (ss *SingleStep) GetParent() int32 {
	return ss.Parent
}

func (ss *SingleStep) Write(out *DataOutputX) error {
	var err error
	_, err = out.WriteDecimal32(ss.Parent)
	_, err = out.WriteDecimal32(ss.Index)
	_, err = out.WriteDecimal32(ss.StartTime)
	_, err = out.WriteDecimal32(0)

	return err
}

func (ss *SingleStep) Read(in *DataInputX) (Step, error) {
	//TODO SingleStep.Read()
	return ss, nil
}

func StepsToBytes(steps []Step) []byte {
	if steps == nil {
		return nil
	}
	dout := NewDataOutputX(nil)
	for _, step := range steps {
		if step != nil {
			dout.WriteStep(step)
		}
	}
	return dout.Bytes()
}


