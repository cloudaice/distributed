package main

import (
	"time"
)

var (
	global int
)

type CirTh struct {
	STS int64 // TimeStamp start want to do this Job.
	Do  func() error
}

func NewCirTh(ts int64) *CirTh {
	return &CirTh{
		STS: ts,
		Do: func() error {
			global++
			return nil
		},
	}
}

type State int

const (
	WANTIN State = iota
	IMIN
	NORMAL
)

type Action int

const (
	NEEDIN Action = iota
	ACKIN
	DOIN
)

type Event struct {
	Sender   *P
	Receiver *P
	TS       int64
	Name     int
}

//P mwans a processor, one processor just have a groutine.
type P struct {
	In chan Event
	ID int
	TS int64
	ST int
	WS []Event
	G  map[int]*P
	CT *CirTh
}

func NewP(id int) {
	p := &P{
		In: make(chan Event, 1),
		ID: id,
		TS: time.Now().Unix(),
		ST: NORMAL,
		G:  make(map[int]*P),
	}
	p.addP(p)
}

// addP add a processor in group,
// it always called by reciever self.
func (p *P) addP(p *P) {
	p.G[p.ID] = p
}

// delP delete a processoer with id in group,
// it always called by reciever self.
func (p *P) delP(id int) {
	delete(p.G, id)
}

func (p *P) Run() {
	go func() {
		for {
			e := <-p.In
			switch e.Name {
			case NEEDIN:
				switch p.ST {
				case WANTIN:
					if p.TS <= e.TS {
						p.TS = e.TS + 1
					}
					msg := Event{
						Sender:   p,
						Receiver: e.Sender,
						TS:       p.TS,
						Name:     ACKIN,
					}
					if p.CT.STS > e.TS {
						Send(msg)
					} else {
						p.WS = append(p.WS, msg)
					}

				case IMIN:
					msg := Event{
						Sender:   p,
						Receiver: e.Sender,
						TS:       p.TS,
						Name:     ACKIN,
					}

				}
			case NORMAL:
				msg := Event{
					Sender:   p,
					Receiver: e.Sender,
					TS:       p.TS,
					Name:     ACKIN,
				}
				Send(msg)
			}
		}
	}()
}

// Send write event into itself.
func (p *P) Send(e Event) {
	e.TS = P.TS
	p.In <- e
}

// Send send event to p.
func Send(e Event) {
	e.Receiver.Send(e)
}
