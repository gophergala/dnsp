package dnsp_test

import (
	"testing"
	"time"

	"github.com/gophergala/dnsp"
)

func TestInvalidOptions(t *testing.T) {
	t.Parallel()

	s, err := dnsp.NewServer(dnsp.Options{
		Poll: -time.Second, // negative poll
	})
	if err == nil {
		t.Error("expected an error, got nil")
	}

	if s != nil {
		t.Errorf("expected nil, got %+v", s)
	}
}

func TestListenAndServe(t *testing.T) {
	t.Parallel()

	t.Skip()

	s, err := dnsp.NewServer(dnsp.Options{
		Bind: ":0",
	})
	if err != nil {
		t.Errorf("%T %#v", err, err)
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		if err := s.Shutdown(); err != nil {
			t.Error(err)
		}

		done <- struct{}{}
	}()

	if err = s.ListenAndServe(); err != nil {
		t.Fatal(err)
	}

	<-done // waint until the above goroutine finishes
}
