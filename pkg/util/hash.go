// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	// SHA-256 is FIPS-approved and works under GODEBUG=fips140=only.
	// This site is a non-cryptographic deterministic ID hash; the algorithm
	// was migrated from SHA-1 to SHA-256 (truncated to 20 bytes for K8s
	// label-value width compatibility) so the operator can run under strict
	// FIPS mode without panicking. See pkg/util/fips/gate.go and
	// docs/security_hardening.md §3.
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash/fnv"

	dumper "github.com/sanity-io/litter"
)

func serializeUnrepeatable(obj interface{}) []byte {
	b := bytes.Buffer{}
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(obj)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}

	return b.Bytes()
}

func serializeRepeatable(obj interface{}) []byte {
	d := dumper.Options{
		Separator: " ",
	}
	return []byte(d.Sdump(obj))
}

// HashIntoString returns a deterministic 40-char hex digest used as a
// non-cryptographic object fingerprint / K8s label value (see Fingerprint and
// labeler.MakeObjectVersion). NOT a security control: the digest is only used
// to compare two serialized object representations for equality. The 40-char
// width is part of the contract — K8s label values must be ≤63 chars and the
// `clickhouse.altinity.com/object-version` label is compared verbatim across
// reconciles, so any change to width or algorithm forces a one-time STS roll
// on operator upgrade as every existing object's label value re-hashes.
//
// SHA-256 (truncated to 20 bytes / 40 hex chars) is used so the operator
// works under GODEBUG=fips140=only — SHA-1 panics in strict mode.
func HashIntoString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:20])
}

// HashIntoInt hashes bytes and returns int version of the hash
func HashIntoInt(b []byte) int {
	if len(b) == 0 {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write(b)
	return int(h.Sum32())
}

// HashIntoIntTopped hashes bytes and return int version of the ash topped with top
func HashIntoIntTopped(b []byte, top int) int {
	return HashIntoInt(b) % top
}
