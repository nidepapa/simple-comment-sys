/* Date:10/4/23-2023
 * Author:fanzhao
 */
package dao

import (
	"context"
	"strings"
)

func (d *dao) CreateCommentReply(ctx context.Context, m *commentreply.CommentReply) (rs *commentreply.CommentReply, err error) {
	//
	_, err = commentreply.
		Create(d.db).Debug().
		Table(commentreply.TableName(m.SetID)).
		SetCommentReply(m).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after create and return
	rs, err = commentreply.
		Find(d.db.Master()).
		Table(commentreply.TableName(m.SetID)).
		Where(
			commentreply.CrIDEQ(m.CrID),
		).
		One(ctx)
	return rs, nil
}

func (d *dao) UpdateCommentReply(ctx context.Context, req *api.UpdateCommentReplyReq) (rs *commentreply.CommentReply, err error) {
	update := commentreply.Update(d.db).Table(commentreply.TableName(req.GetCommentReply().GetSetID()))

	for _, v := range req.GetUpdateMask() {
		switch v {
		case "commentreply.cr_id":
			update.SetCrID(req.GetCommentReply().GetCrID())
		case "commentreply.mid":
			update.SetMid(req.GetCommentReply().GetMid())
		case "commentreply.to_mid":
			update.SetToMid(req.GetCommentReply().GetToMid())
		case "commentreply.project_id":
			update.SetProjectID(req.GetCommentReply().GetProjectID())
		case "commentreply.comment_id":
			update.SetCommentID(req.GetCommentReply().GetCommentID())
		case "commentreply.parent_id":
			update.SetParentID(req.GetCommentReply().GetParentID())
		case "commentreply.content":
			update.SetContent(req.GetCommentReply().GetContent())
		case "commentreply.mentions":
			update.SetMentions(req.GetCommentReply().GetMentions())
		case "commentreply.state":
			update.SetState(req.GetCommentReply().GetState())
		}
	}
	_, err = update.
		Where(
			commentreply.CrIDEQ(req.GetCommentReply().GetCrID()),
		).
		Save(ctx)
	if err != nil {
		errors.Wrapf(err, "update error:%s", update.Sql())
		return nil, err
	}
	// query after update and return
	rs, err = commentreply.
		Find(d.db.Master()).
		Table(commentreply.TableName(req.GetCommentReply().GetSetID())).
		Where(
			commentreply.CrIDEQ(req.GetCommentReply().GetCrID()),
		).
		One(ctx)
	return rs, nil
}

func (d *dao) GetCommentReply(ctx context.Context, m *commentreply.CommentReply) (rs *commentreply.CommentReply, err error) {
	return commentreply.Find(d.db).Debug().Table(commentreply.TableName(m.SetID)).Where(commentreply.CrIDEQ(m.CrID), commentreply.NotDeleted()).One(ctx)
}

func (d *dao) ListCommentReplys(ctx context.Context, commentid, setid, offset, size int64, orderby []string) (list []*commentreply.CommentReply, total int64, err error) {
	find := commentreply.
		Find(d.db).
		Table(commentreply.TableName(setid)).
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
			commentreply.CommentIDEQ(commentid),
			commentreply.NotDeleted(),
		)
	}
	list, err = find.All(ctx)
	if err != nil {
		errors.Wrapf(err, "list error:%s", find.Sql())
		return
	}
	total, err = commentreply.
		Find(d.db).
		Table(commentreply.TableName(setid)).
		Count().
		Int64(ctx)
	if err != nil {
		return
	}
	return
}

func (d *dao) SetCacheCommentReplies(c context.Context, key int64, value []*commentreply.CommentReply) (err error) {
	return errors.Wrapf(d.AddCacheCommentReplies(c, key, value), "redis set error key %s", key)
}

func (d *dao) GetCacheCommentReplies(c context.Context, key int64) (res []*commentreply.CommentReply, err error) {
	res, err = d.CacheCommentReplies(c, key)
	err = errors.Wrapf(err, "redis get error key %s", key)
	return
}
