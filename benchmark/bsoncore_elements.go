// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package benchmark

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/bsontype"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func readDocumentUsingElementsMethod(ctx context.Context, tm TimerManager, iters int, dataSet string) error {
	raw, err := loadSourceDocument(getProjectRoot(), perfDataDir, bsonDataDir, dataSet)
	if err != nil {
		return err
	}
	r, _ := bson.Marshal(raw)
	d := bsoncore.Document(r)

	tm.ResetTimer()

	for i := 0; i < iters; i++ {
		innerReadDocumentUsingElementsMethod(d)
	}
	return nil
}

func innerReadDocumentUsingElementsMethod(d bsoncore.Document) error {
	elems, err := d.Elements()
	if err != nil {
		return err
	}
	for _, elem := range elems {
		if elem.Value().Type == bsontype.EmbeddedDocument {
			innerReadDocumentUsingElementsMethod(elem.Value().Document())
		}
	}
	return nil
}

func readDocumentUsingReadElementFunction(ctx context.Context, tm TimerManager, iters int, dataSet string) error {
	raw, err := loadSourceDocument(getProjectRoot(), perfDataDir, bsonDataDir, dataSet)
	if err != nil {
		return err
	}
	r, _ := bson.Marshal(raw)
	d := bsoncore.Document(r)

	tm.ResetTimer()

	for i := 0; i < iters; i++ {
		innerReadDocumentUsingReadElementFunction(d)
	}
	return nil
}

func innerReadDocumentUsingReadElementFunction(d bsoncore.Document) error {
	length, rem, ok := bsoncore.ReadLength(d)
	if !ok {
		return bsoncore.NewInsufficientBytesError(d, rem)
	}

	length -= 4

	var elem bsoncore.Element
	for length > 1 {
		elem, rem, ok = bsoncore.ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return bsoncore.NewInsufficientBytesError(d, rem)
		}
		if err := elem.Validate(); err != nil {
			return err
		}

		if elem.Value().Type == bsontype.EmbeddedDocument {
			innerReadDocumentUsingReadElementFunction(elem.Value().Document())
		}
	}
	return nil
}

func ReadDocumentUsingElementsMethodFlat(ctx context.Context, tm TimerManager, iters int) error {
	return readDocumentUsingElementsMethod(ctx, tm, iters, flatBSONData)
}

func ReadDocumentUsingElementsMethodDeep(ctx context.Context, tm TimerManager, iters int) error {
	return readDocumentUsingElementsMethod(ctx, tm, iters, deepBSONData)
}

func ReadDocumentUsingReadElementFunctionFlat(ctx context.Context, tm TimerManager, iters int) error {
	return readDocumentUsingReadElementFunction(ctx, tm, iters, flatBSONData)
}

func ReadDocumentUsingReadElementFunctionDeep(ctx context.Context, tm TimerManager, iters int) error {
	return readDocumentUsingReadElementFunction(ctx, tm, iters, deepBSONData)
}
