package librsync

import (
	"bytes"
	"io"

	"github.com/Pirellik/simple-rdiff/rollsum"
)

type Delta struct {
	chunks []chunk
}

func NewDelta(in io.Reader, s *Signature) (*Delta, error) {
	delta := Delta{}
	rSum := rollsum.New()
	block := make([]byte, s.blockLength)
	singleByte := make([]byte, 1)
	prevChunkType := chunkTypeReusable

	for {
		if prevChunkType == chunkTypeReusable {
			n, err := in.Read(block)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			if n < len(block) {
				block = block[:n]
			}
			rSum.Init(block)
		} else {
			if _, err := in.Read(singleByte); err != nil {
				if err == io.EOF {
					delta.addChunk(&modified{data: block[1:]})
					break
				}
				return nil, err
			}
			rSum.Roll(block[0], singleByte[0])
			block = append(block[1:], singleByte...)
		}
		chunk := getChunk(block, rSum.Sum(), s)
		delta.addChunk(chunk)
		prevChunkType = chunk.chunkType()
	}
	return &delta, nil
}

func ReadDelta(in io.Reader) (*Delta, error) {
	delta := Delta{}
	for {
		chunk, err := readChunk(in)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		delta.chunks = append(delta.chunks, chunk)
	}
	return &delta, nil
}

func (d *Delta) Patch(base io.ReadSeeker, out io.Writer) error {
	for _, c := range d.chunks {
		if err := c.patch(base, out); err != nil {
			return err
		}
	}
	return nil
}

func (d *Delta) Write(out io.Writer) error {
	for _, c := range d.chunks {
		if err := c.write(out); err != nil {
			return err
		}
	}
	return nil
}

func (d *Delta) addChunk(c chunk) {
	if len(d.chunks) == 0 {
		d.chunks = append(d.chunks, c)
		return
	}
	if d.chunks[len(d.chunks)-1].chunkType() == c.chunkType() {
		d.chunks[len(d.chunks)-1].append(c)
		return
	}
	d.chunks = append(d.chunks, c)
}

func getChunk(block []byte, weakSum uint32, s *Signature) chunk {
	blockID, ok := s.weakSignaturesToBlockID[weakSum]
	if ok && bytes.Equal(s.strongSignatures[int(blockID)], computeStrongChecksum(block)) {
		return &reusable{
			startPosition: blockID * uint64(len(block)),
			length:        uint64(len(block)),
		}
	}
	return &modified{
		data: block[:1],
	}
}
