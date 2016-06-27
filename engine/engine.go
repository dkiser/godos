// Copyright © 2016 Domingo Kiser
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package engine

import "time"

// Engine provides a construct around launching an attack
// function every LoopTicker time for a total of RunTicker time
type Engine struct {
	// total run ticker
	RunTicker *time.Ticker
	// attak loop ticker
	LoopTicker *time.Ticker
	// total go routines to use
	Concurrency int
	// user attack func to call
	AttackFunc func()
	// time when engine started
	TimeStart time.Time
	// time when engine stopped
	TimeStop time.Time
	// duration of engine run
	TimeDuration time.Duration
}

// NewEngine allocates a an Engine type for wrapping attack runs provided by caller
func NewEngine(runtime time.Duration, loopTime time.Duration, concurrency int, attackFunc func()) (*Engine, error) {
	e := &Engine{
		RunTicker:   time.NewTicker(runtime),
		LoopTicker:  time.NewTicker(loopTime),
		Concurrency: concurrency,
		AttackFunc:  attackFunc,
	}
	// Only log the warning severity or above.
	return e, nil
}

//
func (e *Engine) goAttack() {
	//id := uuid.NewRandom()
	// launch attack runners
	for i := 0; i < e.Concurrency; i++ {
		//fmt.Println("starting attacker #", id)
		go e.AttackFunc()
	}
}

// Start runs the attack function loop
func (e *Engine) Start() {
	e.TimeStart = time.Now()
	// launch first iteration
	e.goAttack()

	// main attack func loop
	for {
		select {
		case <-e.RunTicker.C:
			//fmt.Println("Timer expired")
			e.Stop()
			return
		case <-e.LoopTicker.C:
			//fmt.Println("Ticker ticked")
			// launch attack runners
			e.goAttack()
		}
	}
}

// Close will free up resources of the Engine
func (e *Engine) Close() {
	e.RunTicker.Stop()
	e.LoopTicker.Stop()
}

// Stop will trigger launch of attack funcs to stop
func (e *Engine) Stop() {
	e.TimeStop = time.Now()
	e.TimeDuration = e.TimeStop.Sub(e.TimeStart)
}
