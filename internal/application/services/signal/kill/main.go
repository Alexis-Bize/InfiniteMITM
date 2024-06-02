// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package MITMApplicationKillSignalService

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	cleanupFns      []func()
	cleanupRegMutex sync.Mutex
	initOnce        sync.Once
)

func Init() {
	initOnce.Do(handleSignals)
}

func Register(cleanupFn func()) {
	cleanupRegMutex.Lock()
	defer cleanupRegMutex.Unlock()
	cleanupFns = append(cleanupFns, cleanupFn)
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-c
		log.Printf("Received signal: %v", sig)
		executeCleanup()
		os.Exit(0)
	}()
}

func executeCleanup() {
	cleanupRegMutex.Lock()
	defer cleanupRegMutex.Unlock()
	for _, fn := range cleanupFns {
		fn()
	}
}