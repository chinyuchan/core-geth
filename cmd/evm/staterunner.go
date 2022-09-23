// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"

	"github.com/urfave/cli/v2"
)

var stateTestCommand = &cli.Command{
	Action:    stateTestCmd,
	Name:      "statetest",
	Usage:     "executes the given state tests",
	ArgsUsage: "<file>",
	Flags:     []cli.Flag{stateTestEVMCEWASMFlag, stateTestForkFlag},
	Category:  flags.DevCategory,
}

var stateTestEVMCEWASMFlag = &cli.StringFlag{
	Name:     "evmc.ewasm",
	Usage:    "EVMC EWASM configuration",
	Category: flags.DevCategory,
}

var stateTestForkFlag = &cli.StringFlag{
	Name:     "fork",
	Usage:    "Fork to use for the test",
	Category: flags.DevCategory,
}

// StatetestResult contains the execution status after running a state test, any
// error that might have occurred and a dump of the final state if requested.
type StatetestResult struct {
	Name  string      `json:"name"`
	Pass  bool        `json:"pass"`
	Fork  string      `json:"fork"`
	Error string      `json:"error,omitempty"`
	State *state.Dump `json:"state,omitempty"`
}

func stateTestCmd(ctx *cli.Context) error {
	if len(ctx.Args().First()) == 0 {
		return errors.New("path-to-test argument required")
	}
	// Configure the go-ethereum logger
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.Int(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)

	if s := ctx.String(stateTestEVMCEWASMFlag.Name); s != "" {
		log.Info("Running tests with %s=%s", "evmc.ewasm", s)
		vm.InitEVMCEwasm(s)
	}

	// Configure the EVM logger
	config := &logger.Config{
		EnableMemory:     !ctx.Bool(DisableMemoryFlag.Name),
		DisableStack:     ctx.Bool(DisableStackFlag.Name),
		DisableStorage:   ctx.Bool(DisableStorageFlag.Name),
		EnableReturnData: !ctx.Bool(DisableReturnDataFlag.Name),
	}
	var (
		tracer   vm.EVMLogger
		debugger *logger.StructLogger
	)
	switch {
	case ctx.Bool(MachineFlag.Name):
		tracer = logger.NewJSONLogger(config, os.Stderr)

	case ctx.Bool(DebugFlag.Name):
		debugger = logger.NewStructLogger(config)
		tracer = debugger

	default:
		debugger = logger.NewStructLogger(config)
	}
	// Load the test content from the input file
	src, err := os.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}
	var tests map[string]tests.StateTest
	if err = json.Unmarshal(src, &tests); err != nil {
		return err
	}
	// Iterate over all the tests, run them and aggregate the results
	cfg := vm.Config{
		Tracer:           tracer,
		Debug:            ctx.Bool(DebugFlag.Name) || ctx.Bool(MachineFlag.Name),
		EWASMInterpreter: ctx.String(stateTestEVMCEWASMFlag.Name),
	}
	results := make([]StatetestResult, 0, len(tests))
	for key, test := range tests {
		for _, st := range test.Subtests(nil) {
			if ctx.String(stateTestForkFlag.Name) != "" && ctx.String(stateTestForkFlag.Name) != st.Fork {
				continue
			}
			// Run the test and aggregate the result
			result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
			_, s, err := test.Run(st, cfg, false)
			// print state root for evmlab tracing
			if ctx.Bool(MachineFlag.Name) && s != nil {
				fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%x\"}\n", s.IntermediateRoot(false))
			}
			if err != nil {
				// Test failed, mark as so and dump any state to aid debugging
				result.Pass, result.Error = false, err.Error()
				if ctx.Bool(DumpFlag.Name) && s != nil {
					dump := s.RawDump(nil)
					result.State = &dump
				}
			}

			results = append(results, *result)

			// Print any structured logs collected
			if ctx.Bool(DebugFlag.Name) {
				if debugger != nil {
					fmt.Fprintln(os.Stderr, "#### TRACE ####")
					logger.WriteTrace(os.Stderr, debugger.StructLogs())
				}
			}
		}
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	return nil
}
