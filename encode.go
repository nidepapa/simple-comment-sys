/* Date:10/4/23-2023
 * Author:fanzhao
 */
package simple_comment_sys

var _ ecode.Codes

var (
	// 5000001 - 5000999
	SystemError             = ecode.Error(5000001, "系统错误")
	ReqParamErr             = ecode.Error(5000002, "参数错误")
	CmsNoData               = ecode.Error(5000003, "啥都没找到") // cms查询不存在
	RequestInvalid          = ecode.Error(5000004, "请求参数校验错误")
	OpInvalid               = ecode.Error(5000005, "操作非法") // 操作非法
	ThirdPartyServiceErr    = ecode.Error(5000006, "第三方服务错误")
	TOKEN_EXPIRED           = ecode.Error(5000007, "分享链接过期")
	SHARE_NEED_PWD          = ecode.Error(5000008, "分享需要输入密码")
	SHARE_PWD_ERR           = ecode.Error(5000009, "分享密码错误")
	NotPowerUp              = ecode.Error(5000010, "非万粉up")    // 非万粉up
	ProjectSizeNotEnough    = ecode.Error(5000011, "项目容量不足")   // 项目容量不足
	UserSizeNotEnough       = ecode.Error(5000012, "用户容量不足")   // 用户容量不足
	ProjectReachMemberLimit = ecode.Error(5000013, "项目人数已达上限") // 项目人数已达上限
	SHARE_BANNED            = ecode.Error(5000014, "分享链接已禁用")
	VERSION_REACH_LIMIT     = ecode.Error(5000015, "版本数量已达上限")
	SET_STATE_INVALID       = ecode.Error(5000016, "文件未完成，暂不能操作")
	OpItemsLimit            = ecode.Error(5000017, "操作文件数量太多")
	CommentRepliesMiss      = ecode.Error(5000020, "评论回复查询失败了~")

	// oss服务 5105001 - 5105099
	ImportFileSizeOverMax = ecode.Error(5105001, "文件过大")
)


