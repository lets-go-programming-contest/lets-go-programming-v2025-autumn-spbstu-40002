package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type decoratorReg struct {
	fn       func(ctx context.Context, input chan string, output chan string) error
	inputID  string
	outputID string
}

type multiplexerReg struct {
	fn       func(ctx context.Context, inputs []chan string, output chan string) error
	inputIDs []string
	outputID string
}

type separatorReg struct {
	fn        func(ctx context.Context, input chan string, outputs []chan string) error
	inputID   string
	outputIDs []string
}

type conveyerImpl struct {
	size        int
	mtx         sync.RWMutex
	chans       map[string]chan string
	decorators  []decoratorReg
	multiplexes []multiplexerReg
	separators  []separatorReg
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (cvr *conveyerImpl) ensureChanLocked(ident string) {
	if _, ok := cvr.chans[ident]; !ok {
		cvr.chans[ident] = make(chan string, cvr.size)
	}
}

func (cvr *conveyerImpl) getChan(ident string) chan string {
	cvr.mtx.Lock()
	defer cvr.mtx.Unlock()
	if _, ok := cvr.chans[ident]; !ok {
		cvr.chans[ident] = make(chan string, cvr.size)
	}
	return cvr.chans[ident]
}

func (cvr *conveyerImpl) getChanIfExists(ident string) chan string {
	cvr.mtx.RLock()
	defer cvr.mtx.RUnlock()
	return cvr.chans[ident]
}

func (cvr *conveyerImpl) closeAll() {
	cvr.mtx.Lock()
	defer cvr.mtx.Unlock()
	for ident, chn := range cvr.chans {
		if chn != nil {
			close(chn)
			cvr.chans[ident] = nil
		}
	}
}

func (cvr *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	cvr.mtx.Lock()
	cvr.ensureChanLocked(input)
	cvr.ensureChanLocked(output)
	cvr.decorators = append(cvr.decorators, decoratorReg{
		fn:       fn,
		inputID:  input,
		outputID: output,
	})
	cvr.mtx.Unlock()
}

func (cvr *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	cvr.mtx.Lock()
	for _, id := range inputs {
		cvr.ensureChanLocked(id)
	}
	cvr.ensureChanLocked(output)
	cvr.multiplexes = append(cvr.multiplexes, multiplexerReg{
		fn:       fn,
		inputIDs: inputs,
		outputID: output,
	})
	cvr.mtx.Unlock()
}

func (cvr *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	cvr.mtx.Lock()
	cvr.ensureChanLocked(input)
	for _, id := range outputs {
		cvr.ensureChanLocked(id)
	}
	cvr.separators = append(cvr.separators, separatorReg{
		fn:        fn,
		inputID:   input,
		outputIDs: outputs,
	})
	cvr.mtx.Unlock()
}

func (cvr *conveyerImpl) Run(ctx context.Context) error {
	ctxRun, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	var wgroup sync.WaitGroup
	errchan := make(chan error, 1)

	for _, reg := range cvr.decorators {
		inCh := cvr.getChan(reg.inputID)
		outCh := cvr.getChan(reg.outputID)
		wgroup.Add(1)
		go func(fn func(context.Context, chan string, chan string) error, rin chan string, rout chan string) {
			defer wgroup.Done()
			if err := fn(ctxRun, rin, rout); err != nil {
				select {
				case errchan <- err:
				default:
				}
				cancelRun()
			}
		}(reg.fn, inCh, outCh)
	}

	for _, reg := range cvr.multiplexes {
		var ins []chan string
		for _, id := range reg.inputIDs {
			ins = append(ins, cvr.getChan(id))
		}
		outCh := cvr.getChan(reg.outputID)
		wgroup.Add(1)
		go func(fn func(context.Context, []chan string, chan string) error, ins []chan string, rout chan string) {
			defer wgroup.Done()
			if err := fn(ctxRun, ins, rout); err != nil {
				select {
				case errchan <- err:
				default:
				}
				cancelRun()
			}
		}(reg.fn, ins, outCh)
	}

	for _, reg := range cvr.separators {
		inCh := cvr.getChan(reg.inputID)
		var outs []chan string
		for _, id := range reg.outputIDs {
			outs = append(outs, cvr.getChan(id))
		}
		wgroup.Add(1)
		go func(fn func(context.Context, chan string, []chan string) error, rin chan string, routs []chan string) {
			defer wgroup.Done()
			if err := fn(ctxRun, rin, routs); err != nil {
				select {
				case errchan <- err:
				default:
				}
				cancelRun()
			}
		}(reg.fn, inCh, outs)
	}

	done := make(chan struct{})
	go func() {
		wgroup.Wait()
		close(done)
	}()

	select {
	case err := <-errchan:
		cancelRun()
		wgroup.Wait()
		cvr.closeAll()
		return err
	case <-ctxRun.Done():
		wgroup.Wait()
		cvr.closeAll()
		return ctxRun.Err()
	case <-done:
		cvr.closeAll()
		return nil
	}
}

func (cvr *conveyerImpl) Send(input string, data string) error {
	chn := cvr.getChanIfExists(input)
	if chn == nil {
		return ErrChanNotFound
	}
	defer func() { _ = recover() }()
	chn <- data
	return nil
}

func (cvr *conveyerImpl) Recv(output string) (string, error) {
	chn := cvr.getChanIfExists(output)
	if chn == nil {
		return "", ErrChanNotFound
	}
	val, ok := <-chn
	if !ok {
		return "undefined", nil
	}
	return val, nil
}
