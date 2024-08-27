package converter

import (
	"fmt"
	"os"
	"testing"
)

func TestResolved(t *testing.T) {
	tests := []struct {
		originalString string
		expectString   string
		envs           map[string]string // for render <% ENV %>
	}{
		{
			originalString: "test-ubuntu-20.04-linux",
			expectString:   "test-ubuntu-20.04-linux",
		},
		{
			originalString: "<%      ENV    %>-linux",
			expectString:   "test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			originalString: "secret/foo/<%      ENV    %>-linux",
			expectString:   "secret/foo/test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			originalString: "<% ENV    %>-linux",
			expectString:   "test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			originalString: "<% ENV %>-<% DIST %>-<% VER %>-linux",
			expectString:   "test-ubuntu-20.04-linux",
			envs: map[string]string{
				"ENV":  "test",
				"DIST": "ubuntu",
				"VER":  "20.04",
			},
		},
		{
			originalString: "<%      ENV    -linux",
			expectString:   "<%      ENV    -linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			originalString: "<%      NOTSETVAR  %>-linux",
			expectString:   "-linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.originalString, func(t *testing.T) {
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("os.Setenv() returned an unexpected error: got: %v", err)
				}
			}
			out := resolved(tt.originalString)
			if out != tt.expectString {
				t.Errorf("resolved() returned an unexpected string: got: %s, want: %s", out, tt.expectString)
			}
		})
	}
}

func TestGetVaultSecretKey(t *testing.T) {
	tests := []struct {
		desc       string
		secretPath string
		expectKey  string
		err        error
	}{
		{
			desc:       "invalid secret path",
			secretPath: "secret/foo/bar",
			expectKey:  "",
			err:        fmt.Errorf(illegalVaultPath, "secret/foo/bar"),
		},
		{
			desc:       "path",
			secretPath: "secret/data/bar",
			expectKey:  "bar",
			err:        nil,
		},
		{
			desc:       "path",
			secretPath: "secret/data/bar/foo",
			expectKey:  "bar/foo",
			err:        nil,
		},
		{
			desc:       "path",
			secretPath: "kv-app/dep1/data/bar/foo",
			expectKey:  "bar/foo",
			err:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			out, err := getVaultSecretKey(tt.secretPath)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("test case name %s", tt.desc)
					t.Errorf("getVaultSecretKey() returned an unexpected error: got: %v, want: %v", err, tt.err)
				}
			} else {
				if out != tt.expectKey {
					t.Errorf("test case name %s", tt.desc)
					t.Errorf("getVaultSecretKey() returned an unexpected string: got: %s, want: %s", out, tt.expectKey)
				}
			}
		})
	}
}

func TestResolveAngleBrackets(t *testing.T) {
	tests := []struct {
		originalString string
		expectString   string
		err            error
	}{
		{
			originalString: "test-ubuntu-20.04-linux",
			expectString:   "test-ubuntu-20.04-linux",
		},
		{
			originalString: "<A>-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			originalString: "<   A   >-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			originalString: "<A    >-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			originalString: "<   A>-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			originalString: "sn0rt-<A>-linux",
			expectString:   "sn0rt-{{ .A }}-linux",
		},
		{
			originalString: "sn0rt-<A>-<B>",
			expectString:   "sn0rt-{{ .A }}-{{ .B }}",
		},
		{
			originalString: "sn0rt-<A",
			expectString:   "sn0rt-<A",
			err:            fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `syntax error: unclosed '<'`),
		},
		{
			originalString: "<A>-<B> <C>",
			expectString:   "{{ .A }}-{{ .B }} {{ .C }}",
		},
		{
			originalString: "sn0rt-<A>-<B> <C>",
			expectString:   "sn0rt-{{ .A }}-{{ .B }} {{ .C }}",
		},
		{
			originalString: "<MYSQL_PASSWD>",
			expectString:   "{{ .MYSQL_PASSWD }}",
		},
		{
			originalString: "password = <MYSQL_PASSWD>",
			expectString:   "password = {{ .MYSQL_PASSWD }}",
		},
		{
			originalString: `
[client]
host = example.com
user = < USER >
password = <MYSQL_PASSWD>
port = 4000`,
			expectString: `
[client]
host = example.com
user = {{ .USER }}
password = {{ .MYSQL_PASSWD }}
port = 4000`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.originalString, func(t *testing.T) {
			out, err := resolveAngleBrackets(tt.originalString)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("resolveAngleBrackets() returned an unexpected error: got: %v, want: %v", err, tt.err)
				}
			} else {
				if out != tt.expectString {
					t.Errorf("resolveAngleBrackets() returned an unexpected string: got: %s, want: %s", out, tt.expectString)
				}
			}
		})
	}
}

func TestAddQuotesForCurlyBraces(t *testing.T) {
	tests := []struct {
		originalString string
		expectString   string
		err            error
	}{
		{
			originalString: "test-ubuntu-20.04-linux",
			expectString:   "test-ubuntu-20.04-linux",
		},
		{
			originalString: "{{ .A }}-linux",
			expectString:   `"{{ .A }}-linux"`,
		},
		{
			originalString: `sn0rt-{{ .A }}-linux`,
			expectString:   `"sn0rt-{{ .A }}-linux"`,
		},
		{
			originalString: `sn0rt-{{ .A }}-{{ .B }}`,
			expectString:   `"sn0rt-{{ .A }}-{{ .B }}"`,
		},
		{
			originalString: `{{ .A }}-{{ .B }} {{ .C }}`,
			expectString:   `"{{ .A }}-{{ .B }}" "{{ .C }}"`,
		},
		{
			originalString: `sn0rt-{{ .A }}-{{ .B }} {{ .C }}`,
			expectString:   `"sn0rt-{{ .A }}-{{ .B }}" "{{ .C }}"`,
		},
		{
			originalString: `{{ .MYSQL_PASSWD }}`,
			expectString:   `"{{ .MYSQL_PASSWD }}"`,
		},
		{
			originalString: `password = {{ .MYSQL_PASSWD }}`,
			expectString:   `password = "{{ .MYSQL_PASSWD }}"`,
		},
		{
			originalString: `[client]
host = example.com
user = {{ .USER }}
password = {{ .MYSQL_PASSWD }}
port = 4000`,
			expectString: `[client]
host = example.com
user = "{{ .USER }}"
password = "{{ .MYSQL_PASSWD }}"
port = 4000`,
		},
		{
			originalString: `[client]
host = example.com
user = {{ .USER }}-password
password = linux-{{ .MYSQL_PASSWD }}
port = 4000`,
			expectString: `[client]
host = example.com
user = "{{ .USER }}-password"
password = "linux-{{ .MYSQL_PASSWD }}"
port = 4000`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.originalString, func(t *testing.T) {
			out := addQuotesCurlyBraces(tt.originalString)
			if out != tt.expectString {
				t.Errorf("resolveAngleBrackets() returned an unexpected string: got: %v, want: %s", out, tt.expectString)
			}
		})
	}
}
