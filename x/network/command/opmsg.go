// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package command

import (
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
				mainDoc = bsoncore.AppendValueElement(mainDoc, elem.Key(), elem.Value())
			}
		case wiremessage.SectionDocumentSequence:
			docs := make([]bsoncore.Value, len(converted.Documents))
			for i, doc := range converted.Documents {
				docs[i] = bsoncore.Value{Type: bsontype.Type(doc[0]), Data: doc[2:]}
			}
			mainDoc = bsoncore.BuildArrayElement(mainDoc, converted.Identifier, docs...)
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

	return rdr, extractError(rdr)
}
