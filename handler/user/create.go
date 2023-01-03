package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tzh.com/web/handler"
	"tzh.com/web/model"
	"tzh.com/web/pkg/errno"
	"tzh.com/web/util"
)

// @Summary 创建用户
// @Description 插入新用户到数据库中
// @Tags user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data": {}}"
// @Router /user [post]
func Create(ctx *gin.Context) {
	logrus.WithField(
		"X-Request-Id", util.GetReqID(ctx),
	).Info("用户创建函数被调用")
	// 将 request body 绑定到一个结构体中
	var r CreateRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		handler.SendResponse(ctx, errno.New(errno.ErrBind, err), nil)
		return
	}
	logrus.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// 验证结构
	if err := u.Validate(); err != nil {
		handler.SendResponse(ctx, errno.New(errno.ErrValidation, err), nil)
		return
	}

	// 加密密码
	if err := u.Encrypt(); err != nil {
		handler.SendResponse(ctx, errno.New(errno.ErrEncrypt, err), nil)
		return
	}

	// 插入用户到数据库中
	if err := u.Create(); err != nil {
		handler.SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	resp := CreateResponse{
		Username: r.Username,
	}
	handler.SendResponse(ctx, nil, resp)
}
