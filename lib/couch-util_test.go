package lib

import (
	"encoding/json"
	"log"
	"testing"
)

func TestDecode1(t *testing.T) {
	// example from https://lists.apache.org/thread.html/fb5be24cfc9df28c788ab28973bc09828f2d0c4380d8594921cdf619@%3Cuser.couchdb.apache.org%3E
	decoded, err := DecodeUpdateSeq("g1AAAAI-eJyd0EsOgjAQBuAGiI-dN9C9LmrBwqzkJtrSNkgQV6z1JnoTvYneBEvbhA0aMU1mkj6-_NMSITTJfYFm2anOcsFT10mpTzyG-LxpmiL32eqoN8aEAcWE9dz_jPCFrnzrHGQchiFM4kSgaV0JqQ6VFF-AtAV2DggMgCEGxrNhQfatc3bOyDiKUalg2EBVoCu66KapazcUh41e69-GssjNIvcWWRokk2oNofwj0MNazy4QFURhGQ0J9LKI-SHPIBHEgiak51nxBhxnrRk")
	if err != nil {
		t.Error(err)
	}

	log.Printf("decoded:\n%v\n", decoded)

	numberOfShards := 8
	if len(decoded) != numberOfShards {
		b, _ := json.Marshal(decoded)
		t.Errorf("expected %s to contain %d tuples", b, numberOfShards)
	}
}

func TestDecode2(t *testing.T) {
	decoded, err := DecodeUpdateSeq("g1AAAAFzeJzLYWBg4MhgTmEQTc4vTc5ISXIwNDfSMzTTMzK20DM0zAFKMyUyJMn___8_K5GBgMIkBSCZZA9Wy0hIrQNIbTxx5iaA1NYTpTaPBUgyNAApoPL5xKpfAFG_n1j1ByDq7xOr_gFEPcj9WQAWgmVX")
	if err != nil {
		t.Error(err)
	}

	log.Printf("decoded:\n%v\n", decoded)

	numberOfShards := 8
	if len(decoded) != numberOfShards {
		b, _ := json.Marshal(decoded)
		t.Errorf("expected %s to contain %d tuples", b, numberOfShards)
	}
}
