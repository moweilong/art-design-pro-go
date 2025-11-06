package conversion

import (
	"github.com/moweilong/milady/pkg/core"

	"github.com/moweilong/art-design-pro-go/internal/apiserver/model"
	v1 "github.com/moweilong/art-design-pro-go/pkg/api/apiserver/v1"
)

// SecretMToSecretV1 converts a SecretM object from the internal model
// to a Secret object in the v1 API format.
func SecretMToSecretV1(secretModel *model.SecretM) *v1.Secret {
	var secret v1.Secret
	_ = core.CopyWithConverters(&secret, secretModel)
	return &secret
}

// SecretV1ToSecretM converts a Secret object from the v1 API format
// to a SecretM object in the internal model.
func SecretV1ToSecretM(secret *v1.Secret) *model.SecretM {
	var secretModel model.SecretM
	_ = core.CopyWithConverters(&secretModel, secret)
	return &secretModel
}
