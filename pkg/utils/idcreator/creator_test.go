// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idcreator

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCreator(t *testing.T) {

	currentLogIdCreator := idCreator

	InitCreator(uint16(rand.Int31()))

	assert.NotEqual(t, currentLogIdCreator, idCreator)
}

func TestNextID(t *testing.T) {

	id1, err := idCreator.NextID()
	assert.Nil(t, err)
	assert.Greater(t, id1, uint64(0))

	id2, err := idCreator.NextID()
	assert.Nil(t, err)
	assert.Greater(t, id2, uint64(0))

	assert.NotEqual(t, id1, id2)
}
