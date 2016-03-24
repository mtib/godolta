package deltal

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

// Encoder of delta-l
type Encoder struct {
	Counter     uint64
	passarray   []byte
	FileReader  io.ReadSeeker
	Checksum    []byte
	UseChecksum bool
	headerpos   int
	last        uint8
}

// NewEncoderReader creates an Encoder
func NewEncoderReader(reader io.ReadSeeker, password string, checksum bool) (*Encoder, error) {
	full, err := ioutil.ReadAll(reader)
	reader.Seek(0, 0)
	return &Encoder{passarray: stringToArray(password), FileReader: reader, UseChecksum: checksum, Checksum: rustHash(full)}, err
}

// NewEncoder returns pointer to decoder reader interface
func NewEncoder(file, password string, checksum bool) (*Encoder, error) {
	f, _ := os.Open(file)
	return NewEncoderReader(f, password, checksum)
}

type compressedByte struct {
	data   []byte
	offset int64
}

func (c *compressedByte) Read(b []byte) (n int, err error) {
	begin := c.offset
	c.offset += int64(len(b))
	if c.offset > int64(len(c.data)) {
		err = io.EOF
		c.offset = int64(len(c.data))
	}
	n = int(c.offset - begin)
	for k := begin; k < c.offset; k++ {
		b[k] = c.data[k]
	}
	return
}

func (c *compressedByte) Seek(delta int64, from int) (to int64, err error) {
	switch from {
	case 0:
		c.offset = delta
	case 2:
		c.offset = int64(len(c.data) - 1) // minus one?
	case 1:
		c.offset += delta
	}
	if c.offset < 0 || c.offset > int64(len(c.data)) {
		err = io.ErrUnexpectedEOF
	}
	to = c.offset
	return
}

// NewCompressedEncoderReader gzips the file before encryption
func NewCompressedEncoderReader(reader io.ReadSeeker, password string, checksum bool) (*Encoder, error) {
	var b bytes.Buffer
	gzwriter := gzip.NewWriter(&b)
	io.Copy(gzwriter, reader)
	gzwriter.Flush()
	gzwriter.Close()
	buf, err := ioutil.ReadAll(&b)
	return &Encoder{passarray: stringToArray(password), FileReader: &compressedByte{buf, 0}, UseChecksum: checksum, Checksum: rustHash(buf)}, err
}

func stringToArray(p string) (pa []byte) {
	if p != "" {
		pa = rustHash([]byte(p))
	} else {
		pa = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return
}

// Reset resets the encoder, synonym to Seek(0,0)
func (d *Encoder) Reset() {
	d.Counter = 0
	d.headerpos = 0
	d.FileReader.Seek(0, 0)
	d.last = 0
}

func (d *Encoder) Read(b []byte) (n int, err error) {
	// Writing Checksum
	for (d.UseChecksum && d.headerpos < 12) || (!d.UseChecksum && d.headerpos < 4) {
		if d.headerpos > 3 {
			if d.UseChecksum {
				b[n] = d.Checksum[d.headerpos-4]
			} else {
				return
			}
		} else {
			b[n] = header(d.headerpos, d.UseChecksum)
		}
		d.headerpos++
		n++
		if n == len(b) {
			return
		}
	}

	// Encrypting File
	b2 := make([]byte, len(b)-n)
	n2, err := d.FileReader.Read(b2)
	for k := range b2 {
		s := b2[k]
		b[k+n] = uint8(b2[k] + d.last + d.passarray[d.Counter%uint64(len(d.passarray))])
		d.Counter++
		d.last = s
	}
	n += n2
	return
}

// FastEncrypt is an easy-call function
func FastEncrypt(file, pass string) {
	filein, _ := os.Open(file)
	encoder, _ := NewEncoderReader(filein, pass, true)
	encryptedData, _ := ioutil.ReadAll(encoder)
	ioutil.WriteFile(file+".delta", encryptedData, os.ModePerm)
}

func benchEncrypter(file, pass string) (func(), *os.File) {
	filein, _ := os.Open(file)
	encoder, _ := NewEncoderReader(filein, pass, true)
	return func() {
		ioutil.ReadAll(encoder)
		encoder.Reset()
	}, filein
}

var (
	checkHeader    = []byte{206, 148, 76, 10}
	nonCheckHeader = []byte{206, 148, 108, 10}
)

func header(pos int, check bool) uint8 {
	if check {
		return checkHeader[pos]
	}
	return nonCheckHeader[pos]
}
