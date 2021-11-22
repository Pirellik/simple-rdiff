package librsync

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Pirellik/simple-rdiff/rollsum"
)

type Signature struct {
	blockLength             uint32
	strongSignatures        [][]byte
	weakSignaturesToBlockID map[uint32]uint64
}

func NewSignature(in io.Reader, blockLen uint32) (*Signature, error) {
	if blockLen < sha256.Size {
		return nil, fmt.Errorf("too small block size, min size = %d", sha256.Size)
	}
	sig := Signature{
		blockLength:             blockLen,
		weakSignaturesToBlockID: make(map[uint32]uint64),
	}
	buffer := make([]byte, blockLen)

	for {
		n, err := in.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		block := buffer[:n]
		weakSig := computeRollingChecksum(block)
		strongSig := computeStrongChecksum(block)
		sig.weakSignaturesToBlockID[weakSig] = uint64(len(sig.strongSignatures))
		sig.strongSignatures = append(sig.strongSignatures, strongSig)
	}
	return &sig, nil
}

func ReadSignature(in io.Reader) (*Signature, error) {
	var blockLength uint32
	if err := binary.Read(in, binary.BigEndian, &blockLength); err != nil {
		return nil, err
	}
	strongSigs := [][]byte{}
	weakSigToBlockID := map[uint32]uint64{}
	for {
		var weakSig uint32
		if err := binary.Read(in, binary.BigEndian, &weakSig); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		strongSig := make([]byte, sha256.Size)
		n, err := in.Read(strongSig)
		if err != nil {
			return nil, err
		}
		if n != int(sha256.Size) {
			return nil, fmt.Errorf("too short strong hash, got = %d, want = %d", n, sha256.Size)
		}
		weakSigToBlockID[weakSig] = uint64(len(strongSigs))
		strongSigs = append(strongSigs, strongSig)
	}

	return &Signature{
		blockLength:             blockLength,
		strongSignatures:        strongSigs,
		weakSignaturesToBlockID: weakSigToBlockID,
	}, nil
}

func (s *Signature) Write(out io.Writer) error {
	if err := binary.Write(out, binary.BigEndian, s.blockLength); err != nil {
		return err
	}
	weakSigs := make([]uint32, len(s.strongSignatures))
	for sig, ID := range s.weakSignaturesToBlockID {
		weakSigs[int(ID)] = sig
	}
	for i, weakSig := range weakSigs {
		if err := binary.Write(out, binary.BigEndian, weakSig); err != nil {
			return err
		}
		if _, err := out.Write(s.strongSignatures[i]); err != nil {
			return err
		}
	}
	return nil
}

func computeRollingChecksum(in []byte) uint32 {
	rSum := rollsum.New()
	rSum.Init(in)
	return rSum.Sum()
}

func computeStrongChecksum(in []byte) []byte {
	strongHash := sha256.New()
	strongHash.Write(in)
	return strongHash.Sum(nil)
}
