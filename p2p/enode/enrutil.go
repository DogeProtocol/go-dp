package enode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/DogeProtocol/dp/p2p/enr"
	"github.com/DogeProtocol/dp/rlp"
	"io"
	"net"
	"strconv"
	"strings"
)

// dumpRecord creates a human-readable description of the given node record.
func DumpRecord(out io.Writer, r *enr.Record) {
	n, err := New(ValidSchemes, r)
	if err != nil {
		fmt.Fprintf(out, "INVALID: %v\n", err)
	} else {
		fmt.Fprintf(out, "Node ID: %v\n", n.ID())
		DumpNodeURL(out, n)
	}
	kv := r.AppendElements(nil)[1:]
	fmt.Fprintf(out, "Record has sequence number %d and %d key/value pairs.\n", r.Seq(), len(kv)/2)
	fmt.Fprint(out, DumpRecordKV(kv, 2))
}

func DumpNodeURL(out io.Writer, n *Node) {
	var key PqPubKey
	if n.Load(&key) != nil {
		return // no secp256k1 public key
	}
	fmt.Fprintf(out, "URLv4:   %s\n", n.URLv4())
}

func DumpRecordKV(kv []interface{}, indent int) string {
	// Determine the longest key name for alignment.
	var out string
	var longestKey = 0
	for i := 0; i < len(kv); i += 2 {
		key := kv[i].(string)
		if len(key) > longestKey {
			longestKey = len(key)
		}
	}
	// Print the keys, invoking formatters for known keys.
	for i := 0; i < len(kv); i += 2 {
		key := kv[i].(string)
		val := kv[i+1].(rlp.RawValue)
		pad := longestKey - len(key)
		out += strings.Repeat(" ", indent) + strconv.Quote(key) + strings.Repeat(" ", pad+1)
		formatter := attrFormatters[key]
		if formatter == nil {
			formatter = formatAttrRaw
		}
		fmtval, ok := formatter(val)
		if ok {
			out += fmtval + "\n"
		} else {
			out += hex.EncodeToString(val) + " (!)\n"
		}
	}
	return out
}

// parseNode parses a node record and verifies its signature.
func ParseNode(source string) (*Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return ParseV4(source)
	}
	r, err := ParseRecord(source)
	if err != nil {
		return nil, err
	}
	return New(ValidSchemes, r)
}

// parseRecord parses a node record from hex, base64, or raw binary input.
func ParseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)
	if d, ok := DecodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := DecodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}
	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)
	return &r, err
}

func DecodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}
	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)
	return dec, err == nil
}

func DecodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}
	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)
	return dec[:n], err == nil
}

// attrFormatters contains formatting functions for well-known ENR keys.
var attrFormatters = map[string]func(rlp.RawValue) (string, bool){
	"id":   formatAttrString,
	"ip":   formatAttrIP,
	"ip6":  formatAttrIP,
	"tcp":  formatAttrUint,
	"tcp6": formatAttrUint,
}

func formatAttrRaw(v rlp.RawValue) (string, bool) {
	s := hex.EncodeToString(v)
	return s, true
}

func formatAttrString(v rlp.RawValue) (string, bool) {
	content, _, err := rlp.SplitString(v)
	return strconv.Quote(string(content)), err == nil
}

func formatAttrIP(v rlp.RawValue) (string, bool) {
	content, _, err := rlp.SplitString(v)
	if err != nil || len(content) != 4 && len(content) != 6 {
		return "", false
	}
	return net.IP(content).String(), true
}

func formatAttrUint(v rlp.RawValue) (string, bool) {
	var x uint64
	if err := rlp.DecodeBytes(v, &x); err != nil {
		return "", false
	}
	return strconv.FormatUint(x, 10), true
}
