// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package demotools provides a set of tools that help you write code examples.
//
// **Pausable**
//
// The pausable interface creates an easy to mock pausing object for testing.
package demotools

import "time"

// IPausable defines the interface for pausable objects.
type IPausable interface {
	Pause(secs int)
}

// Pause holds the pausable object.
type Pause struct{}

// Pause waits for the specified number of seconds.
func (pausable Pause) Pause(secs int) {
	time.Sleep(time.Duration(secs) * time.Second)
}
