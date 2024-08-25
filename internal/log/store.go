package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

// Represents a Store File
// this i s the file we store records in
type store struct {
	*os.File
	mu sync.Mutex
	buf *bufio.Writer
	size uint64
}

// newStore initializes a store from a given file
// if the file already exists, we can recreate it with
// its contents
func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store {
		File: f,
		size: size,
		buf: bufio.NewWriter(f),
	}, nil
}

// Write to the store
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size // current size

	// Write the length of the bytes onto the buffer
	// This comes in handy when reading, each "offset" 
	// should be a marker that tells you how many bytes
	// we will need to read to get all of the content
	if err := binary.Write(s.buf, enc, uint64 (len(p))); err != nil {
		return 0,0, err
	}

	// Write the actual bytes to the buffer
	w, err := s.buf.Write(p)
	if err != nil {
		return 0,0,err
	}
	w += lenWidth

	// update the size by the # of bytes we wrote
	s.size += uint64(w)

	// return the number of bytes we wrote, the position of the
	// written entry
	return uint64(w), pos, nil
}

// Returns the record at a given position in the store
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//Flush the buffer beforehand
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// allocate a byte array of size word
	size := make([]byte, lenWidth)

	// Recall that pos is the marker that indicates
	// how many bytes of data that follows this position.
	// Update size to be that # of bytes that pos is marking
	if _, err := s.File.ReadAt(size, int64 (pos)); err != nil {
		return nil, err
	}
	
	// Now read the actual # of bytes worth of data
	b := make([]byte, enc.Uint64(size))
	if _,err := s.File.ReadAt(b, int64 (pos + lenWidth)); err != nil {
		return nil, err
	}

	// Return the data in bytes
	return b, nil
}

// Read len(p) # of bytes at the specific offset in Store
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Flushing the buffer
	if err := s.buf.Flush(); err != nil {
		return 0,err
	}

	return s.File.ReadAt(p, off)
}

// Close any persisted buffer data before closing the file 
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}