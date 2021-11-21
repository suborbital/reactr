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
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("https allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("IP allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://100.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})
}

func TestAllowedDomains(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowedDomains = []string{"example.com", "another.com", "*.hello.com", "tomorrow.*", "100.*.12.13"}

	t.Run("example.com:8080 allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com:8080", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("example.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("another.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://another.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("double wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://goodmorning.goodbye.hello.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("end wildcard allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://tomorrow.eu", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
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
		req, _ := http.NewRequest(http.MethodGet, "http://100.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://101.12.13.14", nil)

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
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("athird.com allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://athird.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("IP allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://100.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})
}

func TestBlockedWithCNAME(t *testing.T) {
	rules := defaultHTTPRules()
	rules.BlockedDomains = []string{"hosting.gitbook.io"}

	t.Run("Resolved CNAME blocked", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://atmo.suborbital.dev", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})
}

func TestDisallowedIPs(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowIPs = false

	t.Run("IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://100.11.12.13", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Localhost IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Private IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://192.168.0.11", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Loopback IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[::1]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Localhost IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[fe80::1]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Private IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[fd00::2f00]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Public IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[2604:a880:cad:d0::dff:7001]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("domain allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://friday.com", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})

	t.Run("localhost allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
		}
	})
}

func TestDisallowedLocal(t *testing.T) {
	rules := defaultHTTPRules()
	rules.AllowPrivate = false

	t.Run("Loopback IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[::1]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Localhost IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[fe80::1]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Private IPv6 disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://[fd00::2f00]", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Resolves to Private disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://local.suborbital.network", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Resolves to Private (with port) disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://local.suborbital.network:8081", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Private disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Localhost IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Private IP disallowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://192.168.0.11", nil)

		if err := rules.requestIsAllowed(req); err == nil {
			t.Error("error did not occur, should have")
		}
	})

	t.Run("Resolves to public allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://suborbital.dev", nil)

		if err := rules.requestIsAllowed(req); err != nil {
			t.Error("error occurred, should not have:", err)
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
			t.Error("error occurred, should not have:", err)
		}
	})
}
