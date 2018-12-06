package lib

import (
	"encoding/base64"
	"fmt"
	"github.com/okeuday/erlang_go/src/erlang"
	"math/big"
	"regexp"
	"sort"
	"strings"
)

type UpdateSequence struct {
	Node  string
	Range []*big.Int
	Seq   int
}

func convertNode(t interface{}) string {
	return string(t.(erlang.OtpErlangAtom))
}

func convertRange(t interface{}) []*big.Int {
	value := t.(erlang.OtpErlangList).Value
	decoded := make([]*big.Int, len(value))
	for i, r := range value {
		switch rType := r.(type) {
		case uint8:
			decoded[i] = big.NewInt(int64(rType))
		case int32:
			decoded[i] = big.NewInt(int64(rType))
		case *big.Int:
			decoded[i] = rType
		default:
			fmt.Printf("%v == %T\n", rType, rType)
			// todo return err
			decoded[i] = big.NewInt(0)
		}
	}
	return decoded
}

func convertSeq(t interface{}) int {
	switch rType := t.(type) {
	case uint8:
		return int(rType)
	case erlang.OtpErlangTuple:
		return int(erlang.OtpErlangTuple(rType)[0].(uint8))
	default:
		fmt.Printf("%v == %T\n", rType, rType)
		// todo return err
		return -1
	}
}

func convert(terms interface{}) []UpdateSequence {
	res := make([]UpdateSequence, 0)
	value := terms.(erlang.OtpErlangList).Value
	for _, term := range value {
		conv := term.(erlang.OtpErlangTuple)
		//log.Printf("c: %v", conv)

		t := UpdateSequence{
			Node:  convertNode(conv[0]),
			Range: convertRange(conv[1]),
			Seq:   convertSeq(conv[2]),
		}

		//log.Printf("converted: %v", t)
		res = append(res, t)
	}
	return res
}

func DecodeUpdateSeq(updateSeq string) ([]UpdateSequence, error) {
	//start := time.Now()
	//defer func() {
	//	fmt.Printf("Duration: %s\n", time.Since(start))
	//}()

	encoded := string(regexp.MustCompile("^\\d+\\-").ReplaceAll([]byte(updateSeq), []byte{}))

	b1 := strings.Replace(encoded, "-", "+", -1)
	b2 := strings.Replace(b1, "_", "/", -1)

	padding := strings.Repeat("=", (4-len(b2)%4)%4)

	var sb strings.Builder
	sb.WriteString(b2)
	sb.WriteString(padding)

	payload := make([]byte, base64.StdEncoding.DecodedLen(len(sb.String())))
	n, err := base64.StdEncoding.Decode(payload, []byte(sb.String()))
	if err != nil {
		return nil, err
	}
	terms, err := erlang.BinaryToTerm(payload[:n])
	if err != nil {
		return nil, err
	}
	decoded := convert(terms)
	sort.Slice(decoded, func(i, j int) bool {
		return decoded[i].Range[0].Cmp(decoded[j].Range[0]) < 0
	})
	return decoded, nil
}
