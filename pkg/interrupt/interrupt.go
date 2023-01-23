package interrupt

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// InterruptHandler helps set up an interrupt handler that can be cleanly shut
// down through the io.Closer interface.
type InterruptHandler struct {
	sig chan os.Signal
	wg  sync.WaitGroup
}

type handlerFunc func(ih *InterruptHandler, sig os.Signal)

// NewInterruptHandler returns a new interrupt handler that will invoke cancel
// if any of the signals provided are received.
func NewInterruptHandler(cancel context.CancelFunc, sigs ...os.Signal) io.Closer {
	intrh := &InterruptHandler{
		sig: make(chan os.Signal, 1),
	}

	count := 0
	handlerFunc := func(ih *InterruptHandler, sig os.Signal) {
		count++
		switch count {
		case 1:
			// Prevent un-terminated ^C character in terminal.
			fmt.Println()

			ih.wg.Add(1)
			go func() {
				defer ih.wg.Done()
				cancel()
			}()

		default:
			fmt.Println("Received another interrupt before graceful shutdown, terminating...")

			syscallSig, ok := sig.(syscall.Signal)
			if !ok {
				os.Exit(-1)
			}

			// Fatal errors exit with 128+n, where "n" is the syscall.Signal code.
			os.Exit(128 + int(syscallSig))
		}
	}

	intrh.Handle(handlerFunc, sigs...)
	return intrh
}

// Close closes its signal receiver and waits for its handlers to exit cleanly.
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
func (ih *InterruptHandler) Handle(handler handlerFunc, sigs ...os.Signal) {
	signal.Notify(ih.sig, sigs...)
	ih.wg.Add(1)
	go func() {
		defer ih.wg.Done()
		for sig := range ih.sig {
			handler(ih, sig)
		}
		signal.Stop(ih.sig)
	}()
}
