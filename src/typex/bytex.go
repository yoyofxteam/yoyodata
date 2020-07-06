package typex

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"errors"
	"io/ioutil"
)

// GzippedText is a []byte which transparently gzips data being submitted to
// a database and ungzips data being Scanned from a database.
type GzippedText []byte

// Value implements the driver.Valuer interface, gzipping the raw value of
// this GzippedText.
func (g GzippedText) Value() (driver.Value, error) {
	b := make([]byte, 0, len(g))
	buf := bytes.NewBuffer(b)
	w := gzip.NewWriter(buf)
	w.Write(g)
	w.Close()
	return buf.Bytes(), nil

}

// Scan implements the sql.Scanner interface, ungzipping the value coming off
// the wire and storing the raw result in the GzippedText.
func (g *GzippedText) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("Incompatible type for GzippedText")
	}
	reader, err := gzip.NewReader(bytes.NewReader(source))
	if err != nil {
		return err
	}
	defer reader.Close()
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	*g = GzippedText(b)
	return nil
}