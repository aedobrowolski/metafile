package metafile

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"testing"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func TestCodec(t *testing.T) {
	// setup the test cases
	var codec Codec = GobCodec{}

	var b bool
	var i int
	var c int32
	var f float64
	var s string
	var ai []int
	var msi map[string]int
	var xyz struct{ X, Y, Z int }

	testCases := []struct {
		desc  string
		vptr  interface{}
		value interface{}
	}{
		{"bool", &b, true},
		{"int", &i, 42},
		{"int32", &c, 'x'},
		{"float64", &f, 1.0123456789},
		{"string", &s, "hello"},
		{"[]int", &ai, []int{0, 1, 2, 3}},
		{"map[string]int", &msi, map[string]int{"one": 1, "two": 2}},
		{"struct", &xyz, struct{ X, Y, Z int }{3, 4, 5}},
	}

	// run the tests that ensure data is properly round tripped
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// execute
			b, err := codec.Encode(tC.value)
			must(err)
			must(codec.Decode(b, tC.vptr))
			// compare
			elem := reflect.ValueOf(tC.vptr).Elem().Interface()
			if !reflect.DeepEqual(elem, tC.value) {
				t.Errorf("Decode(Encode(v)) != v; got %+v [%T], expected %+v [%T]", elem, elem, tC.value, tC.value)
			}
		})
	}
}

func Test(t *testing.T) {
	testCases := []interface{}{
		struct{ X, Y, Z int }{3, 4, 5},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%+v", tC), func(t *testing.T) {
		})
	}
}

// Prove that you can write to and read from a single buffer
// Put:
// * The gob code will write to the buffer to encode the value
// * The persist code will read the buffer to persist it to the db file, using bucket/key
// Get:
// * The persist code will read the db file at bucket/key to get a []byte value
// * The gob code will decode the []byte value into a pointer

func ExampleGobCodec_Encode() {
	type P struct {
		X, Y, Z int
		Name    string
	}
	type Q struct {
		X, Y *int32
		Name string
	}

	var bb0 bytes.Buffer // Stand-in for a network connection

	// Encode (send) two consecutive values on one stream.
	var enc = gob.NewEncoder(&bb0)
	must(enc.Encode("Hello"))
	must(enc.Encode(P{3, 4, 5, "Pythagoras"}))
	// Using the same stream, encode a single consecutive value.
	enc = gob.NewEncoder(&bb0)
	must(enc.Encode(P{3, 4, 5, "Aristotle"}))

	// copy the written bytes into an io.Reader
	bb1 := bytes.NewBuffer(bb0.Bytes())
	var p P
	var q Q
	var s string

	// Decode (receive) two consecutive values.
	var dec = gob.NewDecoder(bb1)
	must(dec.Decode(&s))
	must(dec.Decode(&q))
	// Use a new decoder to decode one more consecutive value on the same stream.
	dec = gob.NewDecoder(bb1)
	must(dec.Decode(&p))

	// Confirm that the data was properly decoded
	fmt.Printf("s: %s\n", s)
	fmt.Printf("q: {X:%d Y:%d Name:%s}\n", *q.X, *q.Y, q.Name)
	fmt.Printf("p: {X:%d Y:%d Z:%d Name:%s}\n", p.X, p.Y, p.Z, p.Name)
	// Output:
	// s: Hello
	// q: {X:3 Y:4 Name:Pythagoras}
	// p: {X:3 Y:4 Z:5 Name:Aristotle}
}
