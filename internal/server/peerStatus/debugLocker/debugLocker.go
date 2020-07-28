package debugLocker

import (
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type MyLocker interface {
	Lock(string) func()
	Unlock(string)
}

type normalLocker struct {
	mu sync.Mutex
}

type debugLocker struct {
	mu        sync.Mutex
	smu       sync.Mutex
	reference []string
}

func NewMutex() MyLocker {
	if tools.IsDevelopmennt() {
		return &debugLocker{
			reference: make([]string, 0, 20),
		}
	} else {
		return &normalLocker{}
	}
}

func (l *normalLocker) Lock(_ string) func() {
	l.mu.Lock()
	return func() {
		l.Unlock("")
	}
}

func (l *debugLocker) Lock(who string) func() {
	l.smu.Lock()
	l.reference = append(l.reference, who)
	if len(l.reference) > 1 {
		tools.Debug("Lock wait queue[%d]: %v", len(l.reference), l.reference)
	}
	l.smu.Unlock()

	l.mu.Lock()

	return func() {
		l.Unlock(who)
	}
}

func (l *normalLocker) Unlock(_ string) {
	l.mu.Unlock()
}

func (l *debugLocker) Unlock(who string) {
	// tools.Error("Unlock - %s", who)

	l.smu.Lock()
	if len(l.reference) == 0 || l.reference[0] != who {
		tools.Die("Not time for unlock: %s (queue: %v)", who, l.reference)
	}
	l.reference = l.reference[1:]
	l.smu.Unlock()

	l.mu.Unlock()
}
