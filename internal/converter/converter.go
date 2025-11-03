package converter

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// convertToUTF8 は様々なエンコーディング（JIS、Shift_JIS、EUC-JP、UTF-8）からUTF-8に変換します
func ConvertToUTF8(data []byte) (string, error) {
	// まずUTF-8として試す
	if isValidUTF8(data) {
		return string(data), nil
	}

	// ISO-2022-JP (JIS) として試す
	decoder := japanese.ISO2022JP.NewDecoder()
	decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// Shift_JIS として試す
	decoder = japanese.ShiftJIS.NewDecoder()
	decoded, err = io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// EUC-JP として試す
	decoder = japanese.EUCJP.NewDecoder()
	decoded, err = io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// どのエンコーディングでも失敗した場合、元のデータをそのまま返す
	return string(data), nil
}

// isValidUTF8 はバイト列が有効なUTF-8かどうかをチェックします
func isValidUTF8(data []byte) bool {
	// 文字列に変換して、元のバイト列と比較
	s := string(data)
	return len(s) > 0 && string([]byte(s)) == string(data)
}
