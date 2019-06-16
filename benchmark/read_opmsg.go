// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package benchmark

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/x/mongo/driver/description"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/network/command"
	"go.mongodb.org/mongo-driver/x/network/wiremessage"
)

func readOpMsgDecoding(ctx context.Context, tm TimerManager, iters int, dataSet string) error {
	d, err := loadSourceDocument(getProjectRoot(), perfDataDir, bsonDataDir, dataSet)
	if err != nil {
		return err
	}

	s := bson.D{
		{"cursor", bson.D{
			{"firstBatch", bson.A{d}},
			{"id", int64(0)},
			{"ns", "namespace"}},
		},
		{"ok", int64(1)},
	}
	r, _ := bson.Marshal(s)

	wm := wiremessage.Msg{
		Sections: []wiremessage.Section{
			wiremessage.SectionBody{
				PayloadType: wiremessage.SingleDocument,
				Document:    bson.Raw(r),
			},
		},
	}
	rc := command.Read{}
	desc := description.SelectedServer{}

	tm.ResetTimer()

	for i := 0; i < iters; i++ {
		_, err := rc.Decode(desc, wm).Result()
		if err != nil {
			return errors.New("decoding failed")
		}
	}
	return nil
}

func ReadOpMsgFlatDecoding(ctx context.Context, tm TimerManager, iters int) error {
	return readOpMsgDecoding(ctx, tm, iters, flatBSONData)
}

func ReadOpMsgDeepDecoding(ctx context.Context, tm TimerManager, iters int) error {
	return readOpMsgDecoding(ctx, tm, iters, deepBSONData)
}
