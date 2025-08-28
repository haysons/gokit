package crypto

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESEncryptAndDecrypt(t *testing.T) {
	plaintext := []byte("This is a test message for AES-GCM encryption.")
	key := []byte("12345678901234567890123456789012") // 32 bytes key for AES-256

	ciphertext, err := AESEncrypt(plaintext, key)
	assert.NoError(t, err, "AESEncrypt should not return an error")

	decryptedText, err := AESDecrypt(ciphertext, key)
	assert.NoError(t, err, "AESDecrypt should not return an error")

	assert.True(t, bytes.Equal(plaintext, decryptedText), "Decrypted text should be the same as the original plaintext")
}

func TestAESDecryptInvalidCiphertext(t *testing.T) {
	invalidCiphertext := []byte("invalidciphertext")

	_, err := AESDecrypt(invalidCiphertext, []byte("12345678901234567890123456789012"))
	assert.Error(t, err, "AESDecrypt should return an error for invalid ciphertext")
}

func TestAESDecryptWithWrongKey(t *testing.T) {
	plaintext := []byte("This is a test message for AES-GCM encryption.")
	key := []byte("12345678901234567890123456789012")
	ciphertext, err := AESEncrypt(plaintext, key)
	assert.NoError(t, err, "AESEncrypt should not return an error")

	wrongKey := []byte("wrongkey123456789012345678901234")
	_, err = AESDecrypt(ciphertext, wrongKey)
	assert.Error(t, err, "AESDecrypt should return an error for incorrect key")
}

func TestAESDecryptWithShortCiphertext(t *testing.T) {
	invalidCiphertext := []byte("short")

	_, err := AESDecrypt(invalidCiphertext, []byte("12345678901234567890123456789012"))
	assert.Error(t, err, "AESDecrypt should return an error for ciphertext that is too short")
}

func TestGenerateRSAKeyPair(t *testing.T) {
	_, _, err := GenerateRSAKeyPair(1024)
	assert.Error(t, err, "Should fail for key size less than 2048 bits")

	priv, pub, err := GenerateRSAKeyPair(2048)
	assert.NoError(t, err, "Key generation should succeed for 2048 bits")
	assert.NotNil(t, priv, "Private key should not be nil")
	assert.NotNil(t, pub, "Public key should not be nil")
}

func TestPEMEncodeDecode(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048)
	assert.NoError(t, err)

	pubPEM, err := PublicKeyToPEM(pub)
	assert.NoError(t, err)
	assert.NotEmpty(t, pubPEM, "PublicKey PEM should not be empty")

	decodedPub, err := PEMToPublicKey(pubPEM)
	assert.NoError(t, err)
	assert.Equal(t, pub.E, decodedPub.E, "Decoded public key exponent must match")
	assert.Equal(t, pub.N, decodedPub.N, "Decoded public key modulus must match")

	privPEM, err := PrivateKeyToPEM(priv)
	assert.NoError(t, err)
	assert.NotEmpty(t, privPEM, "PrivateKey PEM should not be empty")

	decodedPriv, err := PEMToPrivateKey(privPEM)
	assert.NoError(t, err)
	assert.Equal(t, priv.PublicKey.E, decodedPriv.PublicKey.E, "Decoded private key public exponent must match")
	assert.Equal(t, priv.PublicKey.N, decodedPriv.PublicKey.N, "Decoded private key modulus must match")
}

func TestRSAEncryptDecrypt(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048)
	assert.NoError(t, err)

	plaintext := []byte("This is a test message for RSA-OAEP encryption.")

	ciphertext, err := RSAEncrypt(plaintext, pub)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext, "Ciphertext should not be empty")

	decryptedText, err := RSADecrypt(ciphertext, priv)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(plaintext, decryptedText), "Decrypted text should match the original plaintext")
}

func TestRSASignVerify(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048)
	assert.NoError(t, err)

	data := []byte("This is a test message for RSA-PSS signing.")

	signature, err := RSASign(data, priv)
	assert.NoError(t, err)
	assert.NotEmpty(t, signature, "Signature should not be empty")

	err = RSAVerify(data, signature, pub)
	assert.NoError(t, err, "Signature verification should succeed")

	modifiedData := []byte("This is a modified message for RSA-PSS signing.")
	err = RSAVerify(modifiedData, signature, pub)
	assert.Error(t, err, "Signature verification should fail for modified data")
}

func TestPEMDecodeFailures(t *testing.T) {
	invalidPubPEM := []byte("-----BEGIN PUBLIC KEY-----\ninvalid\n-----END PUBLIC KEY-----")
	_, err := PEMToPublicKey(invalidPubPEM)
	assert.Error(t, err, "PEMToPublicKey should error on invalid PEM data")

	invalidPrivPEM := []byte("-----BEGIN RSA PRIVATE KEY-----\ninvalid\n-----END RSA PRIVATE KEY-----")
	_, err = PEMToPrivateKey(invalidPrivPEM)
	assert.Error(t, err, "PEMToPrivateKey should error on invalid PEM data")

	_, err = PEMToPublicKey([]byte{})
	assert.Error(t, err, "Empty PEM should return error for public key")
	_, err = PEMToPrivateKey([]byte{})
	assert.Error(t, err, "Empty PEM should return error for private key")
}

func TestHMACSignAndVerify(t *testing.T) {
	key := []byte("my-secret-key")
	message := []byte("important-message")

	signature := HMACSign(message, key)
	assert.NotEmpty(t, signature, "HMAC signature should not be empty")

	valid := HMACVerify(message, signature, key)
	assert.True(t, valid, "HMAC verification should succeed with correct key and message")

	wrongKey := []byte("wrong-key")
	valid = HMACVerify(message, signature, wrongKey)
	assert.False(t, valid, "HMAC verification should fail with incorrect key")

	wrongMessage := []byte("tampered-message")
	valid = HMACVerify(wrongMessage, signature, key)
	assert.False(t, valid, "HMAC verification should fail with tampered message")

	// Empty key and message
	var emptyKey []byte
	var emptyMessage []byte
	emptySig := HMACSign(emptyMessage, emptyKey)
	assert.True(t, HMACVerify(emptyMessage, emptySig, emptyKey), "HMAC should work with empty message and key")
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "securePassword123"

	// Test hashing
	hash, err := HashPassword(password)
	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hash, "Hashed password should not be empty")

	// Test verification with correct password
	match, err := VerifyPassword(password, hash)
	assert.NoError(t, err, "Verification with correct password should not return an error")
	assert.True(t, match, "Password should match the hash")

	// Test verification with incorrect password
	wrongPassword := "wrongPassword123"
	match, err = VerifyPassword(wrongPassword, hash)
	assert.NoError(t, err, "Verification with incorrect password should not return an error")
	assert.False(t, match, "Wrong password should not match the hash")

	// Test with empty password
	emptyHash, err := HashPassword("")
	assert.NoError(t, err, "Hashing empty password should not return an error")
	assert.NotEmpty(t, emptyHash, "Hashed empty password should not be empty")

	match, err = VerifyPassword("", emptyHash)
	assert.NoError(t, err, "Verifying empty password should not return an error")
	assert.True(t, match, "Empty password should match its hash")
}

func TestGenerateSalt(t *testing.T) {
	length := 16
	salt, err := SecureRandom(length)
	assert.NoError(t, err, "SecureRandom should not return an error")
	assert.Equal(t, length, len(salt), "Generated salt length should be %d", length)

	salt2, err := SecureRandom(length)
	assert.NoError(t, err, "SecureRandom (second call) should not return an error")
	assert.NotEqual(t, salt, salt2, "Subsequent generated salts should be different")
}

func TestDeriveKey(t *testing.T) {
	password := []byte("password123")
	keyLen := 32

	salt, err := SecureRandom(16)
	assert.NoError(t, err, "SecureRandom should not return error for valid input")

	derivedKey1, err := DeriveKey(password, salt, keyLen)
	assert.NoError(t, err, "DeriveKey should not return error for valid parameters")
	assert.Equal(t, keyLen, len(derivedKey1), "Derived key length should be equal to keyLen")

	derivedKey2, err := DeriveKey(password, salt, keyLen)
	assert.NoError(t, err, "DeriveKey (second call) should not return error")
	assert.Equal(t, derivedKey1, derivedKey2, "Derived keys should be identical for same parameters")

	_, err = DeriveKey(password, nil, keyLen)
	assert.Error(t, err, "DeriveKey should return error when salt is empty")
}

func TestDifferentSaltsProduceDifferentKeys(t *testing.T) {
	password := []byte("password123")
	keyLen := 32

	salt1, err := SecureRandom(16)
	assert.NoError(t, err, "SecureRandom for salt1 should not return error")
	salt2, err := SecureRandom(16)
	assert.NoError(t, err, "SecureRandom for salt2 should not return error")

	key1, err := DeriveKey(password, salt1, keyLen)
	assert.NoError(t, err, "DeriveKey with salt1 should not return error")
	key2, err := DeriveKey(password, salt2, keyLen)
	assert.NoError(t, err, "DeriveKey with salt2 should not return error")
	assert.NotEqual(t, key1, key2, "Derived keys should be different when using different salts")
}

func TestStretchKey(t *testing.T) {
	secret := []byte("low-entropy-secret")
	keyLen := 32

	salt := []byte("0123456789abcdef")

	key1, err := StretchKey(secret, salt, keyLen)
	assert.NoError(t, err, "StretchKey should not return an error with valid salt")
	assert.Equal(t, keyLen, len(key1), "Derived key should have length %d", keyLen)

	key2, err := StretchKey(secret, salt, keyLen)
	assert.NoError(t, err, "StretchKey (second call) should not return an error")
	assert.True(t, bytes.Equal(key1, key2), "Keys derived with the same parameters must match")

	anotherSalt := []byte("fedcba9876543210")
	key3, err := StretchKey(secret, anotherSalt, keyLen)
	assert.NoError(t, err, "StretchKey with a different salt should not return an error")
	assert.False(t, bytes.Equal(key1, key3), "Derived keys with different salts should be different")

	emptySalt := []byte("")
	_, err = StretchKey(secret, emptySalt, keyLen)
	assert.Error(t, err, "StretchKey should return an error when salt is empty")
}

func TestSHA256Hash(t *testing.T) {
	got := SHA256Hash([]byte("hello"))
	expected := sha256.Sum256([]byte("hello"))
	if !bytes.Equal(got, expected[:]) {
		t.Errorf("SHA256Hash() = %x, want %x", got, expected)
	}

	gotNil := SHA256Hash(nil)
	expectedNil := sha256.Sum256([]byte(""))
	if !bytes.Equal(gotNil, expectedNil[:]) {
		t.Errorf("SHA256Hash(nil) = %x, want %x", gotNil, expectedNil)
	}
}

func TestSHA1Hash(t *testing.T) {
	got := SHA1Hash([]byte("hello"))
	expected := sha1.Sum([]byte("hello"))
	if !bytes.Equal(got, expected[:]) {
		t.Errorf("SHA1Hash() = %x, want %x", got, expected)
	}

	gotNil := SHA1Hash(nil)
	expectedNil := sha1.Sum([]byte(""))
	if !bytes.Equal(gotNil, expectedNil[:]) {
		t.Errorf("SHA1Hash(nil) = %x, want %x", gotNil, expectedNil)
	}
}

func TestMD5Hash(t *testing.T) {
	got := MD5Hash([]byte("hello"))
	expected := md5.Sum([]byte("hello"))
	if !bytes.Equal(got, expected[:]) {
		t.Errorf("MD5Hash() = %x, want %x", got, expected)
	}

	gotNil := MD5Hash(nil)
	expectedNil := md5.Sum([]byte(""))
	if !bytes.Equal(gotNil, expectedNil[:]) {
		t.Errorf("MD5Hash(nil) = %x, want %x", gotNil, expectedNil)
	}
}
