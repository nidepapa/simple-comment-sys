/* Date:10/4/23-2023
 * Author:fanzhao
 */
package model

import "time"

type Comment struct {
	ID         int64           `json:"id"`                // 自增ID
	CommentID  int64           `json:"comment_id,string"` // 批注id
	ProjectID  int64           `json:"project_id,string"` // 项目ID
	SetID      int64           `json:"set_id,string"`     // 文件ID
	Mid        int64           `json:"mid"`               // 发表者mid
	Annotation string          `json:"annotation"`        // 批注，json
	Content    string          `json:"content"`           // 评论内容
	Mentions   string          `json:"mentions"`          // @的mid,分割
	State      int64           `json:"state"`             // 批注状态 0 未处理 -100删除
	BeginPoint int64           `json:"begin_point"`       // 视频批注开始时间点 ms
	EndPoint   int64           `json:"end_point"`         // 视频批注结束时间点 ms
	Snapshot   string          `json:"snapshot"`          // 截图
	AnnoPic    string          `json:"anno_pic"`          // 批注图
	Ctime      time.Time       `json:"ctime"`             // 创建时间
	Mtime      time.Time       `json:"mtime"`             // 修改时间
	HasReplies bool            `json:"has_replies"`
	Replies    []*CommentReply `json:"replies"`
}

type CommentReply struct {
	ID        int64     `json:"id"`                // 自增ID
	CrID      int64     `json:"cr_id,string"`      // 回复id
	SetID     int64     `json:"set_id,string"`     // 文件ID
	ProjectID int64     `json:"project_id,string"` // 项目ID
	CommentID int64     `json:"comment_id,string"` // 批注ID
	ParentID  int64     `json:"parent_id,string"`  // 回复的父级评论ID
	Mid       int64     `json:"mid" `              // 发表者mid
	ToMid     int64     `json:"to_mid"`            // 回复对象的mid
	Content   string    `json:"content"`           // 评论内容
	Mentions  string    `json:"mentions"`          // @的mid,分割
	State     int64     `json:"state"`             // 批注状态 0 未处理 -100删除
	Ctime     time.Time `json:"ctime"`             // 创建时间
	Mtime     time.Time `json:"mtime"`             // 修改时间
}

type AddCommentReq struct {
	Mid        int64  `json:"mid"`                                   // 发表者mid
	ProjectID  int64  `json:"project_id,string" validate:"required"` // 项目ID
	SetID      int64  `json:"set_id,string" validate:"required" `    // 文件ID
	ParentID   int64  `json:"parent_id,string"`                      // 回复的父级评论ID
	CommentID  int64  `json:"comment_id,string"`                     // 批注id
	Annotation string `json:"annotation"`                            // 批注，json
	Content    string `json:"content"`                               // 评论内容
	Mentions   string `json:"mentions"`                              // @的mid,分割
	State      int64  `json:"state"`                                 // 批注状态 0 未处理 -100删除
	BeginPoint int64  `json:"begin_point" default:"-1"`              // 视频批注开始时间点 ms
	EndPoint   int64  `json:"end_point" default:"-1"`                // 视频批注结束时间点 ms
	Snapshot   string `json:"snapshot"`                              // 截图
	AnnoPic    string `json:"anno_pic"`                              // 批注图
	ToMid      int64  `json:"to_mid"`                                // 回复对象的mid
}

type UpdateCommentReq struct {
	CommentID int64 `json:"comment_id,string"`                 // 批注id
	CrID      int64 `json:"cr_id,string"`                      // 回复id
	SetID     int64 `json:"set_id,string" validate:"required"` // 文件ID
	State     int64 `json:"state"`                             // 批注状态 0 未处理 -100删除
	Mid       int64 `json:"mid"`                               // 发表者mid
}

type ListCommentsReq struct {
	Page int64 `form:"page" ,json:"page"`
	// default 20
	PageSize int64 `form:"page_size" ,json:"page_size"`
	// order by  for example :  [name] [-id]  -表示：倒序排序
	Orderby []string `form:"orderby" ,json:"orderby"`
	// 过滤条件需要自定义 for example  query name has
	ProjectID int64 `form:"project_id" ,json:"project_id,string" validate:"required"`
	SetID     int64 `form:"set_id" ,json:"set_id,string" validate:"required"`
}

type ProjectUserDTO struct {
	// project-based user info
}

type ListCommentsResp struct {
	List []*Comment `json:"list"`
	//	UserList   []*User    `json:"user_list"`
	UserList   map[int64]*ProjectUserDTO `json:"user_list,string"`
	TotalCount int64                     `json:"total_count,string"`
	PageCount  int64                     `json:"page_count,string"`
}

type GetExportCommentReq struct {
	ExportID string `form:"export_id" json:"export_id" validate:"required"`
}

type GetExportCommentResp struct {
	Download string `json:"download" form:"download" validate:"required"`
}
