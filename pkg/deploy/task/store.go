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

package task

import (
	"fmt"
	"sync"
)

type Store interface {
	GetTask(name string) Task
	// Add a task, if the task already exists (task name is same), return an error
	AddTask(task Task) error
	// Update a task, if the task doesn't exist (task name is same), return an error
	UpdateTask(task Task) error
	// Update a task, if the task doesn't exist (task name is same), add it.
	UpdateOrAddTask(task Task) error
}

// A Store implementation via map
type cache struct {
	sync.RWMutex
	m map[string]Task
}

func (c *cache) GetTask(name string) Task {
	c.RLock()
	defer c.RUnlock()

	return c.m[name]
}

func (c *cache) AddTask(task Task) error {
	name := task.GetName()
	if name == "" {
		return fmt.Errorf("failed to add task: task name can't be empty")
	}

	c.Lock()
	defer c.Unlock()

	_, ok := c.m[name]
	if ok {
		return fmt.Errorf("failed to add task: task is already existed")
	}

	c.m[name] = task

	return nil
}

func (c *cache) UpdateTask(task Task) error {
	name := task.GetName()
	if name == "" {
		return fmt.Errorf("failed to update task: task name can't be empty")
	}

	c.Lock()
	defer c.Unlock()

	_, ok := c.m[name]
	if !ok {
		return fmt.Errorf("failed to update task: task doesn't exist")
	}

	c.m[name] = task

	return nil
}

func (c *cache) UpdateOrAddTask(task Task) error {
	name := task.GetName()
	if name == "" {
		return fmt.Errorf("failed to update or create task: task name can't be empty")
	}

	c.Lock()
	defer c.Unlock()

	c.m[name] = task

	return nil
}

// global cache store
var cacheStore *cache

func init() {
	cacheStore = &cache{
		m: make(map[string]Task),
	}
}

func GetGlobalCacheStore() Store {
	return cacheStore
}
