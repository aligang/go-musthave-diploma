package status

const INVALID = "INVALID"
const PROCESSING = "PROCESSING"
const PROCESSED = "PROCESSED"
const NEW = "NEW"

func IsSupported(status string) bool {
	for _, s := range []string{INVALID, PROCESSED, PROCESSING, NEW} {
		if status == s {
			return true
		}
	}
	return false
}

func RequiresTracking(status string) bool {
	for _, s := range []string{PROCESSING, NEW} {
		if status == s {
			return true
		}
	}
	return false
}
