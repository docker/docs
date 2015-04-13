// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#if defined(SQLITE_AMALGAMATION) && defined(SQLITE_HAS_CODEC)

#include "codec.h"

// codec.go exports.
int go_codec_init(const CodecCtx*,void**,char**);
int go_codec_reserve(void*);
void go_codec_resize(void*,int,int);
void *go_codec_exec(void*,void*,u32,int);
void go_codec_get_key(void*,void**,int*);
void go_codec_free(void*);

// sqlite3_key sets the codec key for the main database.
SQLITE_API int sqlite3_key(sqlite3 *db, const void *pKey, int nKey) {
	return sqlite3_key_v2(db, 0, pKey, nKey);
}

// sqlite3_key_v2 sets the codec key for the specified database.
SQLITE_API int sqlite3_key_v2(sqlite3 *db, const char *zDb, const void *pKey, int nKey) {
	int iDb = 0;
	int rc;
	sqlite3_mutex_enter(db->mutex);
	if (zDb && zDb[0]) {
		iDb = sqlite3FindDbName(db, zDb);
	}
	if (iDb < 0) {
		rc = SQLITE_ERROR;
		sqlite3Error(db, SQLITE_ERROR, "unknown database: %s", zDb);
	} else {
		rc = sqlite3CodecAttach(db, iDb, pKey, nKey);
	}
	rc = sqlite3ApiExit(db, rc);
	sqlite3_mutex_leave(db->mutex);
	return rc;
}

// sqlite3_rekey changes the codec key for the main database.
SQLITE_API int sqlite3_rekey(sqlite3 *db, const void *pKey, int nKey) {
	return sqlite3_rekey_v2(db, 0, pKey, nKey);
}

// sqlite3_rekey_v2 changes the codec key for the specified database.
SQLITE_API int sqlite3_rekey_v2(sqlite3 *db, const char *zDb, const void *pKey, int nKey) {
	int iDb = 0;
	int rc;
	sqlite3_mutex_enter(db->mutex);

	rc = SQLITE_ERROR;
	sqlite3Error(db, SQLITE_ERROR, "rekey is not implemented");

	rc = sqlite3ApiExit(db, rc);
	sqlite3_mutex_leave(db->mutex);
	return rc;
}

// sqlite3_activate_see isn't used by Go codecs, but it needs to be linked in.
SQLITE_API void sqlite3_activate_see(const char *zPassPhrase) {}

// sqlite3CodecAttach initializes the codec, reserves space at the end of each
// page, and attaches the codec to the specified database.
int sqlite3CodecAttach(sqlite3 *db, int iDb, const void *pKey, int nKey) {
	Btree *pBt = db->aDb[iDb].pBt;
	Pager *pPager = sqlite3BtreePager(pBt);
	CodecCtx ctx;
	void *pCodec = 0;
	char *zErrMsg = 0;
	int rc;

	// An empty KEY clause in an ATTACH statement disables the codec and SQLite
	// doesn't support codecs for in-memory databases.
	if (nKey <= 0 || pPager->memDb) return SQLITE_OK;

	ctx.db    = db;
	ctx.zPath = sqlite3BtreeGetFilename(pBt);
	ctx.zName = db->aDb[iDb].zName;
	ctx.nBuf  = sqlite3BtreeGetPageSize(pBt);
	ctx.nRes  = sqlite3BtreeGetReserve(pBt);
	ctx.pKey  = pKey;
	ctx.nKey  = nKey;

	sqlite3BtreeEnter(pBt);
	ctx.fixed = (pBt->pBt->btsFlags & BTS_PAGESIZE_FIXED) != 0;
	sqlite3BtreeLeave(pBt);

	if ((rc=go_codec_init(&ctx, &pCodec, &zErrMsg)) != SQLITE_OK) {
		sqlite3Error(db, rc, (zErrMsg ? "%s" : 0), zErrMsg);
		free(zErrMsg);
	} else if (pCodec) {
		int nRes = go_codec_reserve(pCodec);
		if (nRes != ctx.nRes && nRes >= 0) {
			rc = sqlite3BtreeSetPageSize(pBt, -1, nRes, 0);
		}
		if (rc != SQLITE_OK) {
			go_codec_free(pCodec);
			sqlite3Error(db, rc, "unable to reserve page space for the codec");
		} else {
			sqlite3PagerSetCodec(pPager, go_codec_exec, go_codec_resize, go_codec_free, pCodec);
		}
	}
	return rc;
}

// sqlite3CodecGetKey returns the codec key for the specified database.
void sqlite3CodecGetKey(sqlite3 *db, int iDb, void **pKey, int *nKey) {
	void *pCodec = sqlite3PagerGetCodec(sqlite3BtreePager(db->aDb[iDb].pBt));
	*pKey = 0;
	*nKey = 0;
	if (pCodec) {
		go_codec_get_key(pCodec, pKey, nKey);
	}
}

#endif
