package chopener

import (
	"reflect"
	"sync"
	"testing"
)

func TestOpen_Buffered(t *testing.T) {
	var (
		arr1     = [...]int{1, 2, 3, 4, 5}
		arr2     = [...]int{6, 7, 8, 9, 10}
		expected = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		got      = make([]int, 0, len(expected))
		size     = len(arr1)
		i        int
	)

	ch := make(chan int, size)

	for _, el := range &arr1 {
		ch <- el
	}

	for val := range ch {
		got = append(got, val)

		i++
		if i == size {
			close(ch)
		}
	}

	Open(&ch)

	for _, el := range &arr2 {
		ch <- el
	}

	i = 0
	for val := range ch {
		got = append(got, val)

		i++
		if i == size {
			close(ch)
		}
	}

	if !reflect.DeepEqual(expected[:], got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestOpen_Unbuffered(t *testing.T) {
	t.Parallel()

	ch := make(chan string)

	var (
		expected = []string{"Hello", "world", "foo", "bar"}
		got      = make([]string, 0, len(expected))
		mu       sync.Mutex
	)

	mu.Lock()
	go func() {
		ch <- expected[0]
		ch <- expected[1]
		close(ch)

		Open(&ch)

		mu.Unlock()

		ch <- expected[2]
		ch <- expected[3]
		close(ch)
	}()

	for val := range ch {
		got = append(got, val)
	}

	mu.Lock()

	for val := range ch {
		got = append(got, val)
	}

	mu.Unlock()

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
