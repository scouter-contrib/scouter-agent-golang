package cachemap

import (
	"strconv"
	"testing"
)

func TestCacheSet(t *testing.T) {
	cacheMap := New(10)

	//add "1" to "9"
	for i := 0; i < 10; i++ {
		cacheMap.Add(i, strconv.Itoa(i))
	}

	if cacheMap.Size() != 10 {
		t.Error("size error.")
	}

	//max size preserved
	cacheMap.Add(3, "3")
	if cacheMap.Size() != 10 {
		t.Error("size error.")
	}

	//0 exist
	if (!cacheMap.Contains(0)) {
		t.Error("contains error.")
	}

	//0 removed
	cacheMap.Add(1000, "1000")
	if cacheMap.Size() != 10 {
		t.Error("size error.")
	}

	if cacheMap.Contains(0) {
		t.Error("contains error: 0 was removed.")
	}

	//1 exist
	if !cacheMap.Contains(1) {
		t.Error("contains error: 1 should be exist.")
	}

	//1 removed
	cacheMap.Add(1001, "1001")
	if cacheMap.Size() != 10 {
		t.Error("size error.")
	}

	if cacheMap.Contains(1) {
		t.Error("contains error: 1 was removed.")
	}

	if !cacheMap.Contains(1001) {
		t.Error("contains error: 1001 should be exist.")
	}

	if cacheMap.Get(1001) != "1001" {
		t.Error("no matched value.")
	}
}

//func TestValues(t *testing.T) {
//	cacheMap := New(9000)
//
//	go func() {
//		for {
//			time.Sleep(10 * time.Nanosecond)
//			rand := rand.Intn(30000)
//			cacheMap.Add(rand, rand)
//		}
//	}()
//
//	go func() {
//		for {
//			values := cacheMap.GetValues()
//			fmt.Println(strconv.Itoa(len(values)))
//			time.Sleep(10 * time.Millisecond)
//		}
//	}()
//
//	time.Sleep(2000  * time.Millisecond)
//}
