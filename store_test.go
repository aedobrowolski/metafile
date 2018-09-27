package metafile

import (
	"strings"
	"testing"
)

func Test_Get_Basic(t *testing.T) {
	store := newStore("", nil)
	put := "bar"
	if err := store.Put("", "foo", put); err != nil {
		t.Errorf("Error calling Put: %s", err)
	}
	var got string
	if ok, err := store.Get("", "foo", &got); err != nil {
		t.Errorf("Error calling Get: %s", err)
	} else if ok == false {
		t.Errorf("Error calling Get: %s", "ok is false")
	}
	if got != put {
		t.Errorf("Result of Get: `%s`, Expecting Put: `%s`", got, put)
	}
}

func Test_Get_MissingKey(t *testing.T) {
	store := newStore("", nil)
	put := "dummy"
	if err := store.Put("", "foo", put); err != nil {
		t.Errorf("Error calling Put: %s", err)
	}
	got := put
	if ok, err := store.Get("", "missing", &got); err != nil {
		t.Errorf("Error calling Get: %s", err)
	} else if ok != false {
		t.Errorf("Error calling Get on missing key: %s", "ok is true")
	}
	if got != put {
		t.Errorf("Result of Get: `%s`, expecting unchanged value", got)
	}
}

// The following variations all represent the same bucket "x"
// This test also checks the Put() and Get() of integer values.
func Test_Get_BucketVariations(t *testing.T) {
	testCases := []string{"x", "x/", "/x", "./x", "//x", "x//", "/./x", "y/../x"}
	store := newStore("", nil)
	for i, path := range testCases {
		t.Run(path, func(t *testing.T) {
			store.Put("x", "key", i)
			var j int
			ok, err := store.Get(path, "key", &j)
			if err != nil {
				t.Error("Get: ", err)
			} else if !ok {
				t.Errorf("Get: could not find key `key` in path `%s`", path)
			} else if i != j {
				t.Errorf("Result of Get: `%d`, Expecting Put: `%d`", j, i)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	path := Clean("X")
	store := newStore("", nil)
	keys := []string{"a", "b", "c", "d", "e"}
	// store numbers in 5 keys
	for i, k := range keys {
		store.Put(path, k, i)
	}
	if got, exp := len(store.meta[path]), len(keys); got != exp {
		t.Errorf("Put: wrong number of keys found (%d), expecting %d", got, exp)
	}
	// delete all the keys
	for _, k := range keys {
		store.Delete(path, k)
	}
	if got, exp := len(store.meta[path]), 0; got != exp {
		t.Errorf("Put: wrong number of keys found (%d), expecting %d", got, exp)
	}
}

func Test_Empty(t *testing.T) {
	path1 := Clean("X")
	path2 := Clean("Y")
	path3 := Clean("X/Z") // nested path

	store := newStore("", nil)
	keys := []string{"a", "b", "c", "d", "e"}
	// store numbers in 5 keys
	for i, k := range keys {
		store.Put(path1, k, i)
		store.Put(path2, k, i)
		store.Put(path3, k, i)
	}
	if got, exp := len(store.meta[path1]), len(keys); got != exp {
		t.Errorf("Put: wrong number of keys found (%d), expecting %d", got, exp)
	}
	// delete all the keys in X and X/Z by calling Empty(X)
	store.Empty(path1)
	if got, exp := len(store.meta[path1]), 0; got != exp {
		t.Errorf("Put: wrong number of keys in path '%s' found (%d), expecting %d", path1, got, exp)
	}
	if got, exp := len(store.meta[path2]), len(keys); got != exp {
		t.Errorf("Put: wrong number of keys in path '%s' found (%d), expecting %d", path2, got, exp)
	}
	if got, exp := len(store.meta[path3]), 0; got != exp {
		t.Errorf("Put: wrong number of keys in path '%s' found (%d), expecting %d", path3, got, exp)
	}
}

func Test_Move(t *testing.T) {
	path1 := Clean("X")
	path2 := Clean("Y")

	store := newStore("", nil)
	keys := []string{"a", "b", "c", "d", "e"}
	// store numbers in 5 keys
	for i, k := range keys {
		store.Put(path1, k, i)
	}
	if got, exp := len(store.meta[path1]), len(keys); got != exp {
		t.Errorf("Put: wrong number of keys found (%d), expecting %d", got, exp)
	}
	store.Put(path2, "dummy", 42) // should be wiped by move

	// move all the keys in X to Y by calling Move(X, Y)
	store.Move(path1, path2)
	if got, exp := len(store.meta[path1]), 0; got != exp {
		t.Errorf("Put: wrong number of keys in path '%s' found (%d), expecting %d", path1, got, exp)
	}
	if got, exp := len(store.meta[path2]), len(keys); got != exp {
		t.Errorf("Put: wrong number of keys in path '%s' found (%d), expecting %d", path2, got, exp)
	}
}

func Test_Put_valid(t *testing.T) {
	path1 := Clean("X")

	falseStore := newStore("", func(s string) bool { return false })
	if err := falseStore.Put(path1, "key", 0); err == nil {
		t.Error("Put: invalid path, expecting error but got nil.")
	}

	trueStore := newStore("", func(s string) bool { return true })
	if err := trueStore.Put(path1, "key", 0); err != nil {
		t.Errorf("Put: valid path, expecting no error but got %s.", err)
	}

	xStore := newStore("", func(s string) bool { return strings.HasPrefix(s, "x") })
	if err := xStore.Put("xpath", "key", 0); err != nil {
		t.Errorf("Put: valid x path, expecting no error but got %s.", err)
	}
	if err := xStore.Put("ypath", "key", 0); err == nil {
		t.Errorf("Put: invalid y path, expecting error but got %s.", err)
	}
}
