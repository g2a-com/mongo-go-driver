// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package command

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/network/wiremessage"
)

func decodeCommandOpMsg(msg wiremessage.Msg) (bson.Raw, error) {
	var mainDoc bsoncore.Document
	idx, mainDoc := bsoncore.ReserveLength(mainDoc)

	for _, section := range msg.Sections {
		switch converted := section.(type) {
		case wiremessage.SectionBody:
			doc := bsoncore.Document(converted.Document)
			elems, err := doc.Elements()
			if err != nil {
				return nil, err
			}
			for _, elem := range elems {
				mainDoc = bsoncore.AppendHeader(mainDoc, elem.Value().Type, elem.Key())
				mainDoc = append(mainDoc, elem.Value().Data...)
			}
		case wiremessage.SectionDocumentSequence:
			mainDoc = bsoncore.AppendHeader(mainDoc, bsontype.Array, converted.Identifier)
			idx, mainDoc := bsoncore.ReserveLength(mainDoc)
			for i, doc := range converted.Documents {
				val := bsoncore.Value{Type: bsontype.Type(doc[0]), Data: doc[2:]}
				mainDoc = bsoncore.AppendHeader(mainDoc, val.Type, strconv.Itoa(i))
				mainDoc = append(mainDoc, val.Data...)
			}
			mainDoc = append(mainDoc, 0x00)
			mainDoc = bsoncore.UpdateLength(mainDoc, idx, int32(len(mainDoc[idx:])))
		}
	}
	mainDoc, err := bsoncore.AppendDocumentEnd(mainDoc, idx)
	if err != nil {
		return nil, err
	}

	rdr := bson.Raw(mainDoc)
	err = rdr.Validate()
	if err != nil {
		return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
	}

	err = extractError(rdr)
	if err != nil {
		return nil, err
	}
	return rdr, nil
}
