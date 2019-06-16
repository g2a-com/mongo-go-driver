// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package benchmark

import "testing"

func BenchmarkReadDocumentUsingElementsMethodFlat(b *testing.B) {
	WrapCase(ReadDocumentUsingElementsMethodFlat)(b)
}
func BenchmarkReadDocumentUsingElementsMethodDeep(b *testing.B) {
	WrapCase(ReadDocumentUsingElementsMethodDeep)(b)
}
func BenchmarkReadDocumentUsingReadElementFunctionFlat(b *testing.B) {
	WrapCase(ReadDocumentUsingReadElementFunctionFlat)(b)
}
func BenchmarkReadDocumentUsingReadElementFunctionDeep(b *testing.B) {
	WrapCase(ReadDocumentUsingReadElementFunctionDeep)(b)
}
