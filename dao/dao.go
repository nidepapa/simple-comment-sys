/* Date:10/4/23-2023
 * Author:fanzhao
 */
package dao

import (
	"context"
	"database/sql"
	"os/user"
	comment "simple-comment-sys/model"
	"time"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

//go:generate kratos tool btsgen
type _bts interface {
	// bts:  -cache_err=continue
	FileInfo(c context.Context, keys []int64) (map[int64][]*fileinfo.FileInfo, error)
}

//go:generate kratos tool redisgen
type _redis interface {
	// redis: -key=fileInfo -struct_name=dao -batch=50 -max_group=5 -batch_err=continue
	CacheFileInfo(c context.Context, keys []int64) (map[int64][]*fileinfo.FileInfo, error)
	// redis: -key=fileInfo -expire=d.demoExpire -encode=json -struct_name=dao
	AddCacheFileInfo(c context.Context, values map[int64][]*fileinfo.FileInfo) error
	// redis: -key=fileInfo -struct_name=dao
	DelCacheFileInfo(c context.Context, keys []int64) error
	// redis: -key=KeyCommentExport -struct_name=dao
	CacheCommentExportState(c context.Context, key string) (string, error)
	// redis:  -key=KeyCommentExport -expire=d.commentExportExpire -struct_name=dao
	AddCacheCommentExportState(c context.Context, key string, value string) error
	// redis: -key=KeyCommentID -struct_name=dao
	CacheCommentReplies(c context.Context, key int64) ([]*commentreply.CommentReply, error)
	// redis:  -key=KeyCommentID -expire=d.commentReplyExpire -struct_name=dao -encode=json
	AddCacheCommentReplies(c context.Context, key int64, value []*commentreply.CommentReply) error
	// redis: -key=KeyCommentID -struct_name=dao
	DelCacheCommentReplies(c context.Context, keys []int64) error
}

// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)

	CreateComment(ctx context.Context, m *comment.Comment) (rs *comment.Comment, err error)
	UpdateComment(ctx context.Context, req *api.UpdateCommentReq) (rs *comment.Comment, err error)
	GetComment(ctx context.Context, m *comment.Comment) (rs *comment.Comment, err error)
	ListComments(ctx context.Context, setid, projectid, offset, size int64, orderby []string) (list []*comment.Comment, total int64, err error)
	CreateCommentReply(ctx context.Context, m *commentreply.CommentReply) (rs *commentreply.CommentReply, err error)
	UpdateCommentReply(ctx context.Context, req *api.UpdateCommentReplyReq) (rs *commentreply.CommentReply, err error)
	GetCommentReply(ctx context.Context, m *commentreply.CommentReply) (rs *commentreply.CommentReply, err error)
	ListCommentReplys(ctx context.Context, commentid, setid, offset, size int64, orderby []string) (list []*commentreply.CommentReply, total int64, err error)
	SetCacheCommentExportState(c context.Context, key string, value string) (err error)
	GetCacheCommentExportState(c context.Context, key string) (res string, err error)
	SetCacheCommentReplies(c context.Context, key int64, value []*commentreply.CommentReply) (err error)
	GetCacheCommentReplies(c context.Context, key int64) (res []*commentreply.CommentReply, err error)

	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
}

// dao dao.
type dao struct {
	db                  *sql.DB
	redis               *redis.Redis
	mc                  *memcache.Memcache
	cache               *fanout.Fanout
	demoExpire          int32
	commentExportExpire int64
	commentReplyExpire  int64
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire          xtime.Duration
		CommentExportExpire xtime.Duration
		CommentReplyExpire  xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db:                  db,
		redis:               r,
		mc:                  mc,
		cache:               fanout.New("cache"),
		demoExpire:          int32(time.Duration(cfg.DemoExpire) / time.Second),
		commentExportExpire: int64(time.Duration(cfg.CommentExportExpire) / time.Second),
		commentReplyExpire:  int64(time.Duration(cfg.CommentExportExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}

func (d *dao) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = d.db.Begin(ctx)
	if err != nil {
		return
	}
	return
}
