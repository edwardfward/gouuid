// Copyright 2019 Edward F. Ward III.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package uuid

// reference https://tools.ietf.org/html/rfc4122#section-4.2.1

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// RFC 4122 uses the Gregorian epoch (15 October 1582) for calculating 100s
// of nanoseconds. Since Go's Time.Duration returns an int64 for nanoseconds,
// which is limited to 250ish years of nanoseconds we need to perform the
// intermediate step of calculating 100s of nanoseconds between the Gregorian
// epoch and Unix epoch. We determine the number of days between 15 October
// 1582 (Julian 2299171) and 1 January 1970 (Julian 2440601) and then multiply
// the number of seconds by 1e7. Nanoseconds are 1e9 but since we are dividing
// by 100, we only need to multiply 1e7).

// Julian dates calculated from http://numerical.recipes/julian.html

const (
	gregorianEpochJulianDays = 2299171 // 15 October 1582
	unixEpochJulianDays      = 2440601 // 1 January 1970
)

var epochDiffNanos100s = uint64((unixEpochJulianDays - gregorianEpochJulianDays) *
	(24 * 60 * 60) * 1e7)

type uuid struct {
	sync.Mutex
	timestamp uint64
	clock     uint16
	count     uint32
	node      []byte
}

var u = uuid{
	timestamp: getNanos100s(),
	clock:     uint16(rand.Uint32()),
	count:     0,
}

func init() {
	// read network interfaces
	interfaces, err := net.Interfaces()
	// if unable to read interfaces, set to random
	if err != nil {
		randomNode := make([]byte, 6)
		// create 48-bit random bits
		rand.Read(randomNode)
		// check to ensure the most significant bit of the random bits is 1
		randomNode[0] = randomNode[0] | 128
		u.node = randomNode
	}

	for _, inter := range interfaces {
		if len(inter.HardwareAddr) != 6 {
			continue
		} else {
			u.node = inter.HardwareAddr
			break
		}
	}
}

func uint32ToBytes(val uint32) []byte {
	result := make([]byte, 4)
	result[0] = byte(val >> 24)
	result[1] = byte(val >> 16)
	result[2] = byte(val >> 8)
	result[3] = byte(val)
	return result
}

func uint16ToBytes(val uint16) []byte {
	result := make([]byte, 2)
	result[0] = byte(val >> 8)
	result[1] = byte(val)
	return result
}

// getNanos100s calculates the 100s of nanoseconds between now(UTC) and the
// Gregorian calendar epoch. Returns 100s of nanoseconds.
func getNanos100s() uint64 {
	return epochDiffNanos100s + uint64(time.Now().In(time.UTC).UnixNano()/100)
}

// NewV1 generates a RFC 4122 Version 1 compliant UUID. Returns 128-bit / 16
// byte array representing the UUID.
func NewV1() []byte {

	u.Lock()
	newTime := getNanos100s()

	if newTime > u.timestamp {
		u.clock++
		u.timestamp = newTime
		u.count = 0
	} else {
		// A high resolution timestamp can be simulated by keeping a count of
		// the number of UUIDs that have been generated with the same value of
		// the system time, and using it to construct the low order bits of the
		// timestamp.  The count will range between zero and the number of
		// 100-nanosecond intervals per system time interval.
		u.count++
		newTime += uint64(u.count)
	}
	clockSequence := u.clock
	u.Unlock()

	timeLow := uint32(0xFFFFFFFF & newTime)
	timeMid := uint16((newTime >> 32) & 0xFFFF)
	timeHiAndVersion := uint16(((newTime >> 48) & 0x0FFF) | 0x1000)
	clockSeqHiAndReserved := uint8((clockSequence >> 8 & 0x3F) | 0x80)
	clockSeqLow := uint8(clockSequence & 0xFF)

	result := make([]byte, 0, 0)[:]
	result = append(result, uint32ToBytes(timeLow)...)
	result = append(result, uint16ToBytes(timeMid)...)
	result = append(result, uint16ToBytes(timeHiAndVersion)...)
	result = append(result, byte(clockSeqHiAndReserved))
	result = append(result, byte(clockSeqLow))
	result = append(result, u.node...)

	return result
}

// NewV4 generates a RFC 4122 Version 4 compliant UUID. Returns 128-bit / 16
// byte array representing the UUID.
func NewV4() []byte {
	/*
		1. Set all the other bits to randomly (or pseudo-randomly) chosen
		values.

		2. Set the two most significant bits (bits 6 and 7) of the
		clock_seq_hi_and_reserved (byte 8) to zero and one, respectively.

		3. Set the four most significant bits (bits 12 through 15) of the
		time_hi_and_version field to the 4-bit version number
	*/

	result := make([]byte, 16)
	rand.Read(result)                     // step 1
	result[8] = (result[8] & 0x3F) | 0x80 // step 2
	result[6] = (result[6] & 0x0F) | 0x40 // step 3

	return result
}

//PrintUUID returns properly formatted UUID string for any RFC 4122 version,
//including the nil UUID.
func PrintUUID(uuid []byte) string {
	if uuid == nil {
		uuid = make([]byte, 16)
	}

	return fmt.Sprintf("%0.8x-%0.4x-%0.4x-%0.2x%0.2x-%0.12x",
		uuid[0:4], uuid[4:6], uuid[6:8],
		uuid[8], uuid[9], uuid[10:16])
}
