package cbor

import (
	"encoding/base64"

	. "github.com/polydawn/go-xlate/tok"
	"github.com/polydawn/go-xlate/tok/fixtures"
)

func bcat(bss ...[]byte) []byte {
	l := 0
	for _, bs := range bss {
		l += len(bs)
	}
	rbs := make([]byte, 0, l)
	for _, bs := range bss {
		rbs = append(rbs, bs...)
	}
	return rbs
}

func b(b byte) []byte { return []byte{b} }

func deB64(s string) []byte {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bs
}

type situation byte

const (
	situationEncoding situation = 0x1
	situationDecoding situation = 0x2
)

var cborFixtures = []struct {
	title    string
	sequence fixtures.Sequence
	serial   []byte
	only     situation
}{
	// Strings.
	{"",
		fixtures.SequenceMap["flat string"],
		bcat(b(0x60+5), []byte(`value`)),
		situationEncoding | situationDecoding,
	},
	{"indefinite length string (single actual hunk)",
		fixtures.SequenceMap["flat string"],
		bcat(b(0x7f), b(0x60+5), []byte(`value`), b(0xff)),
		situationDecoding,
	},
	{"indefinite length string (multiple hunks)",
		fixtures.SequenceMap["flat string"],
		bcat(b(0x7f), b(0x60+2), []byte(`va`), b(0x60+3), []byte(`lue`), b(0xff)),
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["strings needing escape"],
		bcat(b(0x60+17), []byte("str\nbroken\ttabbed")),
		situationEncoding | situationDecoding,
	},

	// Maps.
	{"",
		fixtures.SequenceMap["empty map"],
		bcat(b(0xa0)),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		bcat(b(0xbf), b(0xff)),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["single row map"],
		bcat(b(0xa0+1),
			b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
		),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["single row map"].SansLengthInfo(),
		bcat(b(0xbf),
			b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["duo row map"],
		bcat(b(0xa0+2),
			b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
			b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
		),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		bcat(b(0xbf),
			b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
			b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},

	// Arrays.
	{"",
		fixtures.SequenceMap["empty array"],
		bcat(b(0x80)),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		bcat(b(0x9f), b(0xff)),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["single entry array"],
		bcat(b(0x80+1),
			b(0x60+5), []byte(`value`),
		),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		bcat(b(0x9f),
			b(0x60+5), []byte(`value`),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},
	{"indefinite length with nested indef string",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		bcat(b(0x9f),
			bcat(b(0x7f), b(0x60+5), []byte(`value`), b(0xff)),
			b(0xff),
		),
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["duo entry array"],
		bcat(b(0x80+2),
			b(0x60+5), []byte(`value`),
			b(0x60+2), []byte(`v2`),
		),
		situationEncoding | situationDecoding,
	},
	{"indefinite length",
		fixtures.SequenceMap["duo entry array"].SansLengthInfo(),
		bcat(b(0x9f),
			b(0x60+5), []byte(`value`),
			b(0x60+2), []byte(`v2`),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},

	// Complex / mixed / nested.
	{"all indefinite length",
		fixtures.SequenceMap["array nested in map as non-first and final entry"].SansLengthInfo(),
		bcat(b(0xbf),
			b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
			b(0x60+2), []byte(`ke`), bcat(b(0x9f),
				b(0x60+2), []byte(`oh`),
				b(0x60+4), []byte(`whee`),
				b(0x60+3), []byte(`wow`),
				b(0xff),
			),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},
	{"all indefinite length",
		fixtures.SequenceMap["array nested in map as first and non-final entry"].SansLengthInfo(),
		bcat(b(0xbf),
			b(0x60+2), []byte(`ke`), bcat(b(0x9f),
				b(0x60+2), []byte(`oh`),
				b(0x60+4), []byte(`whee`),
				b(0x60+3), []byte(`wow`),
				b(0xff),
			),
			b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["maps nested in array"],
		bcat(b(0x80+3),
			bcat(b(0xa0+1),
				b(0x60+1), []byte(`k`), b(0x60+1), []byte(`v`),
			),
			b(0x60+4), []byte(`whee`),
			bcat(b(0xa0+1),
				b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
			),
		),
		situationEncoding | situationDecoding,
	},
	{"all indefinite length",
		fixtures.SequenceMap["maps nested in array"].SansLengthInfo(),
		bcat(b(0x9f),
			bcat(b(0xbf),
				b(0x60+1), []byte(`k`), b(0x60+1), []byte(`v`),
				b(0xff),
			),
			b(0x60+4), []byte(`whee`),
			bcat(b(0xbf),
				b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
				b(0xff),
			),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["arrays in arrays in arrays"],
		bcat(b(0x80+1), b(0x80+1), b(0x80+0)),
		situationEncoding | situationDecoding,
	},
	{"all indefinite length",
		fixtures.SequenceMap["arrays in arrays in arrays"].SansLengthInfo(),
		bcat(b(0x9f), b(0x9f), b(0x9f), b(0xff), b(0xff), b(0xff)),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["maps nested in maps"],
		bcat(b(0xa0+1),
			b(0x60+1), []byte(`k`), bcat(b(0xa0+1),
				b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
			),
		),
		situationEncoding | situationDecoding,
	},
	{"all indefinite length",
		fixtures.SequenceMap["maps nested in maps"].SansLengthInfo(),
		bcat(b(0xbf),
			b(0x60+1), []byte(`k`), bcat(b(0xbf),
				b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
				b(0xff),
			),
			b(0xff),
		),
		situationEncoding | situationDecoding,
	},

	// Numbers.
	{"",
		fixtures.Sequence{"integer zero", []Token{{Type: TInt, Int: 0}}},
		deB64("AA=="),
		situationEncoding, // Impossible to decode this because cbor doens't disambiguate positive vs signed ints.
	},
	{"",
		fixtures.Sequence{"integer zero unsigned", []Token{{Type: TUint, Uint: 0}}},
		deB64("AA=="),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.Sequence{"integer one", []Token{{Type: TInt, Int: 1}}},
		deB64("AQ=="),
		situationEncoding, // Impossible to decode this because cbor doens't disambiguate positive vs signed ints.
	},
	{"",
		fixtures.Sequence{"integer one unsigned", []Token{{Type: TUint, Uint: 1}}},
		deB64("AQ=="),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.Sequence{"integer neg 1", []Token{{Type: TInt, Int: -1}}},
		deB64("IA=="),
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.Sequence{"integer neg 100", []Token{{Type: TInt, Int: -100}}},
		deB64("OGM="),
		situationEncoding | situationDecoding,
	},

	// Byte strings.
}
