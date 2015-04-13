// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _CODEC_H_
#define _CODEC_H_

// Codec initialization context.
typedef struct {
	sqlite3 *db;
	const char *zPath;
	const char *zName;
	int nBuf;
	int nRes;
	int fixed;
	const void *pKey;
	int nKey;
} CodecCtx;

// SQLite codec hooks.
int sqlite3CodecAttach(sqlite3*,int,const void*,int);
void sqlite3CodecGetKey(sqlite3*,int,void**,int*);

#endif
