package netio

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
)

var aconf = conf.GetInstance()

type ProfileCollector struct {
	tctx *TraceContext
	steps []netdata.Step
	pos int
	doingDumpStepJob bool
	currentLevel int32
	parentLevel int32
}

func NewProfileCollector(tctx *TraceContext) *ProfileCollector {
	collector := new(ProfileCollector)
	collector.tctx = tctx
	collector.steps = make([]netdata.Step, aconf.ProfileStepMaxKeepInMemoryCount)
	return collector
}

func (c *ProfileCollector) Push(ss netdata.Step) {
	c.checkDumpStep()
	ss.SetIndex(c.currentLevel)
	ss.SetParent(c.parentLevel)
	c.parentLevel = c.currentLevel
	c.currentLevel++
}

func (c *ProfileCollector) Pop(ss netdata.Step) {
	c.checkDumpStep()
	c.parentLevel = ss.GetParent()
	c.Process(ss)
}

func (c *ProfileCollector) Add(ss netdata.Step) {
	c.checkDumpStep()
	ss.SetIndex(c.currentLevel)
	ss.SetParent(c.parentLevel)
	c.currentLevel++
	c.Process(ss)
}

func (c *ProfileCollector) Process(ss netdata.Step) {
	c.checkDumpStep()
	c.steps[c.pos] = ss
	c.pos++
	if c.pos >= len(c.steps) {
		var o = c.steps
		c.steps = make([]netdata.Step, aconf.ProfileStepMaxKeepInMemoryCount)
		c.pos = 0;
		SendProfile(o, c.tctx)
	}
}

func (c *ProfileCollector) Close(ok bool) {
	c.checkDumpStep()
	if c.pos > 0 && ok {
		SendProfile(c.steps[0 : c.pos], c.tctx)
	}
}

func (*ProfileCollector) checkDumpStep() {
	//TODO
	/*
		if(doingDumpStepJob) {
			return;
		}

		DumpStep dumpStep;
		doingDumpStepJob = true;
		while(true) {
			dumpStep = context.temporaryDumpSteps.poll();
			if(dumpStep == null) {
				break;
			}
			add(dumpStep);
		}
		doingDumpStepJob = false;
	 */
}
