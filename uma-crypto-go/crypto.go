package uma_crypto

// TODO(mhr): Dynamic linking?

// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/libs/darwin/amd64 -luma_crypto
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/libs/darwin/arm64 -luma_crypto
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/libs/linux/amd64 -Wl,-Bstatic -luma_crypto -Wl,-Bdynamic
// #cgo linux,arm64 LDFLAGS: -L${SRCDIR}/libs/linux/arm64 -Wl,-Bstatic -luma_crypto -Wl,-Bdynamic
import "C"
import (
	"github.com/uma-universal-money-address/uma-crypto-uniffi/uma-crypto-go/internal"
)

func SignEcdsa(message []byte, privateKey []byte) ([]byte, error) {
	return internal.SignEcdsa(message, privateKey)
}

func VerifyEcdsa(message []byte, signature []byte, publicKey []byte) (bool, error) {
	return internal.VerifyEcdsa(message, signature, publicKey)
}

func EncryptEcies(message []byte, publicKey []byte) ([]byte, error) {
	return internal.EncryptEcies(message, publicKey)
}

func DecryptEcies(message []byte, privateKey []byte) ([]byte, error) {
	return internal.DecryptEcies(message, privateKey)
}
