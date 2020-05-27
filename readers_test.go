package iohelper

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineNumReportingCsvReader(t *testing.T) {
	r := NewLineNumReportingCsvReader(strings.NewReader("a,b,c"))
	assert.Equal(t, 0, r.LineNum())
	record, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, record)
	assert.Equal(t, 1, r.LineNum())
	record, err = r.Read()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, record)
	assert.Equal(t, 2, r.LineNum())

	r = NewLineNumReportingCsvReader(strings.NewReader("a,b,c"))
	r.numLineFieldName = "non-existing"
	assert.PanicsWithValue(
		t,
		"unable to get 'non-existing' from csv.Reader, has csv.Reader been changed/upgraded?",
		func() {
			r.LineNum()
		})
}

func TestBytesReplacingReader(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    []byte
		search   []byte
		replace  []byte
		expected []byte
	}{
		{
			name:     "len(replace) > len(search)",
			input:    []byte{1, 2, 3, 2, 2, 3, 4, 5},
			search:   []byte{2, 3},
			replace:  []byte{4, 5, 6},
			expected: []byte{1, 4, 5, 6, 2, 4, 5, 6, 4, 5},
		},
		{
			name:     "len(replace) < len(search)",
			input:    []byte{1, 2, 3, 2, 2, 3, 4, 5, 6, 7, 8},
			search:   []byte{2, 3, 2},
			replace:  []byte{9},
			expected: []byte{1, 9, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "strip out search, no replace",
			input:    []byte{1, 2, 3, 2, 2, 3, 4, 2, 3, 2, 8},
			search:   []byte{2, 3, 2},
			replace:  []byte{},
			expected: []byte{1, 2, 3, 4, 8},
		},
		{
			name:     "len(replace) == len(search)",
			input:    []byte{1, 2, 3, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5},
			search:   []byte{5, 5},
			replace:  []byte{6, 6},
			expected: []byte{1, 2, 3, 4, 6, 6, 6, 6, 6, 6, 6, 6, 5},
		},
		{
			name:     "double quote -> single quote",
			input:    []byte(`r = NewLineNumReportingCsvReader(strings.NewReader("a,b,c"))`),
			search:   []byte(`"`),
			replace:  []byte(`'`),
			expected: []byte(`r = NewLineNumReportingCsvReader(strings.NewReader('a,b,c'))`),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := NewBytesReplacingReader(bytes.NewReader(test.input), test.search, test.replace)
			result, err := ioutil.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}

	assert.PanicsWithValue(t, "io.Reader cannot be nil", func() {
		NewBytesReplacingReader(nil, []byte{1}, []byte{2})
	})
	assert.PanicsWithValue(t, "search token cannot be nil/empty", func() {
		(&BytesReplacingReader{}).Reset(strings.NewReader("test"), nil, []byte("est"))
	})
}

func createTestInput(length int, numTarget int) []byte {
	rand.Seed(1234) // fixed rand seed to ensure bench stability
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte(rand.Intn(100) + 10) // all regular numbers >= 10
	}
	for i := 0; i < numTarget; i++ {
		for {
			index := rand.Intn(length)
			if b[index] == 7 {
				continue
			}
			b[index] = 7 // special number 7 we will search for and replace with 8.
			break
		}
	}
	return b
}

var testInput70MBLength500Targets = createTestInput(70*1024*1024, 500)
var testInput1KBLength20Targets = createTestInput(1024, 20)
var testInput50KBLength1000Targets = createTestInput(50*1024, 1000)
var testSearchFor = []byte{7}
var testReplaceWith = []byte{8}

func BenchmarkBytesReplacingReader_70MBLength_500Targets(b *testing.B) {
	r := &BytesReplacingReader{}
	for i := 0; i < b.N; i++ {
		r.Reset(bytes.NewReader(testInput70MBLength500Targets), testSearchFor, testReplaceWith)
		_, _ = ioutil.ReadAll(r)
	}
}

func BenchmarkRegularReader_70MBLength_500Targets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ioutil.ReadAll(bytes.NewReader(testInput70MBLength500Targets))
	}
}

func BenchmarkBytesReplacingReader_1KBLength_20Targets(b *testing.B) {
	r := &BytesReplacingReader{}
	for i := 0; i < b.N; i++ {
		r.Reset(bytes.NewReader(testInput1KBLength20Targets), testSearchFor, testReplaceWith)
		_, _ = ioutil.ReadAll(r)
	}
}

func BenchmarkRegularReader_1KBLength_20Targets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ioutil.ReadAll(bytes.NewReader(testInput1KBLength20Targets))
	}
}

func BenchmarkBytesReplacingReader_50KBLength_1000Targets(b *testing.B) {
	r := &BytesReplacingReader{}
	for i := 0; i < b.N; i++ {
		r.Reset(bytes.NewReader(testInput50KBLength1000Targets), testSearchFor, testReplaceWith)
		_, _ = ioutil.ReadAll(r)
	}
}

func BenchmarkRegularReader_50KBLength_1000Targets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ioutil.ReadAll(bytes.NewReader(testInput50KBLength1000Targets))
	}
}
