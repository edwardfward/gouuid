package uuid

import (
	"testing"
)

const NilUUID = "00000000-0000-0000-0000-000000000000"

func TestPrintUUID(t *testing.T) {

	// we need to test whether we can return an empty UUID
	emptyUUID := PrintUUID(nil)
	if emptyUUID != NilUUID {
		t.Errorf("Failed to print a null UUID. Expected: %s, "+
			"Received: %s", NilUUID, emptyUUID)
	}

	// generate 10000 UUIDs and make sure none of them match
	lastUUID := PrintUUID(NewV1())
	for i := 0; i < 10000; i++ {
		newUUID := PrintUUID(NewV1())
		if newUUID == lastUUID {
			t.Errorf("Duplicate UUIDs detected on test %d", i)
		}
	}
}

func TestNewV1(t *testing.T) {
	result := NewV1()
	if result == nil {
		t.Fatalf("returned a nil byte array")
	}

	// check version id is 1
	if result[6]>>4 != 1 {
		t.Fatalf("incorrect version number detected")
	}

	// check clock bits set correctly
	if result[8]>>6 != 2 {
		t.Fatalf("incorrect clock sequence detected")
	}

	// check string properly formatted for UUID
	for i := 0; i < 10; i++ {
		t.Log(PrintUUID(NewV1()))
	}
}

func TestNewV3(t *testing.T) {
	result := NewV3(u.namespace, "test")
	if result == nil {
		t.Fatalf("returned a nil byte array")
	}

	// check version is 3
	if result[6]>>4 != 3 {
		t.Fatalf("incorrect version number detected")
	}

	// check clock sequence bits set correctly
	if result[8]>>6 != 2 {
		t.Fatalf("incorrect clock sequence detected")
	}

	// check string properly formatted for UUID
	for i := 0; i < 10; i++ {
		t.Log(PrintUUID(result))
	}
}

func TestNewV4(t *testing.T) {
	result := NewV4()
	if result == nil {
		t.Fatalf("returned a nil byte array")
	}

	// check version is 4
	if result[6]>>4 != 4 {
		t.Fatalf("incorrect version number detected")
	}

	// check clock sequence bits set correctly
	if result[8]>>6 != 2 {
		t.Fatalf("incorrect clock sequence detected")
	}

	// check string properly formatted for UUID
	for i := 0; i < 10; i++ {
		t.Log(PrintUUID(NewV4()))
	}
}

func TestNewV5(t *testing.T) {
	result := NewV5(u.namespace, "test")
	if result == nil {
		t.Fatalf("returned a nil byte array")
	}

	// check version is 5
	if result[6]>>4 != 5 {
		t.Fatalf("incorrect version number detected")
	}

	// check clock sequence bits set correctly
	if result[8]>>6 != 2 {
		t.Fatalf("incorrect clock sequence detected")
	}

	// check string properly formatted for UUID
	for i := 0; i < 10; i++ {
		t.Log(PrintUUID(NewV4()))
	}
}

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV1()
	}
}

func BenchmarkNewV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV3(u.namespace, "test")
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV5(u.namespace, "test")
	}
}

func BenchmarkPrintUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PrintUUID(NewV1())
	}
}
