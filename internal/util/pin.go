package util

import (
	"crypto/sha1"
	"fmt"
	"math"
	"time"
	"strconv"
)

func GenerateRoomID(digits int, id string) (int, error) {
	// Validate the number of digits
	if digits < 1 || digits > 10 {
		return 0, fmt.Errorf("digits must be between 1 and 10")
	}

	// Combine ID and UNIX timestamp and hash it
	combined := fmt.Sprintf("%s%d", id, time.Now().UnixNano())
	hash := sha1.Sum([]byte(combined))

	// Convert the first 6 bytes of the hash to an integer
	hashInt := int64(hash[0])<<40 | int64(hash[1])<<32 | int64(hash[2])<<24 | int64(hash[3])<<16 | int64(hash[4])<<8 | int64(hash[5])
	modulo := int64(math.Pow10(digits))
	uniqueNumber := hashInt % modulo

	// Format the number to have the desired number of digits
	format := fmt.Sprintf("%%0%dd", digits)
	roomID := fmt.Sprintf(format, uniqueNumber)

	return strconv.Atoi(roomID)
}
