package log

import (
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

