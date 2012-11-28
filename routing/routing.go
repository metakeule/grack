package grack_routing

import (
	. "github.com/metakeule/grack"
	h "github.com/metakeule/grack/http"
)

type RoutingRacker interface {
	h.HTTPRacker // implies Racker, therefor Racker must not be included to avoid duplications
	Routable_() *Routable
}

type dispatcher struct {
	Rack Racker
}

func (d *dispatcher) Rack_() *Rack {
	return d.Rack.Rack_()
}

func (d *dispatcher) HTTPRack_() *h.HTTPRack {
	return d.Rack.(h.HTTPRacker).HTTPRack_()
}

func (d *dispatcher) NewContext() h.HTTPContexter {
	return d.Rack.(h.HTTPRacker).NewContext()
}

// a Routable Rack behaves to the Context like normal Racker but has a Router middleware
// we have the dispatcher only to be able to inherit from different places
type Routable struct {
	*dispatcher
	Position uint
}

func (r *Routable) Routable_() *Routable {
	return r
}

func (r *Routable) Push(i interface{}) {
	r.dispatcher.Rack.Push(i)
}

func (r *Routable) SetResponder(i interface{}) {
	r.dispatcher.Rack.SetResponder(i)
}

//Stack

func (r *Routable) Router(mw interface{}) {
	ra := r
	r.Position = Len(ra)
	Push(ra, mw)
}

func (r *Routable) Clone() Racker {
	ra := r.dispatcher.Rack.Clone()
	d := &dispatcher{Rack: ra}
	return &Routable{dispatcher: d, Position: r.Position}
}

func NewRoutingRack(r Racker) RoutingRacker {
	d := &dispatcher{Rack: r}
	return &Routable{dispatcher: d}
}

func Router(r RoutingRacker, mw interface{}) {
	r.Routable_().Router(mw)
}

func JumpToRouter(c Contexter) {
	r := c.Ctx().Rack.(RoutingRacker).Routable_()
	pos := r.Position
	GoTo(c, pos)
}
