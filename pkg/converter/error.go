package converter

// for secret common check
const (
	ErrCommonEmptyAnnotations                  = "not accept empty annotations of secret: %s"
	ErrCommonNotFoundAVPPath                   = "not found avp.kubernetes.io/path from secret: %s"
	ErrCommonNotAcceptBothSecretDataAndData    = "not accept both Data and stringData Fields of secret: %s"
	ErrCommonNotAcceptNeitherSecretDataAndData = "not accept neither Data and stringData Fields of secret %s"
	ErrCommonNotIncludeAngleBrackets           = "not include any angle brackets of secret: %s"
	ErrCommonNotNeedRefData                    = "not need ref data of secret: %s"
	ErrCommonNotSetEnv                         = "not set ENV: %s"
	NotImplSecretType                          = "not impl %s secret type of secret: %s"
	illegalStoreType                           = "illegal store type: %s"
	illegalVaultPath                           = "illegal vault path: %s"
	illegalCreatePolicy                        = "illegal create policy: %s, only support Owner, Orphan"
	FileContentAngleBracketsParseSyntaxError   = "template syntax error: %s"
)

const (
	ErrBasicAuthNotAllowDataField = "kubernetes.io/basic-auth type should not allow set Data Fields %s"
	ErrBasicAuthWithEmptyUsername = "basic auth secret with empty username: %s"
	ErrBasicAuthWithEmptyPassword = "basic auth secret with empty password: %s"
)

const (
	ErrDockerConfigJsonAcceptOnlyDataFields = "kubernetes.io/dockerconfigjson type should only accept set Data Fields %s"
	ErrDockerConfigJsonAcceptOnlyOneValue   = "kubernetes.io/dockerconfigjson type should only accept one value %s"
)

const (
	ErrTLSNotAllowDataField = "kubernetes.io/tls type should not allow set Data Fields %s"
)
