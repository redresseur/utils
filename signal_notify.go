package utils

import (
	"context"
	"os"
	"os/signal"
)

func WatchSignal(watchList map[os.Signal]func(), ctx context.Context) {
	ch := make(chan os.Signal, 1)
	ss := []os.Signal{}
	for k := range watchList {
		ss = append(ss, k)
	}

	signal.Notify(ch, ss...)
	for {
		select {
		case s, ok := <-ch:
			if !ok {
				return
			}
			watchList[s]()
		case <-ctx.Done():
			return
		}

	}
}
