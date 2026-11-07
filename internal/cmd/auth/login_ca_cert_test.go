package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aarondpn/redmine-cli/v2/internal/cmdutil"
	"github.com/aarondpn/redmine-cli/v2/internal/config"
)

// tlsLoginServer stubs /users/current.json over TLS with a self-signed
// certificate and returns the server and a PEM file holding that certificate.
func tlsLoginServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/current.json" {
			fmt.Fprint(w, `{"user":{"id":1,"login":"jdoe","firstname":"Jane","lastname":"Doe"}}`)
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)
	caFile := filepath.Join(t.TempDir(), "ca.crt")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	if err := os.WriteFile(caFile, pemBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	return srv, caFile
}

func loadProfile(t *testing.T, cfgPath, name string) config.Config {
	t.Helper()
	pc, err := config.LoadProfiles(cfgPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := pc.Profiles[name]
	if !ok {
		t.Fatalf("profile %q not saved: %+v", name, pc.Profiles)
	}
	return p
}

func TestLoginCACert_FirstLoginSavesCACert(t *testing.T) {
	srv, caFile := tlsLoginServer(t)
	f, cfgPath := loginFactory(t, "")
	f.ServerOverride = srv.URL
	f.APIKeyOverride = "k"

	if err := runLoginCmd(f, "--name", "w", "--ca-cert", caFile); err != nil {
		t.Fatal(err)
	}
	if p := loadProfile(t, cfgPath, "w"); p.CACert != caFile || p.APIKey != "k" {
		t.Fatalf("profile = %+v, want ca_cert %s", p, caFile)
	}
}

func TestLoginCACert_RelativePathSavedAbsolute(t *testing.T) {
	srv, caFile := tlsLoginServer(t)
	f, cfgPath := loginFactory(t, "")
	f.ServerOverride = srv.URL
	f.APIKeyOverride = "k"
	t.Chdir(filepath.Dir(caFile))

	if err := runLoginCmd(f, "--name", "w", "--ca-cert", filepath.Base(caFile)); err != nil {
		t.Fatal(err)
	}
	if p := loadProfile(t, cfgPath, "w"); p.CACert != caFile {
		t.Fatalf("ca_cert = %q, want absolute %q", p.CACert, caFile)
	}
}

func TestLoginCACert_OmittedWithoutFlag(t *testing.T) {
	srv := loginServer(t)
	f, cfgPath := loginFactory(t, "")
	f.ServerOverride = srv.URL
	f.APIKeyOverride = "k"

	if err := runLoginCmd(f, "--name", "w"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "ca_cert") {
		t.Fatalf("config contains ca_cert without --ca-cert:\n%s", data)
	}
}

func TestLoginCACert_ReloginKeepsOrClears(t *testing.T) {
	srv, caFile := tlsLoginServer(t)
	plain := loginServer(t)
	for name, tc := range map[string]struct {
		server string
		args   []string
		want   string
	}{
		"kept without flag": {server: srv.URL, args: nil, want: caFile},
		"cleared by empty":  {server: plain.URL, args: []string{"--ca-cert", ""}, want: ""},
	} {
		t.Run(name, func(t *testing.T) {
			f, cfgPath := loginFactory(t, "")
			yaml := "active_profile: w\nprofiles:\n  w:\n    server: " + srv.URL + "\n    api_key: old\n    auth_method: apikey\n    ca_cert: " + caFile + "\n"
			if err := os.WriteFile(cfgPath, []byte(yaml), 0o600); err != nil {
				t.Fatal(err)
			}
			f.ServerOverride = tc.server
			f.APIKeyOverride = "new"

			if err := runLoginCmd(f, append([]string{"--name", "w"}, tc.args...)...); err != nil {
				t.Fatal(err)
			}
			if p := loadProfile(t, cfgPath, "w"); p.APIKey != "new" || p.CACert != tc.want {
				t.Fatalf("profile = %+v, want ca_cert %q", p, tc.want)
			}
		})
	}
}

func TestLoginCACert_UnknownCAHintAndNoConfig(t *testing.T) {
	srv, _ := tlsLoginServer(t)
	f, cfgPath := loginFactory(t, "")
	f.ServerOverride = srv.URL
	f.APIKeyOverride = "k"

	err := runLoginCmd(f, "--name", "w")
	if hint := cmdutil.ErrorHint(err); hint != loginUnknownCAHint {
		t.Fatalf("err = %v, hint = %q, want the unknown CA hint", err, hint)
	}
	if _, statErr := os.Stat(cfgPath); statErr == nil {
		t.Fatal("config written despite failed login")
	}
}

// otherCAFile writes a freshly generated self-signed CA certificate, which
// did not sign the httptest server certificate.
func otherCAFile(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Other CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "other-ca.crt")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoginCACert_WrongCAHint(t *testing.T) {
	srv, _ := tlsLoginServer(t)
	otherCA := otherCAFile(t)
	f, _ := loginFactory(t, "")
	f.ServerOverride = srv.URL
	f.APIKeyOverride = "k"

	err := runLoginCmd(f, "--name", "w", "--ca-cert", otherCA)
	if hint := cmdutil.ErrorHint(err); hint != fmt.Sprintf(loginWrongCAHint, otherCA) {
		t.Fatalf("err = %v, hint = %q, want the wrong CA hint", err, hint)
	}
}

func TestLoginCACert_MissingFileFailsBeforeConnecting(t *testing.T) {
	f, cfgPath := loginFactory(t, "")
	f.ServerOverride = "https://redmine.invalid"
	f.APIKeyOverride = "k"

	err := runLoginCmd(f, "--name", "w", "--ca-cert", filepath.Join(t.TempDir(), "missing.crt"))
	if err == nil || !strings.Contains(err.Error(), "--ca-cert") {
		t.Fatalf("err = %v, want --ca-cert error", err)
	}
	if _, statErr := os.Stat(cfgPath); statErr == nil {
		t.Fatal("config written despite failed login")
	}
}
