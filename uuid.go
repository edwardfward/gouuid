package uuid

// reference https://tools.ietf.org/html/rfc4122#section-4.2.1

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type uuid struct {
	sync.Mutex
	timestamp uint64
	clock     uint16
	count     uint32
	node      []byte
}

// RFC 4122 uses the start of the Gregorian calendar (15 October, 1582) but
// Go's time.Since is limited to 256 years (max value of int64 nanoseconds).
// UUID uses the UNIX calendar (1 January, 1970) instead.
var UnixDate = time.Date(1970,
	time.January, 1, 0, 0, 0, 0,
	time.UTC)

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
		rand.Read(randomNode)
		// create a 48-bit number for random node
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

func getNanos100s() uint64 {
	return uint64(time.Since(UnixDate).Nanoseconds() / 100)
}

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

func PrintUUID(uuid []byte) string {
	if uuid == nil {
		uuid = make([]byte, 16)
	}

	return fmt.Sprintf("%0.8x-%0.4x-%0.4x-%0.2x%0.2x-%0.12x",
		uuid[0:4], uuid[4:6], uuid[6:8],
		uuid[8], uuid[9], uuid[10:16])
}