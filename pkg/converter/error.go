package converter

// for secret common check
const (
	ErrCommonNotEmptyAnnotations               = "not empty annotations: %s"
	ErrCommonNotFoundAVPPath                   = "not found avp.kubernetes.io/path: %s"
	ErrCommonNotAcceptBothSecretDataAndData    = "not accept both Data and stringData Fields %s"
	ErrCommonNotAcceptNeitherSecretDataAndData = "not accept neither Data and stringData Fields %s"
	NotImplSecretType                          = "not impl %s secret type: secret: %s"
	illegalStoreType                           = "illegal store type: %s"
	illegalVaultPath                           = "illegal vault path: %s"
	NotSupportedSecretData                     = "secret support only Data or stringData: %s"
	FileContentAngleBracketsParseSyntaxError   = "template syntax error: %s"
)

const (
	ErrOpaqueNotAllowDataAndStringData      = "kubernetes.io/opaque type should not allow set Data and stringData Fields %s"
	ErrOpaqueNotAllowEmptyDataAndStringData = "kubernetes.io/opaque type should not allow empty Data and stringData Fields %s"
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
