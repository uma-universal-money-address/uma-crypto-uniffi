package internal

/*


// This file was autogenerated by some hot garbage in the `uniffi` crate.
// Trust me, you don't want to mess with it!

#include <stdbool.h>
#include <stdint.h>

// The following structs are used to implement the lowest level
// of the FFI, and thus useful to multiple uniffied crates.
// We ensure they are declared exactly once, with a header guard, UNIFFI_SHARED_H.
#ifdef UNIFFI_SHARED_H
	// We also try to prevent mixing versions of shared uniffi header structs.
	// If you add anything to the #else block, you must increment the version suffix in UNIFFI_SHARED_HEADER_V4
	#ifndef UNIFFI_SHARED_HEADER_V4
		#error Combining helper code from multiple versions of uniffi is not supported
	#endif // ndef UNIFFI_SHARED_HEADER_V4
#else
#define UNIFFI_SHARED_H
#define UNIFFI_SHARED_HEADER_V4
// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V4 in this file.           ⚠️

typedef struct RustBuffer {
	int32_t capacity;
	int32_t len;
	uint8_t *data;
} RustBuffer;

typedef int32_t (*ForeignCallback)(uint64_t, int32_t, RustBuffer, RustBuffer *);

typedef struct ForeignBytes {
	int32_t len;
	const uint8_t *data;
} ForeignBytes;

// Error definitions
typedef struct RustCallStatus {
	int8_t code;
	RustBuffer errorBuf;
} RustCallStatus;

// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V4 in this file.           ⚠️
#endif // def UNIFFI_SHARED_H

void ffi_uma_crypto_b9a_KeyPair_object_free(
	void* ptr,
	RustCallStatus* out_status
);

RustBuffer uma_crypto_b9a_KeyPair_get_public_key(
	void* ptr,
	RustCallStatus* out_status
);

RustBuffer uma_crypto_b9a_KeyPair_get_private_key(
	void* ptr,
	RustCallStatus* out_status
);

RustBuffer uma_crypto_b9a_sign_ecdsa(
	RustBuffer msg,
	RustBuffer private_key_bytes,
	RustCallStatus* out_status
);

int8_t uma_crypto_b9a_verify_ecdsa(
	RustBuffer msg,
	RustBuffer signature_bytes,
	RustBuffer public_key_bytes,
	RustCallStatus* out_status
);

RustBuffer uma_crypto_b9a_encrypt_ecies(
	RustBuffer msg,
	RustBuffer public_key_bytes,
	RustCallStatus* out_status
);

RustBuffer uma_crypto_b9a_decrypt_ecies(
	RustBuffer cipher_text,
	RustBuffer private_key_bytes,
	RustCallStatus* out_status
);

void* uma_crypto_b9a_generate_keypair(
	RustCallStatus* out_status
);

RustBuffer ffi_uma_crypto_b9a_rustbuffer_alloc(
	int32_t size,
	RustCallStatus* out_status
);

RustBuffer ffi_uma_crypto_b9a_rustbuffer_from_bytes(
	ForeignBytes bytes,
	RustCallStatus* out_status
);

void ffi_uma_crypto_b9a_rustbuffer_free(
	RustBuffer buf,
	RustCallStatus* out_status
);

RustBuffer ffi_uma_crypto_b9a_rustbuffer_reserve(
	RustBuffer buf,
	int32_t additional,
	RustCallStatus* out_status
);


*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type rustBuffer struct {
	capacity int
	length   int
	data     unsafe.Pointer
	self     C.RustBuffer
}

func fromCRustBuffer(crbuf C.RustBuffer) rustBuffer {
	return rustBuffer{
		capacity: int(crbuf.capacity),
		length:   int(crbuf.len),
		data:     unsafe.Pointer(crbuf.data),
		self:     crbuf,
	}
}

// asByteBuffer reads the full rust buffer and then converts read bytes to a new reader which makes
// it quite inefficient
// TODO: Return an implementation which reads only when needed
func (rb rustBuffer) asReader() *bytes.Reader {
	b := C.GoBytes(rb.data, C.int(rb.length))
	return bytes.NewReader(b)
}

func (rb rustBuffer) asCRustBuffer() C.RustBuffer {
	return C.RustBuffer{
		capacity: C.int(rb.capacity),
		len:      C.int(rb.length),
		data:     (*C.uchar)(unsafe.Pointer(rb.data)),
	}
}

func stringToCRustBuffer(str string) C.RustBuffer {
	b := []byte(str)
	cs := C.CString(str)
	return C.RustBuffer{
		capacity: C.int(len(b)),
		len:      C.int(len(b)),
		data:     (*C.uchar)(unsafe.Pointer(cs)),
	}
}

func (rb rustBuffer) free() {
	rustCall(func(status *C.RustCallStatus) bool {
		C.ffi_uma_crypto_b9a_rustbuffer_free(rb.self, status)
		return false
	})
}

type bufLifter[GoType any] interface {
	lift(value C.RustBuffer) GoType
}

type bufLowerer[GoType any] interface {
	lower(value GoType) C.RustBuffer
}

type ffiConverter[GoType any, FfiType any] interface {
	lift(value FfiType) GoType
	lower(value GoType) FfiType
}

type bufReader[GoType any] interface {
	read(reader io.Reader) GoType
}

type bufWriter[GoType any] interface {
	write(writer io.Writer, value GoType)
}

type ffiRustBufConverter[GoType any, FfiType any] interface {
	ffiConverter[GoType, FfiType]
	bufReader[GoType]
}

func lowerIntoRustBuffer[GoType any](bufWriter bufWriter[GoType], value GoType) C.RustBuffer {
	// This might be not the most efficient way but it does not require knowing allocation size
	// beforehand
	var buffer bytes.Buffer
	bufWriter.write(&buffer, value)

	bytes, err := io.ReadAll(&buffer)
	if err != nil {
		panic(fmt.Errorf("reading written data: %w", err))
	}

	return stringToCRustBuffer(string(bytes))
}

func liftFromRustBuffer[GoType any](bufReader bufReader[GoType], rbuf rustBuffer) GoType {
	defer rbuf.free()
	reader := rbuf.asReader()
	item := bufReader.read(reader)
	if reader.Len() > 0 {
		// TODO: Remove this
		leftover, _ := io.ReadAll(reader)
		panic(fmt.Errorf("Junk remaining in buffer after lifting: %s", string(leftover)))
	}
	return item
}

func rustCallWithError[U any](converter bufLifter[error], callback func(*C.RustCallStatus) U) (U, error) {
	var status C.RustCallStatus
	returnValue := callback(&status)
	switch status.code {
	case 0:
		return returnValue, nil
	case 1:
		return returnValue, converter.lift(status.errorBuf)
	case 2:
		// when the rust code sees a panic, it tries to construct a rustbuffer
		// with the message.  but if that code panics, then it just sends back
		// an empty buffer.
		if status.errorBuf.len > 0 {
			panic(fmt.Errorf("%s", FfiConverterstringINSTANCE.lift(status.errorBuf)))
		} else {
			panic(fmt.Errorf("Rust panicked while handling Rust panic"))
		}
	default:
		return returnValue, fmt.Errorf("unknown status code: %d", status.code)
	}
}

func rustCall[U any](callback func(*C.RustCallStatus) U) U {
	returnValue, err := rustCallWithError(nil, callback)
	if err != nil {
		panic(err)
	}
	return returnValue
}

func writeInt8(writer io.Writer, value int8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint8(writer io.Writer, value uint8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt16(writer io.Writer, value int16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint16(writer io.Writer, value uint16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt32(writer io.Writer, value int32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint32(writer io.Writer, value uint32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt64(writer io.Writer, value int64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint64(writer io.Writer, value uint64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat32(writer io.Writer, value float32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat64(writer io.Writer, value float64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func readInt8(reader io.Reader) int8 {
	var result int8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint8(reader io.Reader) uint8 {
	var result uint8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt16(reader io.Reader) int16 {
	var result int16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint16(reader io.Reader) uint16 {
	var result uint16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt32(reader io.Reader) int32 {
	var result int32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint32(reader io.Reader) uint32 {
	var result uint32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt64(reader io.Reader) int64 {
	var result int64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint64(reader io.Reader) uint64 {
	var result uint64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat32(reader io.Reader) float32 {
	var result float32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat64(reader io.Reader) float64 {
	var result float64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func init() {

}

type FfiConverteruint8 struct{}

var FfiConverteruint8INSTANCE = FfiConverteruint8{}

func (FfiConverteruint8) lower(value uint8) C.uint8_t {
	return C.uint8_t(value)
}

func (FfiConverteruint8) write(writer io.Writer, value uint8) {
	writeUint8(writer, value)
}

func (FfiConverteruint8) lift(value C.uint8_t) uint8 {
	return uint8(value)
}

func (FfiConverteruint8) read(reader io.Reader) uint8 {
	return readUint8(reader)
}

type FfiDestroyeruint8 struct{}

func (FfiDestroyeruint8) destroy(_ uint8) {}

type FfiConverterbool struct{}

var FfiConverterboolINSTANCE = FfiConverterbool{}

func (FfiConverterbool) lower(value bool) C.int8_t {
	if value {
		return C.int8_t(1)
	}
	return C.int8_t(0)
}

func (FfiConverterbool) write(writer io.Writer, value bool) {
	if value {
		writeInt8(writer, 1)
	} else {
		writeInt8(writer, 0)
	}
}

func (FfiConverterbool) lift(value C.int8_t) bool {
	return value != 0
}

func (FfiConverterbool) read(reader io.Reader) bool {
	return readInt8(reader) != 0
}

type FfiDestroyerbool struct{}

func (FfiDestroyerbool) destroy(_ bool) {}

type FfiConverterstring struct{}

var FfiConverterstringINSTANCE = FfiConverterstring{}

func (FfiConverterstring) lift(cRustBuf C.RustBuffer) string {
	rustBuf := fromCRustBuffer(cRustBuf)
	defer rustBuf.free()

	reader := rustBuf.asReader()
	b, err := io.ReadAll(reader)
	if err != nil {
		panic(fmt.Errorf("reading reader: %w", err))
	}
	return string(b)
}

func (FfiConverterstring) read(reader io.Reader) string {
	length := readInt32(reader)
	buffer := make([]byte, length)
	read_length, err := reader.Read(buffer)
	if err != nil {
		panic(err)
	}
	if read_length != int(length) {
		panic(fmt.Errorf("bad read length when reading string, expected %d, read %d", length, read_length))
	}
	return string(buffer)
}

func (FfiConverterstring) lower(value string) C.RustBuffer {
	return stringToCRustBuffer(value)
}

func (FfiConverterstring) write(writer io.Writer, value string) {
	if len(value) > math.MaxInt32 {
		panic("String is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	write_length, err := io.WriteString(writer, value)
	if err != nil {
		panic(err)
	}
	if write_length != len(value) {
		panic(fmt.Errorf("bad write length when writing string, expected %d, written %d", len(value), write_length))
	}
}

type FfiDestroyerstring struct{}

func (FfiDestroyerstring) destroy(_ string) {}

// Below is an implementation of synchronization requirements outlined in the link.
// https://github.com/mozilla/uniffi-rs/blob/0dc031132d9493ca812c3af6e7dd60ad2ea95bf0/uniffi_bindgen/src/bindings/kotlin/templates/ObjectRuntime.kt#L31

type FfiObject struct {
	pointer      unsafe.Pointer
	callCounter  atomic.Int64
	freeFunction func(unsafe.Pointer, *C.RustCallStatus)
	destroyed    atomic.Bool
}

func newFfiObject(pointer unsafe.Pointer, freeFunction func(unsafe.Pointer, *C.RustCallStatus)) FfiObject {
	return FfiObject{
		pointer:      pointer,
		freeFunction: freeFunction,
	}
}

func (ffiObject *FfiObject) incrementPointer(debugName string) unsafe.Pointer {
	for {
		counter := ffiObject.callCounter.Load()
		if counter <= -1 {
			panic(fmt.Errorf("%v object has already been destroyed", debugName))
		}
		if counter == math.MaxInt64 {
			panic(fmt.Errorf("%v object call counter would overflow", debugName))
		}
		if ffiObject.callCounter.CompareAndSwap(counter, counter+1) {
			break
		}
	}

	return ffiObject.pointer
}

func (ffiObject *FfiObject) decrementPointer() {
	if ffiObject.callCounter.Add(-1) == -1 {
		ffiObject.freeRustArcPtr()
	}
}

func (ffiObject *FfiObject) destroy() {
	if ffiObject.destroyed.CompareAndSwap(false, true) {
		if ffiObject.callCounter.Add(-1) == -1 {
			ffiObject.freeRustArcPtr()
		}
	}
}

func (ffiObject *FfiObject) freeRustArcPtr() {
	rustCall(func(status *C.RustCallStatus) int32 {
		ffiObject.freeFunction(ffiObject.pointer, status)
		return 0
	})
}

type KeyPair struct {
	ffiObject FfiObject
}

func (_self *KeyPair) GetPublicKey() []uint8 {
	_pointer := _self.ffiObject.incrementPointer("*KeyPair")
	defer _self.ffiObject.decrementPointer()

	return FfiConverterSequenceuint8INSTANCE.lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.uma_crypto_b9a_KeyPair_get_public_key(
			_pointer, _uniffiStatus)
	}))

}
func (_self *KeyPair) GetPrivateKey() []uint8 {
	_pointer := _self.ffiObject.incrementPointer("*KeyPair")
	defer _self.ffiObject.decrementPointer()

	return FfiConverterSequenceuint8INSTANCE.lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.uma_crypto_b9a_KeyPair_get_private_key(
			_pointer, _uniffiStatus)
	}))

}

func (object *KeyPair) Destroy() {
	runtime.SetFinalizer(object, nil)
	object.ffiObject.destroy()
}

type FfiConverterKeyPair struct{}

var FfiConverterKeyPairINSTANCE = FfiConverterKeyPair{}

func (c FfiConverterKeyPair) lift(pointer unsafe.Pointer) *KeyPair {
	result := &KeyPair{
		newFfiObject(
			pointer,
			func(pointer unsafe.Pointer, status *C.RustCallStatus) {
				C.ffi_uma_crypto_b9a_KeyPair_object_free(pointer, status)
			}),
	}
	runtime.SetFinalizer(result, (*KeyPair).Destroy)
	return result
}

func (c FfiConverterKeyPair) read(reader io.Reader) *KeyPair {
	return c.lift(unsafe.Pointer(uintptr(readUint64(reader))))
}

func (c FfiConverterKeyPair) lower(value *KeyPair) unsafe.Pointer {
	// TODO: this is bad - all synchronization from ObjectRuntime.go is discarded here,
	// because the pointer will be decremented immediately after this function returns,
	// and someone will be left holding onto a non-locked pointer.
	pointer := value.ffiObject.incrementPointer("*KeyPair")
	defer value.ffiObject.decrementPointer()
	return pointer
}

func (c FfiConverterKeyPair) write(writer io.Writer, value *KeyPair) {
	writeUint64(writer, uint64(uintptr(c.lower(value))))
}

type FfiDestroyerKeyPair struct{}

func (_ FfiDestroyerKeyPair) destroy(value *KeyPair) {
	value.Destroy()
}

type CryptoError struct {
	err error
}

func (err CryptoError) Error() string {
	return fmt.Sprintf("CryptoError: %s", err.err.Error())
}

func (err CryptoError) Unwrap() error {
	return err.err
}

// Err* are used for checking error type with `errors.Is`
var ErrCryptoErrorSecp256k1Error = fmt.Errorf("CryptoErrorSecp256k1Error")

// Variant structs
type CryptoErrorSecp256k1Error struct {
	message string
}

func NewCryptoErrorSecp256k1Error() *CryptoError {
	return &CryptoError{
		err: &CryptoErrorSecp256k1Error{},
	}
}

func (err CryptoErrorSecp256k1Error) Error() string {
	return fmt.Sprintf("Secp256k1Error: %s", err.message)
}

func (self CryptoErrorSecp256k1Error) Is(target error) bool {
	return target == ErrCryptoErrorSecp256k1Error
}

type FfiConverterTypeCryptoError struct{}

var FfiConverterTypeCryptoErrorINSTANCE = FfiConverterTypeCryptoError{}

func (c FfiConverterTypeCryptoError) lift(cErrBuf C.RustBuffer) error {
	errBuf := fromCRustBuffer(cErrBuf)
	return liftFromRustBuffer[error](c, errBuf)
}

func (c FfiConverterTypeCryptoError) lower(value *CryptoError) C.RustBuffer {
	return lowerIntoRustBuffer[*CryptoError](c, value)
}

func (c FfiConverterTypeCryptoError) read(reader io.Reader) error {
	errorID := readUint32(reader)

	message := FfiConverterstringINSTANCE.read(reader)
	switch errorID {
	case 1:
		return &CryptoError{&CryptoErrorSecp256k1Error{message}}
	default:
		panic(fmt.Sprintf("Unknown error code %d in FfiConverterTypeCryptoError.read()", errorID))
	}

}

func (c FfiConverterTypeCryptoError) write(writer io.Writer, value *CryptoError) {
	switch variantValue := value.err.(type) {
	case *CryptoErrorSecp256k1Error:
		writeInt32(writer, 1)
	default:
		_ = variantValue
		panic(fmt.Sprintf("invalid error value `%v` in FfiConverterTypeCryptoError.write", value))
	}
}

type FfiConverterSequenceuint8 struct{}

var FfiConverterSequenceuint8INSTANCE = FfiConverterSequenceuint8{}

func (c FfiConverterSequenceuint8) lift(cRustBuf C.RustBuffer) []uint8 {
	return liftFromRustBuffer[[]uint8](c, fromCRustBuffer(cRustBuf))
}

func (c FfiConverterSequenceuint8) read(reader io.Reader) []uint8 {
	length := readInt32(reader)
	if length == 0 {
		return nil
	}
	result := make([]uint8, 0, length)
	for i := int32(0); i < length; i++ {
		result = append(result, FfiConverteruint8INSTANCE.read(reader))
	}
	return result
}

func (c FfiConverterSequenceuint8) lower(value []uint8) C.RustBuffer {
	return lowerIntoRustBuffer[[]uint8](c, value)
}

func (c FfiConverterSequenceuint8) write(writer io.Writer, value []uint8) {
	if len(value) > math.MaxInt32 {
		panic("[]uint8 is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	for _, item := range value {
		FfiConverteruint8INSTANCE.write(writer, item)
	}
}

type FfiDestroyerSequenceuint8 struct{}

func (FfiDestroyerSequenceuint8) destroy(sequence []uint8) {
	for _, value := range sequence {
		FfiDestroyeruint8{}.destroy(value)
	}
}

func SignEcdsa(msg []uint8, privateKeyBytes []uint8) ([]uint8, error) {

	_uniffiRV, _uniffiErr := rustCallWithError(FfiConverterTypeCryptoError{}, func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.uma_crypto_b9a_sign_ecdsa(FfiConverterSequenceuint8INSTANCE.lower(msg), FfiConverterSequenceuint8INSTANCE.lower(privateKeyBytes), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue []uint8
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterSequenceuint8INSTANCE.lift(_uniffiRV), _uniffiErr
	}

}

func VerifyEcdsa(msg []uint8, signatureBytes []uint8, publicKeyBytes []uint8) (bool, error) {

	_uniffiRV, _uniffiErr := rustCallWithError(FfiConverterTypeCryptoError{}, func(_uniffiStatus *C.RustCallStatus) C.int8_t {
		return C.uma_crypto_b9a_verify_ecdsa(FfiConverterSequenceuint8INSTANCE.lower(msg), FfiConverterSequenceuint8INSTANCE.lower(signatureBytes), FfiConverterSequenceuint8INSTANCE.lower(publicKeyBytes), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue bool
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterboolINSTANCE.lift(_uniffiRV), _uniffiErr
	}

}

func EncryptEcies(msg []uint8, publicKeyBytes []uint8) ([]uint8, error) {

	_uniffiRV, _uniffiErr := rustCallWithError(FfiConverterTypeCryptoError{}, func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.uma_crypto_b9a_encrypt_ecies(FfiConverterSequenceuint8INSTANCE.lower(msg), FfiConverterSequenceuint8INSTANCE.lower(publicKeyBytes), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue []uint8
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterSequenceuint8INSTANCE.lift(_uniffiRV), _uniffiErr
	}

}

func DecryptEcies(cipherText []uint8, privateKeyBytes []uint8) ([]uint8, error) {

	_uniffiRV, _uniffiErr := rustCallWithError(FfiConverterTypeCryptoError{}, func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.uma_crypto_b9a_decrypt_ecies(FfiConverterSequenceuint8INSTANCE.lower(cipherText), FfiConverterSequenceuint8INSTANCE.lower(privateKeyBytes), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue []uint8
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterSequenceuint8INSTANCE.lift(_uniffiRV), _uniffiErr
	}

}

func GenerateKeypair() (*KeyPair, error) {

	_uniffiRV, _uniffiErr := rustCallWithError(FfiConverterTypeCryptoError{}, func(_uniffiStatus *C.RustCallStatus) unsafe.Pointer {
		return C.uma_crypto_b9a_generate_keypair(_uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue *KeyPair
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterKeyPairINSTANCE.lift(_uniffiRV), _uniffiErr
	}

}
