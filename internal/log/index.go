package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth
)

// represents the index file
type index struct {
	file *os.File
	mmap gommap.MMap //memory mapped file
	size uint64 //size of the entry
}

// newIndex creates an index for a given file
func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index {
		file: f,
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	// Save the current size of the file to track amt of data in 
	// the index file.

	// If the current size is > max index size, we will memory map
	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),); err != nil {
			return nil,err
	}

	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ | gommap.PROT_WRITE,
		gommap.MAP_SHARED,); err != nil {
			return nil, err
	}
	return idx, nil
}
 
func (i *index) Close() error {

	// syncs memory-mapped data to persisted file
	// and flushes the persisted file 
	if err := i.mmap.Sync(gommap.MS_ASYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}

	// truncates persisted file to only be as big
	// as the amount of data it contains 
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	
	return i.file.Close()
}

// takes in an offset and returns the records position in the Store file
func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	// if empty, then return err
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	if in == -1 {
		out = uint32((i.size / entWidth) - 1) // get last index
	} else {
		out = uint32(in)
	}
	pos = uint64(out) * entWidth // given the index, get the byte offset
	
	// if out of bounds, return err
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}

	out = enc.Uint32(i.mmap[pos : pos+offWidth]) // read the offset width amt of bytes
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth]) // read the poswidth amt of bytes
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}

	// write the offset and the position onto the memory map
	// offset is size offWidth (4)
	// position is size entWidth  (8)
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += uint64(entWidth) //increment the size of the file by the new entry
	return nil
}


func (i *index) Name() string {
	return i.file.Name()
}
