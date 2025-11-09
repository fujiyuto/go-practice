package jwt

import (
	"fmt"
	"time"
	"encoding/base64"
	"encoding/json"
	"crypto/sha256"

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

	// 発行日時を設定し、有効に設定
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

	return fmt.Sprintf("%s.%s\n", header64Encode, payload64Encode), nil
}