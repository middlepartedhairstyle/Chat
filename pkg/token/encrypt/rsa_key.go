package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/middlepartedhairstyle/HiWe/pkg/file"
	"os"
)

type RSAKey struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewRSAKey(bits int) *RSAKey {
	rsaKey := RSAKey{}
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil
	}
	rsaKey.PrivateKey = key
	rsaKey.PublicKey = key.Public().(*rsa.PublicKey)
	return &rsaKey
}

func ReadRSAKeyFile(path string, pr string, pu string) (*RSAKey, error) {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	key := &RSAKey{}
	//读取公钥
	publicKey, err := loadPublicKeyKeyFromFile(path, pu)
	if err != nil {
		return nil, err
	}
	//读取私钥
	privateKey, err := loadPrivateKeyFromFile(path, pr)
	if err != nil {
		return nil, err
	}
	key.PrivateKey = privateKey
	key.PublicKey = publicKey
	return key, nil
}

func (rsaKey *RSAKey) SaveRSAKeyFile(path string, prFileName string, puFileName string) error {
	//保存私钥
	err := savePrivateKeyToFile(rsaKey.PrivateKey, path, prFileName)
	if err != nil {
		return err
	}
	//保存公钥
	err = savePublicKeyToFile(rsaKey.PublicKey, path, puFileName)
	if err != nil {
		return err
	}
	return nil
}

func (rsaKey *RSAKey) SignPKCS1v15(str string) (string, error) {
	hashed := crypto.SHA256.New()
	hashed.Write([]byte(str))
	digest := hashed.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey.PrivateKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(signature), nil
}

func (rsaKey *RSAKey) VerifyPKCS1v15(str, sign string) error {
	hashed := crypto.SHA256.New()
	hashed.Write([]byte(str))
	digest := hashed.Sum(nil)
	signature, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(rsaKey.PublicKey, crypto.SHA256, digest, signature)
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, path string, filename string) error {
	// 将私钥转换为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 写入文件
	privateFile := file.CreateFile(path, filename)
	if privateFile == nil {
		return nil
	}

	defer privateFile.Close()

	return pem.Encode(privateFile, privateKeyBlock)
}

func savePublicKeyToFile(publicKey *rsa.PublicKey, path string, filename string) error {
	// 将公钥转换为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 写入文件
	publicFile := file.CreateFile(path, filename)
	if publicFile == nil {
		return nil
	}
	defer publicFile.Close()

	return pem.Encode(publicFile, publicKeyBlock)
}

func loadPublicKeyKeyFromFile(path string, filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path + filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to decode public key")
	}
	return rsaKey, nil
}

func loadPrivateKeyFromFile(path string, filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path + filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
