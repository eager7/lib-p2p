// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates the trie for idna operations. The Unicode casing
// algorithm requires the lookup of various properties and mappings for each
// rune. The table generated by this generator combines several of the most
// frequently used of these into a single trie so that they can be accessed
// with a single lookup.
package main

import (
	"fmt"
	"io"
	"log"
	"unicode"
	"unicode/utf8"

	"gx/ipfs/QmVcxhXDbXjNoAdmYBWbY1eU67kQ8eZUHjG4mAYZUtZZu3/go-text/internal/gen"
	"gx/ipfs/QmVcxhXDbXjNoAdmYBWbY1eU67kQ8eZUHjG4mAYZUtZZu3/go-text/internal/triegen"
	"gx/ipfs/QmVcxhXDbXjNoAdmYBWbY1eU67kQ8eZUHjG4mAYZUtZZu3/go-text/internal/ucd"
	"gx/ipfs/QmVcxhXDbXjNoAdmYBWbY1eU67kQ8eZUHjG4mAYZUtZZu3/go-text/unicode/bidi"
)

func main() {
	gen.Init()
	genTables()
	gen.Repackage("gen_trieval.go", "trieval.go", "idna")
	gen.Repackage("gen_common.go", "common_test.go", "idna")
}

var runes = map[rune]info{}

func genTables() {
	t := triegen.NewTrie("idna")

	ucd.Parse(gen.OpenUCDFile("DerivedNormalizationProps.txt"), func(p *ucd.Parser) {
		r := p.Rune(0)
		if p.String(1) == "NFC_QC" { // p.String(2) is "N" or "M"
			runes[r] = mayNeedNorm
		}
	})
	ucd.Parse(gen.OpenUCDFile("UnicodeData.txt"), func(p *ucd.Parser) {
		r := p.Rune(0)

		const cccVirama = 9
		if p.Int(ucd.CanonicalCombiningClass) == cccVirama {
			runes[p.Rune(0)] = viramaModifier
		}
		switch {
		case unicode.In(r, unicode.Mark):
			runes[r] |= modifier | mayNeedNorm
		}
		// TODO: by using UnicodeData.txt we don't mark undefined codepoints
		// that are earmarked as RTL properly. However, an undefined cp will
		// always fail, so there is no need to store this info.
		switch p, _ := bidi.LookupRune(r); p.Class() {
		case bidi.R, bidi.AL, bidi.AN:
			if x := runes[r]; x != 0 && x != mayNeedNorm {
				log.Fatalf("%U: rune both modifier and RTL letter/number", r)
			}
			runes[r] = rtl
		}
	})

	ucd.Parse(gen.OpenUCDFile("extracted/DerivedJoiningType.txt"), func(p *ucd.Parser) {
		switch v := p.String(1); v {
		case "L", "D", "T", "R":
			runes[p.Rune(0)] |= joinType[v] << joinShift
		}
	})

	ucd.Parse(gen.OpenUnicodeFile("idna", "", "IdnaMappingTable.txt"), func(p *ucd.Parser) {
		r := p.Rune(0)

		// The mappings table explicitly defines surrogates as invalid.
		if !utf8.ValidRune(r) {
			return
		}

		cat := catFromEntry(p)
		isMapped := cat == mapped || cat == disallowedSTD3Mapped || cat == deviation
		if !isMapped {
			// Only include additional category information for non-mapped
			// runes. The additional information is only used after mapping and
			// the bits would clash with mapping information.
			// TODO: it would be possible to inline this data and avoid
			// additional lookups. This is quite tedious, though, so let's first
			// see if we need this.
			cat |= category(runes[r])
		}

		s := string(p.Runes(2))
		if s != "" && !isMapped {
			log.Fatalf("%U: Mapping with non-mapping category %d", r, cat)
		}
		t.Insert(r, uint64(makeEntry(r, s))+uint64(cat))
	})

	w := gen.NewCodeWriter()
	defer w.WriteVersionedGoFile("tables.go", "idna")

	gen.WriteUnicodeVersion(w)

	w.WriteVar("mappings", string(mappings))
	w.WriteVar("xorData", string(xorData))

	sz, err := t.Gen(w, triegen.Compact(&normCompacter{}))
	if err != nil {
		log.Fatal(err)
	}
	w.Size += sz
}

var (
	// mappings contains replacement strings for mapped runes, each prefixed
	// with a byte containing the length of the following string.
	mappings = []byte{}
	mapCache = map[string]int{}

	// xorData is like mappings, except that it contains XOR data.
	// We split these two tables so that we don't get an overflow.
	xorData  = []byte{}
	xorCache = map[string]int{}
)

// makeEntry creates a trie entry.
func makeEntry(r rune, mapped string) info {
	orig := string(r)

	if len(orig) != len(mapped) {
		// Store the mapped value as is in the mappings table.
		index := len(mappings)
		if x, ok := mapCache[mapped]; ok {
			index = x
		} else {
			mapCache[mapped] = index
			mappings = append(mappings, byte(len(mapped)))
			mappings = append(mappings, mapped...)
		}
		return info(index) << indexShift
	}

	// Create per-byte XOR mask.
	var b []byte
	for i := 0; i < len(orig); i++ {
		b = append(b, orig[i]^mapped[i])
	}

	// Remove leading 0 bytes, but keep at least one byte.
	for ; len(b) > 1 && b[0] == 0; b = b[1:] {
	}

	if len(b) == 1 {
		return xorBit | inlineXOR | info(b[0])<<indexShift
	}
	mapped = string(b)

	// Store the mapped value as is in the mappings table.
	index := len(xorData)
	if x, ok := xorCache[mapped]; ok {
		index = x
	} else {
		xorCache[mapped] = index
		xorData = append(xorData, byte(len(mapped)))
		xorData = append(xorData, mapped...)
	}
	return xorBit | info(index)<<indexShift
}

// The following code implements a triegen.Compacter that was originally
// designed for normalization. The IDNA table has some similarities with the
// norm table. Using this compacter, together with the XOR pattern approach,
// reduces the table size by roughly 100K. It can probably be compressed further
// by also including elements of the compacter used by cases, but for now it is
// good enough.

const maxSparseEntries = 16

type normCompacter struct {
	sparseBlocks [][]uint64
	sparseOffset []uint16
	sparseCount  int
}

func mostFrequentStride(a []uint64) int {
	counts := make(map[int]int)
	var v int
	for _, x := range a {
		if stride := int(x) - v; v != 0 && stride >= 0 {
			counts[stride]++
		}
		v = int(x)
	}
	var maxs, maxc int
	for stride, cnt := range counts {
		if cnt > maxc || (cnt == maxc && stride < maxs) {
			maxs, maxc = stride, cnt
		}
	}
	return maxs
}

func countSparseEntries(a []uint64) int {
	stride := mostFrequentStride(a)
	var v, count int
	for _, tv := range a {
		if int(tv)-v != stride {
			if tv != 0 {
				count++
			}
		}
		v = int(tv)
	}
	return count
}

func (c *normCompacter) Size(v []uint64) (sz int, ok bool) {
	if n := countSparseEntries(v); n <= maxSparseEntries {
		return (n+1)*4 + 2, true
	}
	return 0, false
}

func (c *normCompacter) Store(v []uint64) uint32 {
	h := uint32(len(c.sparseOffset))
	c.sparseBlocks = append(c.sparseBlocks, v)
	c.sparseOffset = append(c.sparseOffset, uint16(c.sparseCount))
	c.sparseCount += countSparseEntries(v) + 1
	return h
}

func (c *normCompacter) Handler() string {
	return "idnaSparse.lookup"
}

func (c *normCompacter) Print(w io.Writer) (retErr error) {
	p := func(f string, x ...interface{}) {
		if _, err := fmt.Fprintf(w, f, x...); retErr == nil && err != nil {
			retErr = err
		}
	}

	ls := len(c.sparseBlocks)
	p("// idnaSparseOffset: %d entries, %d bytes\n", ls, ls*2)
	p("var idnaSparseOffset = %#v\n\n", c.sparseOffset)

	ns := c.sparseCount
	p("// idnaSparseValues: %d entries, %d bytes\n", ns, ns*4)
	p("var idnaSparseValues = [%d]valueRange {", ns)
	for i, b := range c.sparseBlocks {
		p("\n// Block %#x, offset %#x", i, c.sparseOffset[i])
		var v int
		stride := mostFrequentStride(b)
		n := countSparseEntries(b)
		p("\n{value:%#04x,lo:%#02x},", stride, uint8(n))
		for i, nv := range b {
			if int(nv)-v != stride {
				if v != 0 {
					p(",hi:%#02x},", 0x80+i-1)
				}
				if nv != 0 {
					p("\n{value:%#04x,lo:%#02x", nv, 0x80+i)
				}
			}
			v = int(nv)
		}
		if v != 0 {
			p(",hi:%#02x},", 0x80+len(b)-1)
		}
	}
	p("\n}\n\n")
	return
}
