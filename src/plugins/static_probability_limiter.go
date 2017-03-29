package plugins

import (
	"fmt"
	"math/rand"
	"time"

	"logger"
	"pipelines"
	"task"
)

type staticProbabilityLimiterConfig struct {
	CommonConfig
	PassPr float32 `json:"pass_pr"`
}

func StaticProbabilityLimiterConfigConstructor() Config {
	return new(staticProbabilityLimiterConfig)
}

func (c *staticProbabilityLimiterConfig) Prepare() error {
	err := c.CommonConfig.Prepare()
	if err != nil {
		return err
	}

	if c.PassPr < 0 || c.PassPr > 1 {
		return fmt.Errorf("invalid passing probability %f", c.PassPr)
	}

	if c.PassPr == 0 {
		logger.Warnf("[ZERO passing probablity has been applied, no request could be processed!]")
	}

	if c.PassPr == 1 {
		logger.Warnf("[1.0 passing probablity has been applied, no request could be limited!]")
	}

	return nil
}

type staticProbabilityLimiter struct {
	conf *staticProbabilityLimiterConfig
}

func StaticProbabilityLimiterConstructor(conf Config) (Plugin, error) {
	c, ok := conf.(*staticProbabilityLimiterConfig)
	if !ok {
		return nil, fmt.Errorf("config type want *staticProbabilityLimiterConfig got %T", conf)
	}

	return &staticProbabilityLimiter{
		conf: c,
	}, nil
}

func (l *staticProbabilityLimiter) Prepare(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (l *staticProbabilityLimiter) Run(ctx pipelines.PipelineContext, t task.Task) (task.Task, error) {
	if rand.Float32() < 1.0-l.conf.PassPr {
		t.SetError(fmt.Errorf("service is unavaialbe caused by probability limit"), task.ResultFlowControl)
	}
	return t, nil
}

func (l *staticProbabilityLimiter) Name() string {
	return l.conf.PluginName()
}

func (l *staticProbabilityLimiter) Close() {
	// Nothing to do.
}

////

func init() {
	rand.Seed(time.Now().UnixNano())
}