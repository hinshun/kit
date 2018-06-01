package interrupt

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
)

// InterruptHandler helps set up an interrupt handler that can be cleanly shut
// down through the io.Closer interface.
type InterruptHandler struct {
	sig chan os.Signal
	wg  sync.WaitGroup
}

func NewInterruptHandler(ctx context.Context, sigs ...os.Signal) (io.Closer, context.Context) {
	intrh := &InterruptHandler{
		sig: make(chan os.Signal, 1),
	}
	ctx, cancel := context.WithCancel(ctx)

	handlerFunc := func(count int, ih *InterruptHandler) {
		switch count {
		case 1:
			fmt.Println() // Prevent un-terminated ^C character in terminal

			ih.wg.Add(1)
			go func() {
				defer ih.wg.Done()
				cancel()
			}()

		default:
			fmt.Println("Received another interrupt before graceful shutdown, terminating...")
			os.Exit(-1)
		}
	}

	intrh.Handle(handlerFunc, sigs...)
	return intrh, ctx
}

func (ih *InterruptHandler) Close() error {
	close(ih.sig)
	ih.wg.Wait()
	return nil
}

// Handle starts handling the given signals, and will call the handler callback
// function each time a signal is catched. The function is passed the number of
// times the handler has been triggered in total, as well as the handler itself,
// so that the handling logic can use the handler's wait group to ensure clean
// shutdown when Close() is called.
func (ih *InterruptHandler) Handle(handler func(count int, ih *InterruptHandler), sigs ...os.Signal) {
	signal.Notify(ih.sig, sigs...)
	ih.wg.Add(1)
	go func() {
		defer ih.wg.Done()
		count := 0
		for range ih.sig {
			count++
			handler(count, ih)
		}
		signal.Stop(ih.sig)
	}()
}
