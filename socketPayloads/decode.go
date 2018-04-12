package socketPayloads

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type parser struct {
	err           error
	initialArg    json.RawMessage
	base64Decoded []byte
	buffer        bytes.Buffer
	unzipped      []byte
}

func (p *parser) parse(arg json.RawMessage) ([]byte, error) {
	p.initialArg = arg

	p.base64Decode()
	p.createBuffer()
	p.unzip()

	return p.unzipped, p.err
}

func (p *parser) base64Decode() {

	var base64EncodedStr string

	parseErr := json.Unmarshal(p.initialArg, &base64EncodedStr)

	if parseErr != nil {
		panic(parseErr)
	}

	var decodeErr error
	p.base64Decoded, decodeErr = base64.StdEncoding.DecodeString(base64EncodedStr)

	if decodeErr != nil {
		p.err = fmt.Errorf("unable to decode response contents: %+v", decodeErr)
	}
}

func (p *parser) createBuffer() {
	if p.err != nil {
		return
	}

	_, bufWriteErr := p.buffer.Write(p.base64Decoded)
	if bufWriteErr != nil {
		p.err = fmt.Errorf("unable to write decoded string to buffer: %+v", bufWriteErr)
	}
}

func (p *parser) unzip() {
	if p.err != nil {
		return
	}

	var buf bytes.Buffer

	zr := flate.NewReader(&p.buffer)

	_, readErr := io.Copy(&buf, zr)

	if readErr != nil {
		p.err = fmt.Errorf("unable to unzip response contents. %+v", readErr)
		return
	}

	p.unzipped = make([]byte, buf.Len())

	_, readErr2 := buf.Read(p.unzipped)

	if readErr2 != nil {
		p.unzipped = nil
		p.err = fmt.Errorf("unable to read from buffer %+v", readErr2)
		return
	}

}

//Parse engage the socketpayload parser.
func Parse(arg json.RawMessage) ([]byte, error) {
	parser := parser{}
	return parser.parse(arg)
}
