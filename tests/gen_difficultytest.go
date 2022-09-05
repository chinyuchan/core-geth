// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

// Copyright 2022 the core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ = (*difficultyTestMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (d DifficultyTest) MarshalJSON() ([]byte, error) {
	type DifficultyTest struct {
		ParentTimestamp    math.HexOrDecimal64   `json:"parentTimestamp"`
		ParentDifficulty   *math.HexOrDecimal256 `json:"parentDifficulty"`
		UncleHash          common.Hash           `json:"parentUncles"`
		CurrentTimestamp   math.HexOrDecimal64   `json:"currentTimestamp"`
		CurrentBlockNumber math.HexOrDecimal64   `json:"currentBlockNumber"`
		CurrentDifficulty  *math.HexOrDecimal256 `json:"currentDifficulty"`
		// Chainspec          chainspecRef          `json:"chainspec,omitempty"`
		// Name               string                `json:"name,omitempty"`
	}
	var enc DifficultyTest
	enc.ParentTimestamp = math.HexOrDecimal64(d.ParentTimestamp)
	enc.ParentDifficulty = (*math.HexOrDecimal256)(d.ParentDifficulty)
	enc.UncleHash = d.UncleHash
	enc.CurrentTimestamp = math.HexOrDecimal64(d.CurrentTimestamp)
	enc.CurrentBlockNumber = math.HexOrDecimal64(d.CurrentBlockNumber)
	enc.CurrentDifficulty = (*math.HexOrDecimal256)(d.CurrentDifficulty)
	// enc.Chainspec = d.Chainspec
	// enc.Name = d.Name
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (d *DifficultyTest) UnmarshalJSON(input []byte) error {
	type DifficultyTest struct {
		ParentTimestamp    *math.HexOrDecimal64  `json:"parentTimestamp"`
		ParentDifficulty   *math.HexOrDecimal256 `json:"parentDifficulty"`
		UncleHash          *common.Hash          `json:"parentUncles"`
		CurrentTimestamp   *math.HexOrDecimal64  `json:"currentTimestamp"`
		CurrentBlockNumber *math.HexOrDecimal64  `json:"currentBlockNumber"`
		CurrentDifficulty  *math.HexOrDecimal256 `json:"currentDifficulty"`
		// Chainspec          *chainspecRef         `json:"chainspec,omitempty"`
		// Name               *string               `json:"name,omitempty"`
	}
	var dec DifficultyTest
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentTimestamp != nil {
		d.ParentTimestamp = uint64(*dec.ParentTimestamp)
	}
	if dec.ParentDifficulty != nil {
		d.ParentDifficulty = (*big.Int)(dec.ParentDifficulty)
	}
	if dec.UncleHash != nil {
		d.UncleHash = *dec.UncleHash
	}
	if dec.CurrentTimestamp != nil {
		d.CurrentTimestamp = uint64(*dec.CurrentTimestamp)
	}
	if dec.CurrentBlockNumber != nil {
		d.CurrentBlockNumber = uint64(*dec.CurrentBlockNumber)
	}
	if dec.CurrentDifficulty != nil {
		d.CurrentDifficulty = (*big.Int)(dec.CurrentDifficulty)
	}
	// if dec.Chainspec != nil {
	// 	d.Chainspec = *dec.Chainspec
	// }
	// if dec.Name != nil {
	// 	d.Name = *dec.Name
	// }
	return nil
}
