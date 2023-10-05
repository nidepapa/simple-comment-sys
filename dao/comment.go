/* Date:10/4/23-2023
 * Author:fanzhao
 */
package dao

import (
	"context"
	"fmt"
	comment "simple-comment-sys/model"
	"strings"
)

func KeyCommentExport(id string) string {
	return id
}

func KeyCommentID(id int64) string {
	return fmt.Sprintf("comment_reply_%d", id)
}

// CreateComment CreateComment
func (d *dao) CreateComment(ctx context.Context, m *comment.Comment) (rs *comment.Comment, err error) {
	//
	_, err = comment.
		Create(d.db).Debug().
		Table(comment.TableName(m.SetID)).
		SetComment(m).
		Save(ctx)
	if err != nil {
		errors.Wrapf(err, "create error")
		return nil, err
	}
	// query after create and return
	rs, err = comment.
		Find(d.db.Master()).
		Table(comment.TableName(m.SetID)).
		Where(
			comment.CommentIDEQ(m.CommentID),
		).
		One(ctx)
	return rs, nil
}

func (d *dao) UpdateComment(ctx context.Context, req *api.UpdateCommentReq) (rs *comment.Comment, err error) {
	update := comment.Update(d.db).Table(comment.TableName(req.GetComment().GetSetID()))

	for _, v := range req.GetUpdateMask() {
		switch v {
		case "comment.comment_id":
			update.SetCommentID(req.GetComment().GetCommentID())
		case "comment.mid":
			update.SetMid(req.GetComment().GetMid())
		case "comment.project_id":
			update.SetProjectID(req.GetComment().GetProjectID())
		case "comment.file_id":
			update.SetSetID(req.GetComment().GetSetID())
		case "comment.annotation":
			update.SetAnnotation(req.GetComment().GetAnnotation())
		case "comment.content":
			update.SetContent(req.GetComment().GetContent())
		case "comment.mentions":
			update.SetMentions(req.GetComment().GetMentions())
		case "comment.state":
			update.SetState(req.GetComment().GetState())
		case "comment.begin_point":
			update.SetBeginPoint(req.GetComment().GetBeginPoint())
		case "comment.end_point":
			update.SetEndPoint(req.GetComment().GetEndPoint())
		case "comment.snapshot":
			update.SetSnapshot(req.GetComment().GetSnapshot())
		case "comment.anno_pic":
			update.SetAnnoPic(req.GetComment().GetAnnoPic())
		}
	}
	_, err = update.
		Where(
			comment.CommentIDEQ(req.GetComment().GetCommentID()),
		).
		Save(ctx)
	if err != nil {
		errors.Wrapf(err, "update error:%s", update.Sql())
		return nil, err
	}
	// query after update and return
	rs, err = comment.
		Find(d.db.Master()).
		Table(comment.TableName(req.GetComment().GetSetID())).
		Where(
			comment.CommentIDEQ(req.GetComment().GetCommentID()),
		).
		One(ctx)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (d *dao) GetComment(ctx context.Context, m *comment.Comment) (rs *comment.Comment, err error) {
	return comment.Find(d.db).Debug().Table(comment.TableName(m.SetID)).Where(comment.CommentIDEQ(m.CommentID), comment.NotDeleted()).One(ctx)
}

func (d *dao) ListComments(ctx context.Context, setid, projectid, offset, size int64, orderby []string) (list []*comment.Comment, total int64, err error) {
	find := comment.
		Find(d.db).
		Table(comment.TableName(setid)).
		Offset(int(offset)).
		Limit(int(size))
	for _, v := range orderby {
		if strings.HasPrefix(v, "-") {
			find.OrderDesc(strings.TrimPrefix(v, "-"))
			continue
		}
		find.OrderAsc(v)
	}
	// costom filter
	{
		find.Where(
			comment.SetIDEQ(setid), comment.ProjectIDEQ(projectid),
			comment.NotDeleted(),
		)
	}
	list, err = find.All(ctx)
	if err != nil {
		errors.Wrapf(err, "list error:%s", find.Sql())
		return
	}
	total, err = comment.
		Find(d.db).
		Table(comment.TableName(setid)).
		Count().
		Int64(ctx)
	if err != nil {
		errors.Wrapf(err, "count total error")
		return
	}
	return
}

func (d *dao) SetCacheCommentExportState(c context.Context, key string, value string) (err error) {
	return errors.Wrapf(d.AddCacheCommentExportState(c, key, value), "redis set error key %s", key)
}

func (d *dao) GetCacheCommentExportState(c context.Context, key string) (res string, err error) {
	res, err = d.CacheCommentExportState(c, key)
	err = errors.Wrapf(err, "redis get error key %s", key)
	return
}
