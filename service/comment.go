/* Date:10/4/23-2023
 * Author:fanzhao
 */
package service

import (
	"context"
	"simple-comment-sys/model"
	"strconv"
	"strings"
)

func appendMid(mids string, addstr []string) string {
	var resstr = mids
	for _, v := range addstr {
		if v != "0" {
			resstr = strings.Trim(resstr+","+v, ",")
		}
	}
	return resstr
}
func (s *Service) AddComment(c context.Context, req *model.AddCommentReq) (err error) {
	if req.CommentID == 0 {
		//批注
		var (
			a = new(api.Comment)
		)
		util.SimpleCopyProperties(a, req)
		if _, err1 := s.cmsClient.CreateComment(c, a); err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.CreateComment req:%v", a)
			return
		}
		mentions, err1 := util.SplitInt64(req.Mentions)
		if err1 != nil {
			err = errors.Wrapf(ecode.SystemError, "SplitInt64 error:%v", err1)
			return
		}
		s.sendCommentsNotice(c, req.SetID, req.ProjectID, req.Mid, mentions)
	} else {
		//回复
		var (
			a = new(api.CommentReply)
		)
		util.SimpleCopyProperties(a, req)
		if _, err1 := s.cmsClient.CreateCommentReply(c, a); err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.CreateCommentReply req:%v", a)
			return
		}
		mentions, err1 := util.SplitInt64(req.Mentions)
		if err1 != nil {
			err = errors.Wrapf(ecode.SystemError, "SplitInt64 error:%v", err1)
			return
		}
		s.sendReplyNotice(c, req.SetID, req.Mid, req.ToMid, mentions)
	}
	return
}

func (s *Service) DelComment(c context.Context, req *model.UpdateCommentReq) (err error) {
	var (
		toMid   int64
		content string
	)
	if req.CommentID != 0 {
		//批注
		var (
			a = new(api.UpdateCommentReq)
		)
		comment, err1 := s.cmsClient.GetComment(c, &api.CommentID{CommentID: req.CommentID, SetID: req.SetID})
		if err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.GetComment req:%v", a)
			return
		}
		toMid = comment.Mid
		content = comment.Content
		a.Comment = &api.Comment{
			CommentID: req.CommentID,
			SetID:     req.SetID,
			State:     -100,
		}
		a.UpdateMask = []string{"comment.state"}
		if _, err1 := s.cmsClient.UpdateComment(c, a); err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.UpdateComment req:%v", a)
			return
		}
	} else {
		//回复
		var (
			a = new(api.UpdateCommentReplyReq)
		)
		comment, err1 := s.cmsClient.GetCommentReply(c, &api.CommentReplyID{CrID: req.CrID, SetID: req.SetID})
		if err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.GetCommentReply req:%v", a)
			return
		}
		toMid = comment.Mid
		content = comment.Content
		a.CommentReply = &api.CommentReply{
			CrID:  req.CrID,
			SetID: req.SetID,
			State: -100,
		}
		a.UpdateMask = []string{"commentreply.state"}
		if _, err1 := s.cmsClient.UpdateCommentReply(c, a); err1 != nil {
			err = errors.Wrapf(err1, "cmsClient.UpdateCommentReply req:%v", a)
			return
		}
	}
	s.sendDelCommentNotice(c, req.SetID, req.Mid, toMid, content)
	return
}

func (s *Service) SetCommentState(c context.Context, req *model.UpdateCommentReq) (err error) {
	//批注
	var (
		a = new(api.UpdateCommentReq)
	)
	a.Comment = &api.Comment{
		CommentID: req.CommentID,
		SetID:     req.SetID,
		State:     req.State,
	}
	a.UpdateMask = []string{"comment.state"}
	if _, err1 := s.cmsClient.UpdateComment(c, a); err1 != nil {
		err = errors.Wrapf(err1, "cmsClient.UpdateComment req:%v", a)
		return
	}
	return
}

func (s *Service) ListComments(c context.Context, req *model.ListCommentsReq) (res *model.ListCommentsResp, err error) {
	c = context.Background()
	res = new(model.ListCommentsResp)
	var (
		a               = new(api.ListCommentsReq)
		midstr          string
		mids            = make([]int64, 0)
		projectUserTags = make(map[int64]string)
	)
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		userList, err1 := s.cmsClient.ListProjectUser(c, &api.ListProjectUserReq{
			ProjectId: req.ProjectID,
			State:     constant.VALID,
		})
		if err1 != nil {
			log.Errorc(c, "cmsClient.ListProjectUser error:%v", err1)
			return
		}
		for _, v := range userList.Users {
			projectUserTags[v.Mid] = v.Tag
		}
	}()
	util.SimpleCopyProperties(a, req)
	a.WithReplies = true
	result, err1 := s.cmsClient.ListComments(c, a)
	if err1 != nil {
		err = errors.Wrapf(err1, "cmsClient.ListComments req:%v", a)
		return
	}
	for _, v := range result.Comments {
		tmp := new(model.Comment)
		util.SimpleCopyProperties(tmp, v)
		midstr = appendMid(midstr, []string{v.Mentions,
			strconv.FormatInt(v.Mid, 10)})
		for _, r := range v.CommentReplys {
			tr := new(model.CommentReply)
			util.SimpleCopyProperties(tr, r)
			tmp.Replies = append(tmp.Replies, tr)
			midstr = appendMid(midstr, []string{r.Mentions,
				strconv.FormatInt(r.Mid, 10),
				strconv.FormatInt(r.ToMid, 10)})
		}
		res.List = append(res.List, tmp)
	}

	//补充用户信息
	mids, err = util.SplitInt64(midstr)
	if err != nil {
		err = errors.Wrapf(ecode.SystemError, "SplitInt64 error:%v", err)
		return
	}
	if mids != nil && len(mids) != 0 {

		infoReply, err1 := s.accountClient.Infos3(c, &account.MidsReq{Mids: mids})
		if err1 != nil {
			log.Errorc(c, "accountClient.Infos3 req:%v error:%v", &account.MidsReq{Mids: mids}, err1)
			err = ecode.ThirdPartyServiceErr
			return
		}
		s.wg.Wait()
		res.UserList = make(map[int64]*model.ProjectUserDTO, 0)
		for _, u := range infoReply.Infos {
			ut := new(model.ProjectUserDTO)
			ut.Name = u.GetName()
			ut.Avatar = u.GetFace()
			ut.Tag = projectUserTags[u.Mid]
			res.UserList[u.GetMid()] = ut
		}
	}

	return
}

func (s *Service) ExportComments(c context.Context, req *export.CommentExportMsg) (err error) {
	err = s.dao.SendCommentMsg(c, req.ExportID, req)
	if err != nil {
		return
	}
	_, err = s.cmsClient.SetCommentExportStateCache(c, &api.CacheKV{Key: req.ExportID, Value: strconv.FormatInt(export.EXPORT_PENDING, 10)})
	if err != nil {
		err = errors.Wrapf(err, "cmsClient.SetCommentExportStateCache key:%v,value:%v", req.ExportID, strconv.FormatInt(export.EXPORT_PENDING, 10))
		return
	}
	s.sendExportCommentsNotice(c, req.SetID, req.ProjectID, req.Mid)
	return
}

func (s *Service) GetExportCommentsFile(c context.Context, req *model.GetExportCommentReq) (res *model.GetExportCommentResp, err error) {
	res = new(model.GetExportCommentResp)
	var cache *api.CacheKV
	cache, err = s.cmsClient.GetCommentExportStateCache(c, &api.CacheKV{Key: req.ExportID})
	if err != nil {
		return
	}
	if cache.Value == strconv.FormatInt(export.EXPORT_SUCC, 10) {
		res.Download, err = s.OSS.GetDownloadUrl(&boss.GetDownloadUrlArg{Key: req.ExportID, Expire: 60 * 60 * 24})
		if err != nil {
			err = errors.Wrapf(err, "OSS.GetDownloadUrl key:%v", req.ExportID)
			return
		}
	}
	return
}

func (s *Service) sendCommentsNotice(c context.Context, setID, projectID, mid int64, mentions []int64) {
	var (
		setTitle        string
		projectUserList []int64
		publishName     string
	)
	//get notice info
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		setres, err1 := s.fillSetInfo(c, setID)
		if err1 != nil {
			log.Errorc(c, "s.fillSetInfo error:%v", err1)
			return
		}
		setTitle = setres.Title
	}()
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		userres, err1 := s.ListProjectUser(c, &model.ListProjectUserReq{projectID})
		if err1 != nil {
			log.Errorc(c, "s.ListProjectUser error:%v", err1)
			return
		}
		for _, v := range userres {
			if v.Mid != mid {
				//不发给自己
				projectUserList = append(projectUserList, v.Mid)
			}
		}
	}()
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		infores, err1 := s.accountClient.Info3(c, &account.MidReq{Mid: mid})
		if err1 != nil {
			log.Errorc(c, "accountClient.Info3 error:%v", err1)
			return
		}
		publishName = infores.GetInfo().GetName()
	}()

	//send message
	s.wg.Wait()
	// 项目中成员接收消息
	notices := new(model.MessagePushDataBusSubModel)
	for _, v := range projectUserList {
		notice := &model.MessageModuleCommentFile{
			FileName:    setTitle,
			PublishName: publishName,
			FileId:      setID,
			PublishId:   mid,
			ReceiveId:   v,
		}
		notices.List = append(notices.List, notice.FormatMessageData())
	}
	// 被@的人接收消息
	for _, v := range mentions {
		notice := &model.MessageModuleAt{
			PublishName: publishName,
			FileId:      setID,
			PublishId:   mid,
			ReceiveId:   v,
		}
		notices.List = append(notices.List, notice.FormatMessageData())
	}
	s.MessagePushService(c, notices)
}
func (s *Service) sendReplyNotice(c context.Context, setID, mid, toMid int64, mentions []int64) {
	var (
		publishName string
	)
	if mentions != nil && len(mentions) > 0 {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			infores, err1 := s.accountClient.Info3(c, &account.MidReq{Mid: mid})
			if err1 != nil {
				log.Errorc(c, "accountClient.Info3 error:%v", err1)
				return
			}
			publishName = infores.GetInfo().GetName()
		}()
	}
	//get notice info
	infores, err1 := s.accountClient.Info3(c, &account.MidReq{Mid: mid})
	if err1 != nil {
		log.Errorc(c, "accountClient.Info3 error:%v", err1)
		return
	}
	s.wg.Wait()
	//send message
	//被回复对象消息
	notices := new(model.MessagePushDataBusSubModel)
	notice := &model.MessageModuleReplyFile{
		MemberName: infores.GetInfo().GetName(),
		ReceiveId:  toMid,
		FileId:     setID,
	}
	notices.List = append(notices.List, notice.FormatMessageData())
	// 被@的人接收消息
	for _, v := range mentions {
		notice := &model.MessageModuleAt{
			PublishName: publishName,
			FileId:      setID,
			PublishId:   mid,
			ReceiveId:   v,
		}
		notices.List = append(notices.List, notice.FormatMessageData())
	}
	s.MessagePushService(c, notices)
}
func (s *Service) sendDelCommentNotice(c context.Context, setID, mid, toMid int64, content string) {
	//get notice info
	infores, err1 := s.accountClient.Info3(c, &account.MidReq{Mid: mid})
	if err1 != nil {
		log.Errorc(c, "accountClient.Info3 error:%v", err1)
		return
	}
	//send message
	notices := new(model.MessagePushDataBusSubModel)
	notice := &model.MessageModuleDeleteComment{
		CommentContent: content,
		DelMemberName:  infores.GetInfo().GetName(),
		ReceiveId:      toMid,
		FileId:         setID,
	}
	notices.List = append(notices.List, notice.FormatMessageData())
	s.MessagePushService(c, notices)

}
func (s *Service) sendExportCommentsNotice(c context.Context, setID, projectID, mid int64) {
	var (
		setTitle        string
		projectUserList = make([]int64, 0)
		publishName     string
	)
	//get notice info
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		setres, err1 := s.fillSetInfo(c, setID)
		if err1 != nil {
			log.Errorc(c, "s.fillSetInfo error:%v", err1)
			return
		}
		setTitle = setres.Title
	}()
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		userres, err1 := s.ListProjectUser(c, &model.ListProjectUserReq{projectID})
		if err1 != nil {
			log.Errorc(c, "s.ListProjectUser error:%v", err1)
			return
		}
		for _, v := range userres {
			if v.Role == constant.ADMIN {
				//只发给管理员
				projectUserList = append(projectUserList, v.Mid)
			}
		}
	}()
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		infores, err1 := s.accountClient.Info3(c, &account.MidReq{Mid: mid})
		if err1 != nil {
			log.Errorc(c, "accountClient.Info3 error:%v", err1)
			return
		}
		publishName = infores.GetInfo().GetName()
	}()

	//send message
	s.wg.Wait()
	notices := new(model.MessagePushDataBusSubModel)
	for _, v := range projectUserList {
		notice := &model.MessageModuleExportComment{
			MemberName: publishName,
			FileName:   setTitle,
			ReceiveId:  v,
			FileId:     setID,
		}
		notices.List = append(notices.List, notice.FormatMessageData())
	}
	s.MessagePushService(c, notices)
}
