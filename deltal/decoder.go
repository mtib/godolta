package deltal

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Decoder of delta-l files
type Decoder struct {
	Stream      io.ReadSeeker
	passhash    []byte
	Checksum    []byte
	UseChecksum bool
	Offset      uint64
	passOffset  int
	last        uint8
}

// Init initializes the decoder instance
func (d *Decoder) Init() {
	if d.Offset == 0 {
		h1 := make([]byte, 4)
		d.Stream.Read(h1)
		d.UseChecksum = (h1[2] == 76)
		d.Offset += 4
		if d.UseChecksum {
			d.Checksum = make([]byte, 8)
			d.Stream.Read(d.Checksum)
			d.Offset += 8
		}
	}
}

func (d *Decoder) String() string {
	return fmt.Sprintf("Decoder: (Offset: %X) using Passhash: %v Checksum: %v (%v)", d.Offset, d.passhash, d.Checksum, d.UseChecksum)
}

// Seek as defined in io.Seeker interface
func (d *Decoder) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		d.Offset = uint64(offset)
	case 2:
		ioutil.ReadAll(d)
		fallthrough
	case 1:
		d.Offset += uint64(offset)
	}
	d.Stream.Seek(offset, whence)
	return int64(d.Offset), nil
}

// Read decrypts the Stream of the Decoder,
// make sure to call Init() before trying to Read().
// If you created the Decoder with NewDecoder(), Init() was already called
func (d *Decoder) Read(b []byte) (n int, err error) {
	n, err = d.Stream.Read(b)
	for k := range b {
		if k == n {
			break
		}
		b[k] = byte(b[k] - d.last - d.passhash[(d.passOffset+k)%len(d.passhash)])
		d.last = b[k]
	}
	d.passOffset = (d.passOffset + n) % len(d.passhash)
	d.Offset += uint64(n)
	return
}

// Check compares the Checksum in the file with the actual checksum
func (d *Decoder) Check(file []byte) bool {
	for k, v := range rustHash(file) {
		if v != d.Checksum[k] {
			return false
		}
	}
	return true
}

// NewDecoderStream return initialized delta-l Decoder reading the stream io.Reader
func NewDecoderStream(stream io.ReadSeeker, password string) *Decoder {
	var pass []byte
	if password == "" {
		pass = u64tobyte(0)
	} else {
		pass = rustHash([]byte(password))
	}
	decoder := &Decoder{Stream: stream, passhash: pass}
	decoder.Init()
	return decoder
}
