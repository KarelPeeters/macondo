package strategy

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/domino14/macondo/alphabet"
)

// type OriginalCHD struct { // All are LittleEndian.
//   rl uint32
//   r [rl]uint64
//   il uint32
//   indices [il]uint16
//   el uint32
//   [el]struct {
//     kl uint32
//     vl uint32
//     keys [kl]byte
//     values [vl]byte
//   }
// }

// Proposed struct { // All are LittleEndian.
//   rl, r, il, indies, el unchanged
//   values [el]float32
//   maxLen uint32
//   keys [el][maxLen]byte (keys are right-padded with 0xff)
// }

type OldLeaves struct {
	r           []uint64
	indices     []uint16
	leaveFloats []float32
	maxLength   uint32
	buf         []byte
}

func (olv *OldLeaves) LeaveValue(leave alphabet.MachineWord) float64 {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic; leave was %v\n", leave.UserVisible(alphabet.EnglishAlphabet()))
			// Panic anyway; the recover was just to figure out which leave did it.
			panic("panicking anyway")
		}
	}()
	ll := len(leave)
	if ll == 0 {
		return float64(0)
	}
	for i := 1; i < ll; i++ {
		for j := i; j > 0 && leave[j-1] > leave[j]; j-- {
			leave[j-1], leave[j] = leave[j], leave[j-1]
		}
	}
	if len(leave) <= int(olv.maxLength) {
		// TODO: replace this comment with shrug emoji.
		h := uint64(14695981039346656037)
		for _, c := range leave {
			h ^= uint64(c)
			h *= 1099511628211
		}
		h ^= olv.r[0]
		ri := olv.indices[h%uint64(len(olv.indices))]

		if ri < uint16(len(olv.r)) {

			h = (h ^ olv.r[ri]) % uint64(len(olv.leaveFloats))
			bufp := int(uint64(olv.maxLength) * h)
			if string(olv.buf[bufp:bufp+len(leave)]) == string(leave) &&
				(len(leave) >= int(olv.maxLength) || olv.buf[bufp+len(leave)] == 0xff) {
				return float64(olv.leaveFloats[h])
			}
		}
	}
	// Only will happen if we have a pass. Passes are very rare and
	// we should ignore this a bit since there will be a negative
	// adjustment already from the fact that we're scoring 0.
	return float64(0)
}

// Not very useful, but it's used for logging.
func (olv *OldLeaves) Len() int {
	return len(olv.leaveFloats)
}

func ReadOldLeaves(file io.Reader) (*OldLeaves, error) {
	var rl uint32
	if err := binary.Read(file, binary.LittleEndian, &rl); err != nil {
		return nil, err
	}
	r := make([]uint64, rl)
	if err := binary.Read(file, binary.LittleEndian, &r); err != nil {
		return nil, err
	}
	var il uint32
	if err := binary.Read(file, binary.LittleEndian, &il); err != nil {
		return nil, err
	}
	indices := make([]uint16, il)
	if err := binary.Read(file, binary.LittleEndian, &indices); err != nil {
		return nil, err
	}
	var lenLeaves uint32
	if err := binary.Read(file, binary.LittleEndian, &lenLeaves); err != nil {
		return nil, err
	}
	leaveFloats := make([]float32, lenLeaves)
	if err := binary.Read(file, binary.LittleEndian, &leaveFloats); err != nil {
		return nil, err
	}
	var maxLength uint32
	if err := binary.Read(file, binary.LittleEndian, &maxLength); err != nil {
		return nil, err
	}
	buf := make([]byte, maxLength*lenLeaves)
	if err := binary.Read(file, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}

	return &OldLeaves{
		r:           r,
		indices:     indices,
		leaveFloats: leaveFloats,
		maxLength:   maxLength,
		buf:         buf,
	}, nil
}
