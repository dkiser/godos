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

package cmd

import (
	"fmt"

	"github.com/dkiser/godos/engine"
	"github.com/spf13/cobra"
)

// synfloodCmd vars
var dstIP string
var srcPort, dstPort, concurrency, runtime, pps int

// synfloodCmd helper vars
var goQuitChannel, drainDoneChannel chan bool

// synfloodCmd represents the synflood command
var synfloodCmd = &cobra.Command{
	Use:   "synflood",
	Short: "TCP syn flood attack",
	Long:  `Launch TCP syn flood attack against a destination.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(synfloodCmd)

	synfloodCmd.Flags().StringVarP(&dstIP, "dstIP", "d", "127.0.0.1",
		"dst IP or hostname to attack")
	synfloodCmd.Flags().IntVarP(&dstPort, "dstPort", "p", -1,
		"dst Port or -1 for random")
	synfloodCmd.Flags().IntVarP(&srcPort, "srcPort", "s", -1,
		"src Port or -1 for random")
	synfloodCmd.Flags().IntVarP(&concurrency, "concurrency", "g", 1,
		"number of concurrent go routines to run")
	synfloodCmd.Flags().IntVarP(&runtime, "runtime", "t", 10,
		"number seconds to run")
	synfloodCmd.Flags().IntVarP(&pps, "pps", "x", 10,
		"number of packets to send per second per goroutine")

}

func attack() {
	fmt.Println("running synflood attack")
}

// run the attack
func run() {
	e, _ := engine.NewEngine(runtime, concurrency, 10, attack)
	e.Start()
}
