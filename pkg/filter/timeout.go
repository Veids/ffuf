package filter

import (
	"sync"
)

type TimeoutManager struct {
	Count uint64
	Mutex sync.Mutex
}

type PerDomainTimeout struct {
	Timeouts  map[string]*TimeoutManager
	Threshold uint64
	Mutex     sync.Mutex
}

func (p *PerDomainTimeout) InitHost(host string) *TimeoutManager {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if p.Timeouts[host] == nil {
		var tm = NewTimeoutManager()
		p.Timeouts[host] = tm
		return tm
	}

	return p.Timeouts[host]
}

func (p *PerDomainTimeout) IncreaseTimeoutCount(host string) {
	if p.Timeouts[host] != nil {
		p.Timeouts[host].Increase()
	} else {
		p.InitHost(host).Increase()
	}
}

func (p *PerDomainTimeout) ResetTimeoutCount(host string) {
	if p.Timeouts[host] != nil {
		p.Timeouts[host].Reset()
	}
}

func (p *PerDomainTimeout) IsTimedOut(host string) bool {
	if p.Timeouts[host] != nil {
		return p.Timeouts[host].Count >= p.Threshold
	}

	return false
}

func (p *PerDomainTimeout) IsEnabled() bool {
	return p.Threshold > 0
}

func NewPerDomainTimeout(threshold uint64) *PerDomainTimeout {
	return &PerDomainTimeout{
		Timeouts:  make(map[string]*TimeoutManager),
		Threshold: threshold,
	}
}

func NewTimeoutManager() *TimeoutManager {
	return &TimeoutManager{
		Count: 1,
	}
}

func (t *TimeoutManager) Increase() {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.Count += 1
}

func (t *TimeoutManager) Reset() {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.Count = 0
}
