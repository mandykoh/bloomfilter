package bloomfilter

import (
	"sync"
	"testing"
)

func TestFilter16(t *testing.T) {

	t.Run("Add()", func(t *testing.T) {

		t.Run("increments and returns the count for the value", func(t *testing.T) {
			f := New16Bit(16, 0.1)

			if expected, actual := uint(0), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}

			count := f.Add("hello")

			if expected, actual := uint(1), count; expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(1), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}

			count = f.Add("hello")

			if expected, actual := uint(2), count; expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(2), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
		})

		t.Run("stops incrementing at 65535", func(t *testing.T) {
			f := New16Bit(16, 0.1).(*filter16)

			for i := 0; i < len(f.counters); i++ {
				f.counters[i] = 65535
			}

			f.Add("hello")

			if expected, actual := uint(65535), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
		})

		t.Run("is thread safe", func(t *testing.T) {
			f := New16Bit(16, 0.1)

			start := sync.WaitGroup{}
			start.Add(1)

			wait := sync.WaitGroup{}

			for i := 0; i < 250; i++ {
				wait.Add(1)

				go func() {
					start.Wait()
					f.Add("hello")
					wait.Done()
				}()
			}

			start.Done()
			wait.Wait()

			if expected, actual := uint(250), f.Count("hello"); expected != actual {
				t.Errorf("Expected a total count of %d but got %d; possible concurrency error", expected, actual)
			}
		})
	})

	t.Run("Clear()", func(t *testing.T) {

		t.Run("zeros out all counters", func(t *testing.T) {
			f := New16Bit(16, 0.1).(*filter16)

			for i := 0; i < len(f.counters); i++ {
				f.counters[i] = 65535
			}

			f.Clear()

			for i := 0; i < len(f.counters); i++ {
				if expected, actual := uint16(0), f.counters[i]; expected != actual {
					t.Errorf("Expected counter %d to be %d but got %d", i, expected, actual)
				}
			}
		})
	})

	t.Run("Remove()", func(t *testing.T) {

		t.Run("decrements and returns the count for the value", func(t *testing.T) {
			f := New16Bit(16, 0.1)

			f.Add("hello")
			f.Add("hello")
			f.Add("hello")

			if expected, actual := uint(3), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}

			count := f.Remove("hello")

			if expected, actual := uint(2), count; expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(2), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}

			count = f.Remove("hello")

			if expected, actual := uint(1), count; expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(1), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if expected, actual := uint(0), f.Count("hello 2"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
		})

		t.Run("stops decrementing at zero", func(t *testing.T) {
			f := New16Bit(16, 0.1)

			f.Remove("hello")

			if expected, actual := uint(0), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
		})

		t.Run("is thread safe", func(t *testing.T) {
			f := New16Bit(16, 0.1).(*filter16)

			for i := 0; i < len(f.counters); i++ {
				f.counters[i] = 500
			}

			start := sync.WaitGroup{}
			start.Add(1)

			wait := sync.WaitGroup{}

			for i := 0; i < 250; i++ {
				wait.Add(1)

				go func() {
					start.Wait()
					f.Remove("hello")
					wait.Done()
				}()
			}

			start.Done()
			wait.Wait()

			if expected, actual := uint(250), f.Count("hello"); expected != actual {
				t.Errorf("Expected a total count of %d but got %d; possible concurrency error", expected, actual)
			}
		})
	})

	t.Run("Reset()", func(t *testing.T) {

		t.Run("zeros out counters for one value", func(t *testing.T) {
			f := New16Bit(16, 0.1)

			f.Add("hello")
			f.Add("hello")
			f.Add("hello 2")

			f.Reset("hello")

			if expected, actual := uint(0), f.Count("hello"); expected != actual {
				t.Errorf("Expected count of %d but got %d", expected, actual)
			}
			if actual := f.Count("hello 2"); actual == 0 {
				t.Errorf("Expected non-zero count but got %d", actual)
			}
		})
	})

	t.Run("New16Bit()", func(t *testing.T) {

		t.Run("initialises the number of hashes and the counters", func(t *testing.T) {
			f := New16Bit(32, 0.5).(*filter16)

			if expected, actual := uint(1), f.hashCount; expected != actual {
				t.Errorf("Expected %d hashes but got %d", expected, actual)
			}
			if expected, actual := 46, len(f.counters); expected != actual {
				t.Errorf("Expected %d counters but got %d", expected, actual)
			}

			f = New16Bit(16, 0.2).(*filter16)

			if expected, actual := uint(2), f.hashCount; expected != actual {
				t.Errorf("Expected %d hashes but got %d", expected, actual)
			}
			if expected, actual := 53, len(f.counters); expected != actual {
				t.Errorf("Expected %d counters but got %d", expected, actual)
			}

			f = New16Bit(8, 0.01).(*filter16)

			if expected, actual := uint(6), f.hashCount; expected != actual {
				t.Errorf("Expected %d hashes but got %d", expected, actual)
			}
			if expected, actual := 76, len(f.counters); expected != actual {
				t.Errorf("Expected %d counters but got %d", expected, actual)
			}
		})
	})
}
