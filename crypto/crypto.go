package crypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"
	"io"
)

// AESEncrypt 进行AES-GCM加密，nonce放置于密文之前
func AESEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, 12)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESDecrypt 进行AES-GCM解密，nonce放置于密文之前
func AESDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, 12)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bits < 2048 {
		return nil, nil, errors.New("ras key size too small, must be at least 2048 bits")
	}
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	pubKey := &privKey.PublicKey
	return privKey, pubKey, nil
}

// PublicKeyToPEM 公钥进行pem编码
func PublicKeyToPEM(pubKey *rsa.PublicKey) ([]byte, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pubPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	return pem.EncodeToMemory(pubPEM), nil
}

// PrivateKeyToPEM 私钥进行pem编码
func PrivateKeyToPEM(privKey *rsa.PrivateKey) ([]byte, error) {
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	privPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	return pem.EncodeToMemory(privPEM), nil
}

// PEMToPublicKey 从pem中解码出公钥
func PEMToPublicKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key pem")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an rsa public key")
	}
	return rsaPub, nil
}

// PEMToPrivateKey 从pem中解码出私钥
func PEMToPrivateKey(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key pem")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an rsa private key")
	}
	return rsaPriv, nil
}

// RSAEncrypt 使用公钥进行RSA-OAEP加密
func RSAEncrypt(plaintext []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// RSADecrypt 使用私钥进行RSA-OAEP解密
func RSADecrypt(ciphertext []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// RSASign 使用私钥进行RSA-PSS签名
func RSASign(data []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// RSAVerify 使用公钥进行RSA-PSS验签
func RSAVerify(data, signature []byte, pubKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
}

// HMACSign 使用hmac-sha256进行签名
func HMACSign(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HMACVerify 使用hmac-sha256进行验签
func HMACVerify(data, signature, key []byte) bool {
	expectedMAC := HMACSign(data, key)
	return hmac.Equal(expectedMAC, signature)
}

// HashPassword 使用bcrypt对密码进行哈希处理
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// VerifyPassword 使用bcrypt验证密码是否与哈希匹配
func VerifyPassword(password string, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}

// SecureRandom 生成一个加密安全伪随机数
func SecureRandom(length int) ([]byte, error) {
	random := make([]byte, length)
	_, err := rand.Read(random)
	if err != nil {
		return nil, err
	}
	return random, nil
}

// DeriveKey 使用scrypt进行秘钥派生
func DeriveKey(password, salt []byte, keyLen int) ([]byte, error) {
	if len(salt) == 0 {
		return nil, errors.New("salt must not be empty")
	}
	return scrypt.Key(password, salt, 32768, 8, 1, keyLen)
}

// StretchKey 使hkdf进行密钥拉伸，将一个低熵的秘钥拉伸为一个高质量密钥
func StretchKey(secret, salt []byte, keyLen int) ([]byte, error) {
	if len(salt) == 0 {
		return nil, errors.New("salt must not be empty")
	}
	reader := hkdf.New(sha256.New, secret, salt, nil)
	key := make([]byte, keyLen)
	if _, err := io.ReadFull(reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// SHA256Hash 计算数据的SHA-256哈希值
func SHA256Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// SHA1Hash 计算数据的SHA-1哈希值
func SHA1Hash(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// MD5Hash 计算数据的MD5哈希值
func MD5Hash(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}
