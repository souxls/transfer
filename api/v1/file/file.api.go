package file

import (
	"fmt"
	"net/http"
	"time"

	"transfer/api/types/request"
	"transfer/api/types/response"
	"transfer/api/types/schema"
	"transfer/internal/app/model"
	"transfer/internal/app/service"

	"transfer/internal/app/pkg/logger"
	"transfer/internal/app/pkg/security"
	mURL "transfer/internal/app/pkg/url"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var FileSet = wire.NewSet(wire.Struct(new(API), "*"))

type API struct {
	FileService        *service.File
	PermesssionService *service.Permession
}

func (a *API) Create(c *gin.Context) {

	logger.Infof("start upload file.")

	upload, err := c.FormFile("file")
	if err != nil {
		logger.Warnf("get formfile  error: %s", err)
		response.RetMsg(c, http.StatusBadRequest, "缺少文件", "")

		return
	}

	claims := jwt.ExtractClaims(c)
	user := claims["id"].(string)

	fileID := security.GetUUID()
	nowTime := time.Now()
	file := &model.File{
		Filename:   upload.Filename,
		Fileid:     fileID,
		Bucketname: user,
		Createtime: nowTime,
		Owner:      user,
		Expired:    nowTime.Add(time.Hour * viper.GetViper().GetDuration("MinIO.FileExpired")), //minio 默认是7天
	}

	r, _ := upload.Open()
	defer r.Close()

	if err := a.FileService.Create(c, file, r); err != nil {
		logger.Errorf("file save failed: %s", err.Error())
		response.RetMsg(c, http.StatusInternalServerError, "文件保存失败", "")
		return
	}

	newFile := a.FileService.ByFileid(c, file)

	if newFile != nil {
		retData := schema.FileURL{
			File: file.Filename,
			URL:  mURL.ShortUrlencode(newFile.ID),
		}

		logger.Infof("%s file upload successd.", newFile.Filename)
		response.RetMsg(c, http.StatusOK, "文件上传成功", retData)
		return
	}
	logger.Infof("file upload failed.")
	response.RetMsg(c, http.StatusOK, "文件上传失败, 存储数据获取失败", "")
}

func (a *API) Get(c *gin.Context) {

	logger.Infof("start download file or get file sign.")

	fileID, idOK := c.Params.Get("id")
	if !idOK {

		logger.Debugf("get id error")
		response.RetMsg(c, http.StatusBadRequest, "请求参数错误", "")
		return
	}

	fileSign := c.Query("AuthParam")
	logger.Debugf("request sign param. %s", fileSign)

	file := &model.File{
		Fileid: fileID,
	}

	fileInfo := a.FileService.Info(c, file)

	if fileInfo != nil {
		// 当有 AuthParam参数的时候，开始下载文件。无 AuthParam 获取 file sign。
		if fileSign != "" {
			if a.FileService.CheckSign(c, fileInfo, fileSign) {

				logger.Infof("start download: %s", fileInfo.Filename)

				extraHeaders := map[string]string{
					"Content-Disposition": fmt.Sprintf("attachment; filename=%s", fileInfo.Filename),
					"response-type":       "blob",
				}
				reader := a.FileService.Download(c, fileInfo.Bucketname, fileInfo.Filename)

				readerInfo, _ := reader.Stat()

				c.DataFromReader(http.StatusOK, readerInfo.Size, "application/octet-stream", reader, extraHeaders)
				logger.Infof("download finish: %s", fileInfo.Filename)

				return
			} else {

				logger.Infof("file sign check faild")
				response.RetMsg(c, http.StatusNoContent, "下载URL过期或不存在", "")

				return
			}
		}

		claims := jwt.ExtractClaims(c)
		user := claims["id"].(string)
		userPermession := &model.Permession{
			Fileid:   fileID,
			Username: user,
		}

		if !a.PermesssionService.CheckFile(c, userPermession) {
			response.RetMsg(c, http.StatusForbidden, "权限不足或已过期, 不允许下载", "")
			return
		}

		// 请求无 AuthParam 返回 filesign
		if sign := a.FileService.Sign(c, fileInfo); sign != nil {
			logger.Infof("file sign is %s", sign)
			response.RetMsg(c, http.StatusOK, "获取成功", sign)

			return
		}

		response.RetMsg(c, http.StatusNoContent, "获取失败", "")
		return
	}

	response.RetMsg(c, http.StatusNotFound, "文件不存在或已过期", "")
}

func (a *API) Delete(c *gin.Context) {

	logger.Infof("start delete file.")

	fileId, ok := c.Params.Get("id")
	if !ok {
		response.RetMsg(c, http.StatusBadRequest, "文件 id 格式错误", "")
		return
	}

	claims := jwt.ExtractClaims(c)
	user := claims["id"].(string)

	file := &model.File{
		Owner:  user,
		Fileid: fileId,
	}

	if err := a.FileService.Delete(c, file); err != nil {
		logger.Errorf("delete file failed.", err.Error())
		response.RetMsg(c, http.StatusInternalServerError, "文件删除失败", "")

		return
	}

	logger.Infof("file delete successd. %s", file.Fileid)
	response.RetMsg(c, http.StatusOK, "删除成功", "")
}

func (a *API) Query(c *gin.Context) {

	logger.Infof("fetch file list.")
	claims := jwt.ExtractClaims(c)
	user := claims["id"].(string)

	var paginationParam schema.PaginationParam
	if c.BindQuery(&paginationParam) != nil {
		logger.Debug("没有分页参数")
	}

	fileList, err := a.FileService.QueryShow(c, &model.File{Owner: user}, paginationParam)
	if err != nil {
		logger.Warnf("get file list failed: %s", err)
		response.RetMsg(c, http.StatusNoContent, "列表为空", "")
		return
	}

	logger.Debugf("file list: %s", fileList)
	response.RetMsg(c, http.StatusOK, "", *fileList)
}

func (a *API) CreateAuth(c *gin.Context) {

	logger.Infof("create auth for file")
	id, _ := c.Params.Get("id")

	var users request.RequestUsers

	if err := c.BindJSON(&users); err != nil {
		response.RetMsg(c, http.StatusBadRequest, "数据格式不正确, 需要数组", "")
		return
	}

	nowTime := time.Now()
	for _, user := range users.Data {
		var file *model.File

		// fileid 是 32位的，小于32位的按短url处理
		if len(id) < 32 {
			file = &model.File{
				ID: mURL.ShortUrldecode(id),
			}
		} else {
			file = &model.File{
				Fileid: id,
			}
		}

		fileInfo := a.FileService.Info(c, file)
		logger.Debugf("file info: %s", fileInfo)
		if fileInfo == nil || time.Now().After(fileInfo.Expired) {
			response.RetMsg(c, http.StatusBadRequest, "文件已过期或不存在", "")
			return
		}

		userPermession := &model.Permession{
			Username:   user,
			Shorturl:   fileInfo.ID,
			Fileid:     fileInfo.Fileid,
			Expiredate: nowTime.Add(time.Minute * viper.GetViper().GetDuration("MinIO.UserExpired")),
		}

		if err := a.PermesssionService.Create(c, userPermession); err != nil {
			response.RetMsg(c, http.StatusInternalServerError, fmt.Sprintf("%s授权失败", user), "")
			return
		}
	}

	response.RetMsg(c, http.StatusOK, "授权成功", "")
}

func (a *API) UpdateAuth(c *gin.Context) {

	logger.Infof("start update auth")
	fileId, _ := c.Params.Get("id")
	user, _ := c.Params.Get("userid")
	postExpireDate := &schema.FileExpireDate{}

	if err := c.BindJSON(postExpireDate); err != nil {
		response.RetMsg(c, http.StatusBadRequest, "请求参数错误", "")
		return
	}
	expireDate, err := time.ParseInLocation("2006-01-02 15:04:05", postExpireDate.Expiredate, time.Local)
	if err != nil {
		response.RetMsg(c, http.StatusBadRequest, "时间格式错误。请使用: 2023-04-13 13:59:30", "")
		return
	}
	logger.Debugf("expiredate: %s", expireDate)

	permession := &model.Permession{
		Fileid:     fileId,
		Username:   user,
		Expiredate: expireDate,
	}

	if !a.PermesssionService.Update(c, permession) {
		msg := fmt.Sprintf("为用户 %s 授权失败", permession.Username)
		response.RetMsg(c, http.StatusInternalServerError, msg, "")
		return
	}
	logger.Infof("update auth finish for %s", permession.Fileid)
	response.RetMsg(c, http.StatusOK, "授权成功", "")
}

func (a *API) ByShortUrl(c *gin.Context) {
	// 根据前端返回的短url在权限表中获取fileid,使用fileid从files 表中获取文件详细信息

	logger.Infof("get file info by short url.")
	shortUrl, _ := c.Params.Get("url")
	claims := jwt.ExtractClaims(c)
	user := claims["id"].(string)

	permession := &model.Permession{
		Shorturl: mURL.ShortUrldecode(shortUrl),
		Username: user,
	}

	if permession := a.PermesssionService.Info(c, permession); permession != nil {

		file := &model.File{
			ID: mURL.ShortUrldecode(shortUrl),
		}

		var fileInfoList schema.Files
		if fileInfo := a.FileService.Info(c, file); fileInfo != nil {
			for _, p := range *permession {
				fileInfoList := append(fileInfoList, &schema.FileInfo{
					Filename: fileInfo.Filename,
					Fileid:   fileInfo.Fileid,
					Owner:    fileInfo.Owner,
					Expired:  p.Expiredate,
				})

				logger.Debugf("file info: %s", fileInfoList)
				response.RetMsg(c, http.StatusOK, "", fileInfoList)
				return
			}
		}
	}
	response.RetMsg(c, http.StatusNotFound, "文件不存在或已过期", "")
}

func (a *API) ForDownload(c *gin.Context) {
	// 获取文件下载列表
	logger.Info("fetch file list for download")
	claims := jwt.ExtractClaims(c)
	user := claims["id"].(string)

	userpemession := &model.Permession{
		Username: user,
	}

	if permession := a.PermesssionService.Info(c, userpemession); permession != nil {
		fileList := a.FileService.ForDownload(c, "", permession)

		logger.Debugf("download file list, permession: %s", permession)
		if fileList != nil {
			logger.Debugf("download file list. %s", fileList)
			response.RetMsg(c, http.StatusOK, "", fileList)
			return
		}
	}
	response.RetMsg(c, http.StatusNoContent, "列表为空", "")

}
