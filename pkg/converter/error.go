package converter

const (
	NotSetAVPAnnotations                 = "not set AVP annotations, secret: %s"
	NotEmptyAnnotations                  = "not set AVP annotations, secret: %s"
	NotImplSecretType                    = "not impl %s secret type: secret: %s"
	NotSupportedStoreType                = "illegal store type: %s"
	NotSupportedSecretDataBothStringData = "secret support both Data and stringData: %s"
	NotSupportedSecretData               = "secret support only Data or stringData: %s"
	NotSupportedSecretDataEmpty          = "secret data is empty: %s"
)

const (
	illegalVaultPath = "illegal vault path: %s"
)
