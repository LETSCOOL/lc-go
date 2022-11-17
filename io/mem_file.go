package io

import (
	"fmt"
	"io"
	"io/fs"
	"sync"
	"time"
)

func NewRoMemFile(name string, data []byte) *RoMemFile {
	f := RoMemFile{
		name:        name,
		data:        data,
		createdTime: time.Now(),
	}
	return &f
}

type RoMemFile struct {
	name        string
	data        []byte
	createdTime time.Time
	mutex       sync.Mutex
	pos         int
}

func (m *RoMemFile) Stat() (fs.FileInfo, error) {
	return m, nil
}

func (m *RoMemFile) Read(b []byte) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	n, err := m.readAt(b, int64(m.pos))
	m.pos += n
	return n, err
}

func (m *RoMemFile) readAt(b []byte, off int64) (int, error) {
	if off < 0 || int64(int(off)) < off {
		return 0, fmt.Errorf("invalid argument to read file(%s), off: %d", m.name, off)
	}
	if off > int64(len(m.data)) {
		return 0, io.EOF
	}
	n := copy(b, m.data[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

func (m *RoMemFile) Close() error {
	return nil
}

func (m *RoMemFile) Name() string {
	return m.name
}

func (m *RoMemFile) Size() int64 {
	return int64(len(m.data))
}

func (m *RoMemFile) Mode() fs.FileMode {
	return fs.ModeNamedPipe
}

func (m *RoMemFile) ModTime() time.Time {
	return m.createdTime
}

func (m *RoMemFile) IsDir() bool {
	return false
}

func (m *RoMemFile) Sys() any {
	return nil
}
