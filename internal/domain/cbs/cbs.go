// Package cbs provides core banking system functionality.
package cbs

// Status defines the core banking system status entity.
type Status struct {
	// SystemDate is the current system date in the core banking system.
	SystemDate string
	// IsEOD defines whether the EOD process is started.
	IsEOD bool
	// IsStandIn defines whether the EOD process is a stand-in process.
	IsStandIn bool
}

// NotReady returns true if the EOD process started and not in stand-in mode.
func (s *Status) NotReady() bool {
	return s.IsEOD && !s.IsStandIn
}
