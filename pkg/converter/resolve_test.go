package converter

import (
	"fmt"
	"os"
	"testing"
)

func TestResolved(t *testing.T) {
	tests := []struct {
		name           string
		originalString string
		expectString   string
		envs           map[string]string // for render <% ENV %>
		err            error
	}{
		{
			name:           "simple case",
			originalString: "test-ubuntu-20.04-linux",
			expectString:   "test-ubuntu-20.04-linux",
		},
		{
			name:           "include_space_and_env",
			originalString: "<%      ENV    %>-linux",
			expectString:   "test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			name:           "include_many_space_and_env",
			originalString: "secret/foo/<%      ENV    %>-linux",
			expectString:   "secret/foo/test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			name:           "include_many_space_and_env_2",
			originalString: "<% ENV    %>-linux",
			expectString:   "test-linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			name:           "include_many_env",
			originalString: "<% ENV %>-<% DIST %>-<% VER %>-linux",
			expectString:   "test-ubuntu-20.04-linux",
			envs: map[string]string{
				"ENV":  "test",
				"DIST": "ubuntu",
				"VER":  "20.04",
			},
		},
		{
			name:           "illegal_env",
			originalString: "<%      ENV    -linux",
			expectString:   "<%      ENV    -linux",
			envs: map[string]string{
				"ENV": "test",
			},
		},
		{
			name:           "can_not_been_resolve",
			originalString: "<%      NOTSETVAR  %>-linux",
			err:            fmt.Errorf(ErrCommonNotSetEnv, "NOTSETVAR"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				_ = os.Setenv(k, v)
			}
			out, err := resolved(tt.originalString)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("resolved() returned an unexpected error: got: %v, want: %v", err, tt.err)
				}
			} else {
				if out != tt.expectString {
					t.Errorf("resolved() returned an unexpected string: got: %s, want: %s", out, tt.expectString)
				}
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
		name           string
		originalString string
		expectString   string
		err            error
	}{
		{
			name:           "simple",
			originalString: "test-ubuntu-20.04-linux",
			expectString:   "test-ubuntu-20.04-linux",
		},
		{
			name:           "include_angle_brackets",
			originalString: "<A>-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			name:           "include_many_space_and_angle_brackets",
			originalString: "<   A   >-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			name:           "include_many_space_and_angle_brackets_2",
			originalString: "<A    >-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			name:           "include_many_space_and_angle_brackets_3",
			originalString: "<   A>-linux",
			expectString:   "{{ .A }}-linux",
		},
		{
			name:           "angle_brackets_connected",
			originalString: "sn0rt-<A>-linux",
			expectString:   "sn0rt-{{ .A }}-linux",
		},
		{
			name:           "angle_brackets_connected_2",
			originalString: "sn0rt-<A>-<B>",
			expectString:   "sn0rt-{{ .A }}-{{ .B }}",
		},
		{
			name:           "illegal_angle_brackets",
			originalString: "sn0rt-<A",
			expectString:   "sn0rt-<A",
			err:            fmt.Errorf(FileContentAngleBracketsParseSyntaxError, `unclosed '<'`),
		},
		{
			name:           "multi_angle_brackets",
			originalString: "<A>-<B> <C>",
			expectString:   "{{ .A }}-{{ .B }} {{ .C }}",
		},
		{
			name:           "multi_angle_brackets_2",
			originalString: "sn0rt-<A>-<B> <C>",
			expectString:   "sn0rt-{{ .A }}-{{ .B }} {{ .C }}",
		},
		{
			name:           "simple_case_2",
			originalString: "<MYSQL_PASSWD>",
			expectString:   "{{ .MYSQL_PASSWD }}",
		},
		{
			name:           "include_equal_sign",
			originalString: "password = <MYSQL_PASSWD>",
			expectString:   "password = {{ .MYSQL_PASSWD }}",
		},
		{
			name: "simple_file",
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
		{
			name: "yaml_file",
			originalString: `type: S3
prefix: "test/ubuntu"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: <S3_ACCESS_KEY>
  secret_key: <S3_SECRET_KEY>`,
			expectString: `type: S3
prefix: "test/ubuntu"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: {{ .S3_ACCESS_KEY }}
  secret_key: {{ .S3_SECRET_KEY }}`,
		},
		{
			name:           "include <% ENV %>",
			originalString: "password = <% MYSQL_PASSWD %>",
			expectString:   "password = <% MYSQL_PASSWD %>",
		},
		{
			name: "yaml_file with env and tab",
			originalString: `type: S3
prefix: "test/ubuntu"
config:
	endpoint: "https://s3.amazonaws.com"
	access_key: <% S3_ACCESS_KEY %>
	secret_key: <S3_SECRET_KEY>`,
			expectString: `type: S3
prefix: "test/ubuntu"
config:
	endpoint: "https://s3.amazonaws.com"
	access_key: <% S3_ACCESS_KEY %>
	secret_key: {{ .S3_SECRET_KEY }}`,
		},
		{
			name: "yaml_file with env and space",
			originalString: `type: S3
prefix: "test/ubuntu"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: <% S3_ACCESS_KEY %>
  secret_key: <S3_SECRET_KEY>`,
			expectString: `type: S3
prefix: "test/ubuntu"
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: <% S3_ACCESS_KEY %>
  secret_key: {{ .S3_SECRET_KEY }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := resolveAngleBrackets(tt.originalString)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("resolveAngleBrackets() returned an unexpected error: goot: %v", err)
					t.Errorf("resolveAngleBrackets() returned an unexpected error: want: %v", tt.err)

				}
			} else {
				if out != tt.expectString {
					t.Errorf("resolveAngleBrackets() returned an unexpected string: got: %v, want: %s", out, tt.expectString)
					fmt.Printf("Got length: %d, Want length: %d\n", len(out), len(tt.expectString))
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
			originalString: "<% linux-test %>",
			expectString:   `<% linux-test %>`,
		},
		{
			originalString: "{{ .A }}-linux",
			expectString:   `"{{ .A }}-linux"`,
		},
		{
			originalString: "{{ .A }}-linux-<% ENV %>",
			expectString:   `"{{ .A }}-linux-<% ENV %>"`,
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
		{
			originalString: `type: S3
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: {{ .S3_ACCESS_KEY }}
  secret_key: {{ .S3_SECRET_KEY }}`,
			expectString: `type: S3
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: "{{ .S3_ACCESS_KEY }}"
  secret_key: "{{ .S3_SECRET_KEY }}"`,
		},
		{
			originalString: `type: S3
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: {{ .S3_ACCESS_KEY }}-<% ENV %>
  secret_key: {{ .S3_SECRET_KEY }}`,
			expectString: `type: S3
config:
  endpoint: "https://s3.amazonaws.com"
  access_key: "{{ .S3_ACCESS_KEY }}-<% ENV %>"
  secret_key: "{{ .S3_SECRET_KEY }}"`,
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			out := addQuotesCurlyBraces(tt.originalString)
			if out != tt.expectString {
				t.Errorf("addQuotesCurlyBraces() returned an unexpected string: got: %v, want: %s", out, tt.expectString)
			}
		})
	}
}

func Test_processCommented(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "should_get_empty",
			input: []byte(`# +optional`),
			want:  []byte(``),
		},
		{
			name: "should_get_1_line",
			input: []byte(`# +optional
metav1.ObjectMeta`),
			want: []byte(`metav1.ObjectMeta`),
		},
		{
			name: "should_get_2_yaml",
			input: []byte(`apiVersion: v1
kind: Secret
---
#apiVersion: v1
#kind: Secret
#metadata:
#  name: input7
#  annotations:
#    avp.kubernetes.io/path: "secret/data/test-foo"
#type: kubernetes.io/tls
#data:
#  tls.crt: <TLS_CRT>
#  tls.key: <TLS_KEY>
---
apiVersion: v1
kind: Secret`),
			want: []byte(`apiVersion: v1
kind: Secret
---
---
apiVersion: v1
kind: Secret`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processCommented(tt.input)
			if string(got) != string(tt.want) {
				t.Errorf("processCommented() got:%v, want:%v", string(got), string(tt.want))
			}
		})
	}
}
