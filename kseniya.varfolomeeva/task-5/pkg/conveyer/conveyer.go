package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"golang.org/x/sync/errgroup"
)

var ErrNoSuchChannel = errors.New("no such channel")

const EmptyValue = "undefined"

type ConveyerAPI interface {
	RegisterDataHandler(handler func(context.Context, chan string, chan string) error, inputName, outputName string)
	RegisterCombiner(handler func(context.Context, []chan string, chan string) error, inputNames []string, outputName string)
	RegisterSplitter(handler func(context.Context, chan string, []chan string) error, inputName string, outputNames []string)
	RunPipeline(ctx context.Context) error
	SendData(input string, data string) error
	ReceiveData(output string) (string, error)
}

type pipelineTask func(context.Context) error

type Pipeline struct {
	queueSize  int
	dataPipes  map[string]chan string
	tasks      []pipelineTask
	pipeLock   sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		queueSize: size,
		dataPipes: make(map[string]chan string),
		tasks:     make([]pipelineTask, 0),
		pipeLock:  sync.RWMutex{},
	}
}

func (p *Pipeline) getPipe(name string) chan string {
	p.pipeLock.Lock()
	defer p.pipeLock.Unlock()

	if pipe, exists := p.dataPipes[name]; exists {
		return pipe
	}

	newPipe := make(chan string, p.queueSize)
	p.dataPipes[name] = newPipe
	return newPipe
}

func (p *Pipeline) findPipe(name string) (chan string, error) {
	p.pipeLock.RLock()
	defer p.pipeLock.RUnlock()

	if pipe, exists := p.dataPipes[name]; exists {
		return pipe, nil
	}

	return nil, ErrNoSuchChannel
}

func (p *Pipeline) RegisterDataHandler(
	handler func(context.Context, chan string, chan string) error,
	input, output string,
) {
	inPipe := p.getPipe(input)
	outPipe := p.getPipe(output)

	task := func(ctx context.Context) error {
		return handler(ctx, inPipe, outPipe)
	}

	p.addTask(task)
}

func (p *Pipeline) RegisterCombiner(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inPipes := make([]chan string, len(inputs))
	for i, name := range inputs {
		inPipes[i] = p.getPipe(name)
	}

	outPipe := p.getPipe(output)

	task := func(ctx context.Context) error {
		return handler(ctx, inPipes, outPipe)
	}

	p.addTask(task)
}

func (p *Pipeline) RegisterSplitter(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inPipe := p.getPipe(input)

	outPipes := make([]chan string, len(outputs))
	for i, name := range outputs {
		outPipes[i] = p.getPipe(name)
	}

	task := func(ctx context.Context) error {
		return handler(ctx, inPipe, outPipes)
	}

	p.addTask(task)
}

func (p *Pipeline) SendData(input string, data string) error {
	pipe, err := p.findPipe(input)
	if err != nil {
		return err
	}

	select {
	case pipe <- data:
		return nil
	default:
		return fmt.Errorf("pipe %s is full", input)
	}
}

func (p *Pipeline) ReceiveData(output string) (string, error) {
	pipe, err := p.findPipe(output)
	if err != nil {
		return "", err
	}

	value, ok := <-pipe
	if !ok {
		return EmptyValue, nil
	}

	return value, nil
}

func (p *Pipeline) closeAllPipes() {
	p.pipeLock.Lock()
	defer p.pipeLock.Unlock()

	for _, pipe := range p.dataPipes {
		close(pipe)
	}
}

func (p *Pipeline) RunPipeline(ctx context.Context) error {
	p.pipeLock.RLock()
	defer p.pipeLock.RUnlock()

	group, ctx := errgroup.WithContext(ctx)

	for _, task := range p.tasks {
		currentTask := task
		group.Go(func() error {
			return currentTask(ctx)
		})
	}

	err := group.Wait()
	p.closeAllPipes()

	if err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

func (p *Pipeline) addTask(task pipelineTask) {
	p.pipeLock.Lock()
	defer p.pipeLock.Unlock()
	p.tasks = append(p.tasks, task)
}
