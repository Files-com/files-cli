package lib

import (
	"context"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
)

var PagerEnvVariables = []string{"PAGER"}
var PagerCommands = []string{"less", "more"}

type Pager struct {
	DefaultOut io.Writer
	io.Writer
	UsePager  bool
	cmd       *exec.Cmd
	runOnce   sync.Once
	startOnce sync.Once
	wait      chan bool
}

func (p Pager) Init(previewInput interface{}, defaultOut io.Writer) (*Pager, error) {
	var err error
	p.wait = make(chan bool)
	path, _, err := pagerExecPath()
	if !singleResult(previewInput) && isTerminal(defaultOut) && p.UsePager && err == nil {
		var args []string
		p.cmd = exec.Command(path, args...)
		p.Writer, err = p.cmd.StdinPipe()
		if err != nil {
			return &p, err
		}
		p.cmd.Stdout = defaultOut
	} else {
		p.Writer = defaultOut
	}

	return &p, nil
}

func (p *Pager) Start(cancel context.CancelFunc) {
	p.startOnce.Do(func() {
		if p.cmd != nil {
			p.cmd.Start()
		}
		go func() {
			if p.cmd != nil {
				p.cmd.Wait()
				cancel()
				p.wait <- true
			}
		}()
	})
}

func (p *Pager) Canceled(ctx context.Context) bool {
	select {
	case <-p.wait:
		// pager exited early
		return true
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (p *Pager) Wait() {
	closer, ok := p.Writer.(io.WriteCloser)
	if ok {
		closer.Close()
	}
	if p.cmd != nil {
		<-p.wait
	}
}

func pagerExecPath() (path string, args []string, err error) {
	for _, testVar := range PagerEnvVariables {
		path = os.Getenv(testVar)
		if path != "" {
			// BUG: does not handle multiple spaces, e.g.: "less -s  -R"
			args = strings.Split(path, " ")
			return args[0], args[1:], nil
		}
	}

	// This default only gets used if PagerCommands is empty.
	err = exec.ErrNotFound
	for _, testPath := range PagerCommands {
		path, err = exec.LookPath(testPath)
		if err == nil {
			return path, nil, nil
		}
	}
	return "", nil, err
}

func singleResult(result interface{}) bool {
	switch result.(type) {
	case []interface{}:
		sliceResults, _ := result.([]interface{})
		return len(sliceResults) == 1 || len(sliceResults) == 0
	case map[interface{}]interface{}:
		return true
	case Iter:
		return false
	case nil:
		return true
	default:
		switch reflect.ValueOf(result).Kind().String() {
		case "struct":
			return true
		case "invalid":
			return true
		}
		return false
	}
}
