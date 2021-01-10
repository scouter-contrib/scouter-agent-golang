package cacheset

import (
	"strconv"
	"testing"
)

func TestCacheSet(t *testing.T) {
	cacheSet := New(10)

	//add "1" to "9"
	for i := 0; i < 10; i++ {
		cacheSet.Add(strconv.Itoa(i))
	}

	if cacheSet.Size() != 10 {
		t.Error("size error.")
	}

	//max size preserved
	cacheSet.Add("3")
	if cacheSet.Size() != 10 {
		t.Error("size error.")
	}

	//0 exist
	if (!cacheSet.Contains("0")) {
		t.Error("contains error.")
	}

	//0 removed
	cacheSet.Add("1000")
	if cacheSet.Size() != 10 {
		t.Error("size error.")
	}

	if cacheSet.Contains("0") {
		t.Error("contains error: 0 was removed.")
	}

	//1 exist
	if !cacheSet.Contains("1") {
		t.Error("contains error: 1 should be exist.")
	}

	//1 removed
	cacheSet.Add("1001")
	if cacheSet.Size() != 10 {
		t.Error("size error.")
	}

	if cacheSet.Contains("1") {
		t.Error("contains error: 1 was removed.")
	}

	if !cacheSet.Contains("1001") {
		t.Error("contains error: 1001 should be exist.")
	}
}
