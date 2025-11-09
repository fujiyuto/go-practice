package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"os"
	"time"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/afero"
)

type Header struct {
	Alg string
	Typ string
}

type Payload struct {
	Sub    string
	Name   string
	Iat    int64
	Valid  bool
}

type Jwt struct {
	Header
	Payload
}

func Run(fn string, g bool, a bool) error {
	if g {
		jwtToken, err := generateJwt(fn)
		if err != nil {
			return err
		}

		fmt.Println(jwtToken)
	}

	return nil
}

// JWTの生成
func generateJwt(fn string) (string, error) {
	// JWTの生成はファイル読み込みで行う

	// ファイルシステム初期化
	fs := afero.NewOsFs()

	// jwt.jsonファイル読み込み
	file, err := fs.Open(fn)
	if err != nil {
		fmt.Println("Error: json file could not find.")
	}

	// jsonデータ格納変数
	var jwtData Jwt

	// jsonデコード
	if err := json.NewDecoder(file).Decode(&jwtData); err != nil {
		return "", err
	}
	
	// ヘッダー
	Header := &Header{
		Alg: jwtData.Header.Alg,
		Typ: jwtData.Header.Typ,
	}

	if Header.Alg == "" {
		Header.Alg = "HS256"
	}

	if Header.Typ != "JWT" {
		Header.Typ = "JWT"
	}

	// ペイロード
	Payload := &Payload{
		Sub:  jwtData.Payload.Sub,
		Name: jwtData.Payload.Name,

	}

	// 発行日時(タイムスタンプ)を設定し、有効に設定
	Payload.Iat = time.Now().Unix()
	Payload.Valid = true

	// ヘッダー、ペイロードをjsonエンコードして、base64エンコード
	headerJsonEncode, err := json.Marshal(Header)
	if err != nil {
		return "", err
	}

	header64Encode := base64.StdEncoding.EncodeToString(headerJsonEncode)

	payloadJsonEncode, err := json.Marshal(Payload)
	if err != nil {
		return "", err
	}

	payload64Encode := base64.StdEncoding.EncodeToString(payloadJsonEncode)

	// 署名生成
	signature, err := generateSignature(header64Encode, payload64Encode)

	return strings.Join([]string{header64Encode, payload64Encode, signature}, "."), nil
}

// 署名(signature)の生成
func generateSignature(header string, payload string) (string, error) {
	// .envファイルの読み込み
	err := godotenv.Load(".env")
	
	if err != nil {
		return "", err
	}

	// get secret key
	key := os.Getenv("JWT_SECRET_KEY")

	// create signature
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(strings.Join([]string{header, payload}, ".")))
	signature := mac.Sum(nil)

	return hex.EncodeToString(signature), nil
}