package rcap

import (
	"net/http"
	"testing"
)

func TestDefaultRules(t *testing.T) {
	rules := defaultHTTPRules()

	t.Run("http allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("https allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("IP allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://10.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})
}

func TestAllowedDomains(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowedDomains = []string{"example.com", "another.com", "*.hello.com", "tomorrow.*", "10.*.12.13"}

	t.Run("example.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("another.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://another.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("double wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodmorning.goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("end wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://tomorrow.eu", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("double end wildcard disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://tomorrow.co.uk", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("athird.com disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://athird.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("wildcard IP allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://10.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://11.12.13.14", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})
}

func TestBlockedDomains(t *testing.T) {
	rules := defaultHTTPRules()
	rules.BlockedDomains = []string{"example.com", "another.com", "*.hello.com", "tomorrow.*"}

	t.Run("example.com disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("another.com disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://another.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("wildcard disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("double wildcard disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodnight.goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("end wildcard disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://tomorrow.eu", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("double end wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://tomorrow.co.uk", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("athird.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://athird.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})

	t.Run("IP allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://10.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})
}

func TestDisallowedIPs(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowIPs = false

	t.Run("IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://10.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("domain allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://friday.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})
}

func TestDisallowHTTP(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowHTTP = false

	t.Run("HTTP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("HTTPS allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://friday.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have")
		}
	})
}
