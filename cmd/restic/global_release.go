// +build !debug,!profile

package main

// runDebug is a noop without the debug tag.
func runDebug() error { return nil }
