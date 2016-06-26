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

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

// All engines have the following
type Engine struct {
	// total run time in seconds
	runTime int
	// total go routines to use
	concurrency int
	// time between attackFunc loops
	attackLoopTime int
	// user attack func to call
	attackFunc func()
	// channel used for quitting
	quitChannel chan bool
	// channel used for done complete
	doneChannel chan bool
}

// NewEngine allocates a an Engine type for wrapping attack runs provided by caller
func NewEngine(runtime int, concurrency int, attackLoopTime int, attackFunc func()) (*Engine, error) {
	e := &Engine{
		runTime:        runtime,
		concurrency:    concurrency,
		attackLoopTime: attackLoopTime,
		attackFunc:     attackFunc,
		quitChannel:    make(chan bool),
		doneChannel:    make(chan bool),
	}
	return e, nil
}

// Start launches concurrency/timers for calling attack function
func (e *Engine) Start() {
	fmt.Println("starting engine")

	defer e.close()

	// launch attack runners
	for i := 0; i < e.concurrency; i++ {
		go e.attackRunner()
	}

	// wait and drain attackers
	e.timeAndDrain()
}

// attack goroutine
func (e *Engine) attackRunner() {

	id := uuid.NewRandom()
	tickerLoop := time.NewTicker(time.Second * time.Duration(e.attackLoopTime))
	for {
		// select case to see if we should quit
		select {
		case <-e.quitChannel:
			// we are done, exit loop and goroutine
			fmt.Println("exiting attacker #", id)
			e.doneChannel <- true
			break
		default:
			// launch registered attack function
			fmt.Println("starting attacker #", id)
			e.attackFunc()
			<-tickerLoop.C
		}
	}

}

// general close method
func (e *Engine) close() {
	close(e.doneChannel)
	close(e.quitChannel)
}

// Stop will trigger all goroutines to stop after they are
// done with their respective execution loops
func (e *Engine) Stop() {
	// drain the go routines launched
	for i := 0; i < e.concurrency; i++ {
		select {
		case e.quitChannel <- true:
		default:
		}

		fmt.Println("draining goroutine #", i)
	}

	// wait for all GOs to die
	for i := 0; i < e.concurrency; i++ {
		<-e.doneChannel
	}
}

func (e *Engine) timeAndDrain() {
	// setup ticker to wait for time desired
	fmt.Println("Setting up timer for ", e.runTime, " seconds.")
	ticker := time.NewTicker(time.Second * time.Duration(e.runTime))

	// wait for timer to expire
	<-ticker.C

	e.Stop()
}
