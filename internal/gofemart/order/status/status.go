package status

import "golang.org/x/exp/slices"

const INVALID = "INVALID"
const PROCESSING = "PROCESSING"
const PROCESSED = "PROCESSED"
const NEW = "NEW"

func IsSupported(status string) bool {
	return slices.Contains([]string{INVALID, PROCESSED, PROCESSING, NEW}, status)
}

func RequiresTracking(status string) bool {
	return slices.Contains([]string{PROCESSING, NEW}, status)
}
