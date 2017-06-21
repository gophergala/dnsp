package dnsp_test

import (
	"testing"
	"time"

	"github.com/gophergala/dnsp"
)

func TestValidate(t *testing.T) {
	var testCases = []struct {
		inputOptions dnsp.Options
		wantErr      bool
	}{
		{
			// invalid Net
			dnsp.Options{
				Net:       "something invalid",
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// valid bind
			dnsp.Options{
				Bind:      "example.com:dns",
				BlockedIP: "192.168.1.117",
			},
			false,
		},
		{
			// invalid resolve
			dnsp.Options{
				Resolve: []string{
					"something.com",
				},
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// valid resolve
			dnsp.Options{
				Resolve: []string{
					"0.0.0.0:53",
				},
				BlockedIP: "192.168.1.117",
			},
			false,
		},
		{
			// Poll too short
			dnsp.Options{
				Poll:      time.Millisecond * 900,
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// whitelist and blacklist together
			dnsp.Options{
				Whitelist: "wikipedia.com",
				Blacklist: "badsite.com",
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// invalid whitelist
			dnsp.Options{
				Whitelist: "somethinginvalid",
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// invalid blacklist
			dnsp.Options{
				Blacklist: "somethinginvalid",
				BlockedIP: "192.168.1.117",
			},
			true,
		},
		{
			// valid blockedip
			dnsp.Options{
				BlockedIP: "192.168.1.117",
			},
			false,
		},
		{
			// invalid blockedip
			dnsp.Options{
				BlockedIP: "somethinginvalid",
			},
			true,
		},
	}

	for _, tt := range testCases {
		err := tt.inputOptions.Validate()

		// if we expect an error and there isn't one
		if tt.wantErr && err == nil {
			t.Errorf("expected an error, but err is nil")
		}
		// if we don't expect an error and there is one
		if !tt.wantErr && err != nil {
			t.Errorf("expected error to be nil, but err is: %q", err)
		}
	}
}

func TestPathOrURL(t *testing.T) {
	var testCases = []struct {
		inputPath  string
		wantString string
		wantErr    bool
	}{
		{
			"//userinfo@host/path.com",
			"//userinfo@host/path.com",
			false,
		},
	}

	for _, tt := range testCases {
		str, err := dnsp.PathOrURL(tt.inputPath)

		if str != tt.wantString {
			t.Errorf("wanted %q but got %q", tt.wantString, str)
		}

		// if we expect an error and there isn't one
		if tt.wantErr && err == nil {
			t.Errorf("expected an error, but err is nil")
		}
		// if we don't expect an error and there is one
		if !tt.wantErr && err != nil {
			t.Errorf("expected error to be nil, but err is: %q", err)
		}
	}
}
