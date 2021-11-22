package librsync

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDelta(t *testing.T) {
	tests := []struct {
		desc      string
		giveInput *bytes.Buffer
		giveSig   *Signature
		wantDelta *Delta
	}{
		{
			desc: "should handle no changes",
			giveInput: bytes.NewBuffer([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
				101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
				108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32}),
			giveSig: &Signature{
				blockLength: 32,
				strongSignatures: [][]byte{
					{61, 7, 188, 146, 183, 66, 102, 5, 216, 249, 196, 2, 184, 114, 200, 118, 207, 233, 146, 244, 196, 82, 188, 82, 74, 178, 66, 250, 206, 163, 215, 240},
					{1, 84, 112, 6, 249, 182, 164, 120, 200, 26, 252, 211, 98, 67, 127, 254, 81, 223, 36, 86, 194, 26, 205, 54, 85, 246, 96, 23, 101, 215, 125, 41},
					{134, 176, 225, 187, 226, 118, 88, 57, 49, 158, 133, 226, 87, 193, 5, 129, 20, 56, 212, 158, 60, 234, 21, 240, 68, 11, 190, 154, 195, 62, 165, 28},
				},
				weakSignaturesToBlockID: map[uint32]uint64{
					16646287:   2,
					3276082140: 1,
					3308522449: 0,
				},
			},
			wantDelta: &Delta{
				chunks: []chunk{
					&reusable{
						startPosition: 0,
						length:        66,
					},
				},
			},
		},
		{
			desc: "should handle changes",
			giveInput: bytes.NewBuffer([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
				101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 102, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
				108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32}),
			giveSig: &Signature{
				blockLength: 32,
				strongSignatures: [][]byte{
					{61, 7, 188, 146, 183, 66, 102, 5, 216, 249, 196, 2, 184, 114, 200, 118, 207, 233, 146, 244, 196, 82, 188, 82, 74, 178, 66, 250, 206, 163, 215, 240},
					{1, 84, 112, 6, 249, 182, 164, 120, 200, 26, 252, 211, 98, 67, 127, 254, 81, 223, 36, 86, 194, 26, 205, 54, 85, 246, 96, 23, 101, 215, 125, 41},
					{134, 176, 225, 187, 226, 118, 88, 57, 49, 158, 133, 226, 87, 193, 5, 129, 20, 56, 212, 158, 60, 234, 21, 240, 68, 11, 190, 154, 195, 62, 165, 28},
				},
				weakSignaturesToBlockID: map[uint32]uint64{
					16646287:   2,
					3276082140: 1,
					3308522449: 0,
				},
			},
			wantDelta: &Delta{
				chunks: []chunk{
					&modified{
						data: []byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 102},
					},
					&reusable{
						startPosition: 32,
						length:        34,
					},
				},
			},
		},
		{
			desc: "should handle additions",
			giveInput: bytes.NewBuffer([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
				101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 19, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
				108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32}),
			giveSig: &Signature{
				blockLength: 32,
				strongSignatures: [][]byte{
					{61, 7, 188, 146, 183, 66, 102, 5, 216, 249, 196, 2, 184, 114, 200, 118, 207, 233, 146, 244, 196, 82, 188, 82, 74, 178, 66, 250, 206, 163, 215, 240},
					{1, 84, 112, 6, 249, 182, 164, 120, 200, 26, 252, 211, 98, 67, 127, 254, 81, 223, 36, 86, 194, 26, 205, 54, 85, 246, 96, 23, 101, 215, 125, 41},
					{134, 176, 225, 187, 226, 118, 88, 57, 49, 158, 133, 226, 87, 193, 5, 129, 20, 56, 212, 158, 60, 234, 21, 240, 68, 11, 190, 154, 195, 62, 165, 28},
				},
				weakSignaturesToBlockID: map[uint32]uint64{
					16646287:   2,
					3276082140: 1,
					3308522449: 0,
				},
			},
			wantDelta: &Delta{
				chunks: []chunk{
					&reusable{
						startPosition: 0,
						length:        32,
					},
					&modified{
						data: []byte{19},
					},
					&reusable{
						startPosition: 32,
						length:        34,
					},
				},
			},
		},
		{
			desc: "should handle removals",
			giveInput: bytes.NewBuffer([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
				101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
				108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32}),
			giveSig: &Signature{
				blockLength: 32,
				strongSignatures: [][]byte{
					{61, 7, 188, 146, 183, 66, 102, 5, 216, 249, 196, 2, 184, 114, 200, 118, 207, 233, 146, 244, 196, 82, 188, 82, 74, 178, 66, 250, 206, 163, 215, 240},
					{1, 84, 112, 6, 249, 182, 164, 120, 200, 26, 252, 211, 98, 67, 127, 254, 81, 223, 36, 86, 194, 26, 205, 54, 85, 246, 96, 23, 101, 215, 125, 41},
					{134, 176, 225, 187, 226, 118, 88, 57, 49, 158, 133, 226, 87, 193, 5, 129, 20, 56, 212, 158, 60, 234, 21, 240, 68, 11, 190, 154, 195, 62, 165, 28},
				},
				weakSignaturesToBlockID: map[uint32]uint64{
					16646287:   2,
					3276082140: 1,
					3308522449: 0,
				},
			},
			wantDelta: &Delta{
				chunks: []chunk{
					&modified{
						data: []byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104},
					},
					&reusable{
						startPosition: 32,
						length:        34,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotDelta, err := NewDelta(tc.giveInput, tc.giveSig)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantDelta, gotDelta)
		})
	}
}

func TestReadDelta(t *testing.T) {
	giveBuff := bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 1, 0, 0, 0, 0, 0, 0, 0, 1, 19, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 34})
	wantDelta := &Delta{
		chunks: []chunk{
			&reusable{
				startPosition: 0,
				length:        32,
			},
			&modified{
				data: []byte{19},
			},
			&reusable{
				startPosition: 32,
				length:        34,
			},
		},
	}

	gotDelta, err := ReadDelta(giveBuff)
	assert.NoError(t, err)
	assert.Equal(t, wantDelta, gotDelta)
}

func TestDeltaWrite(t *testing.T) {
	giveDelta := &Delta{
		chunks: []chunk{
			&reusable{
				startPosition: 0,
				length:        32,
			},
			&modified{
				data: []byte{19},
			},
			&reusable{
				startPosition: 32,
				length:        34,
			},
		},
	}
	wantBuff := bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 1, 0, 0, 0, 0, 0, 0, 0, 1, 19, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 34})

	gotBuff := &bytes.Buffer{}
	err := giveDelta.Write(gotBuff)
	assert.NoError(t, err)
	assert.Equal(t, wantBuff, gotBuff)
}

func TestDeltaPatch(t *testing.T) {
	giveBaseBuff := bytes.NewReader([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
		101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
		108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32})
	giveDelta := &Delta{
		chunks: []chunk{
			&reusable{
				startPosition: 0,
				length:        32,
			},
			&modified{
				data: []byte{19},
			},
			&reusable{
				startPosition: 32,
				length:        34,
			},
		},
	}
	wantBuff := bytes.NewBuffer([]byte{104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104,
		101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 19, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108,
		108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32, 104, 101, 108, 108, 111, 32})

	gotBuff := &bytes.Buffer{}
	err := giveDelta.Patch(giveBaseBuff, gotBuff)
	assert.NoError(t, err)
	assert.Equal(t, wantBuff, gotBuff)
}
