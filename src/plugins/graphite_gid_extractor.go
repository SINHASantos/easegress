package plugins

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"common"
	"pipelines"
	"task"
)

type graphiteGidExtractorConfig struct {
	CommonConfig

	GidKey  string `json:"gid_key"`
	DataKey string `json:"data_key"`
}

func GraphiteGidExtractorConfigConstructor() Config {
	return &graphiteGidExtractorConfig{}
}

func (c *graphiteGidExtractorConfig) Prepare() error {
	err := c.CommonConfig.Prepare()
	if err != nil {
		return err
	}

	ts := strings.TrimSpace
	c.GidKey, c.DataKey = ts(c.GidKey), ts(c.DataKey)

	if len(c.GidKey) == 0 {
		return fmt.Errorf("invalid gid key")
	}

	if len(c.DataKey) == 0 {
		return fmt.Errorf("invalid data key")
	}

	return nil
}

type graphiteGidExtractor struct {
	conf *graphiteGidExtractorConfig
}

func GraphiteGidExtractorConstructor(conf Config) (Plugin, error) {
	c, ok := conf.(*graphiteGidExtractorConfig)
	if !ok {
		return nil, fmt.Errorf("config type want *graphiteGidExtractorConfig got %T", conf)
	}

	return &graphiteGidExtractor{
		conf: c,
	}, nil
}

func (e *graphiteGidExtractor) Prepare(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (e *graphiteGidExtractor) extract(t task.Task) (error, task.TaskResultCode, task.Task) {
	dataValue := t.Value(e.conf.DataKey)
	data, ok := dataValue.([]byte)
	if !ok {
		return fmt.Errorf("input %s got wrong value: %#v", e.conf.DataKey, dataValue),
			task.ResultMissingInput, t
	}

	s := bufio.NewScanner(bytes.NewReader(data))
	if !s.Scan() {
		return fmt.Errorf("unexpected EOF"), task.ResultBadInput, t
	}

	fields := common.GraphiteSplit(s.Text(), ".", "#")
	if len(fields) != 4 {
		return fmt.Errorf("graphite data want 4 fields('#'-splitted) got %v",
			len(fields)), task.ResultBadInput, t
	}

	// system application instance hostipv4 hostname
	gid := strings.Join([]string{fields[0], "", fields[1], fields[2], fields[3]}, "")

	var err error
	t, err = task.WithValue(t, e.conf.GidKey, gid)
	if err != nil {
		return err, task.ResultInternalServerError, t
	}

	return nil, t.ResultCode(), t
}

func (e *graphiteGidExtractor) Run(ctx pipelines.PipelineContext, t task.Task) (task.Task, error) {
	err, resultCode, t := e.extract(t)
	if err != nil {
		t.SetError(err, resultCode)
	}
	return t, nil
}

func (e *graphiteGidExtractor) Name() string {
	return e.conf.PluginName()
}

func (e *graphiteGidExtractor) Close() {
	// Nothing to do.
}