package subpub

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

// check correctness of MySubscription
var _ Subscription = (*MySubscription)(nil)

type MySubscription struct {
	sp        *MySubPub // origin subpub
	publisher string
	handler   MessageHandler
}

func (s *MySubscription) Unsubscribe() {
	value := s.sp.subs[s.publisher]

	for i, ms := range value {
		if ms == s {
			value = slices.Delete(value, i, i+1)
		}
	}

	s.sp.subs[s.publisher] = value
}

// check correctness of MySubPub
var _ SubPub = (*MySubPub)(nil)

type MySubPub struct {
	subs     map[string][]*MySubscription // contains list of subscriptions for publishier
	meta     map[*MySubscription][]any    // contains list of sent messages in FIFO order
	queue    chan *MySubscription         // chan of handling subs
	token    sync.Mutex                   //  guard subPub
	finished chan struct{}                // observer of work
}

func NewSubPub() SubPub {
	sp := &MySubPub{subs: make(map[string][]*MySubscription), queue: make(chan *MySubscription),
		finished: make(chan struct{}, 1), meta: make(map[*MySubscription][]any)}

	// run handler
	go func() {
		for {
			sub, exist := <-sp.queue
			if !exist {
				return
			}

			// handle one of subscribers in separete goroutine
			go func() {
				for {
					sp.token.Lock()
					slice := sp.meta[sub]
					sp.token.Unlock()
					for _, msg := range slice {
						sub.handler(msg)
					}
					sp.token.Lock()
					slice2 := sp.meta[sub]
					slice2 = slice2[len(slice):]
					sp.meta[sub] = slice2
					sp.token.Unlock()
					if len(slice2) == 0 {
						break
					}
				}
			}()
		}
	}()

	return sp
}

func (sp *MySubPub) Subscribe(subj string, cb MessageHandler) (Subscription, error) {
	// check if finished
	select {
	case <-sp.finished:
		sp.finished <- struct{}{}
		return nil, fmt.Errorf("already finished")
	default:
	}

	val := sp.subs[subj]
	subs := &MySubscription{sp: sp, publisher: subj, handler: cb}
	sp.subs[subj] = append(val, subs)
	sp.meta[subs] = nil

	return subs, nil
}

func (sp *MySubPub) Publish(subj string, msg any) error {
	// check if finished
	select {
	case <-sp.finished:
		sp.finished <- struct{}{}
		return fmt.Errorf("already finished")
	default:
	}

	val := sp.subs[subj]

	for _, s := range val {
		sp.token.Lock()
		tmp := sp.meta[s]
		if len(tmp) == 0 {
			sp.queue <- s
		}
		tmp = append(tmp, msg)
		sp.meta[s] = tmp
		sp.token.Unlock()
	}

	return nil
}

func (sp *MySubPub) Close(ctx context.Context) error {
	select {
	case <-sp.finished:
		sp.finished <- struct{}{}
		return fmt.Errorf("already finished")
	default:
	}

	defer func() {
		sp.finished <- struct{}{}
		close(sp.queue)
	}()

	if ctx.Done() == nil {
		return nil
	}

	<-ctx.Done()
	return ctx.Err()
}
