package controllers

import (
	"easyfind/internal/apiErr"
	"easyfind/internal/services"
	"easyfind/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
}

// UploadImage godoc
// @Summary 上传图片
// @Description 上传图片文件，支持自动去重，返回图片 URL
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response{data=UploadResponse} "上传成功"
// @Failure 200 {object} response.Response "失败"
// @Router /api/v1/upload/image [post]
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		apiErr.HandleBizError(c, response.InvalidParams, "file is required")
		return
	}

	// 调用 Service
	image, err := services.UploadServiceApp.UploadImage(file)
	if err != nil {
		apiErr.HandleSysError(c, response.CodeError, err)
		return
	}

	response.Success(c, UploadResponse{
		URL:      image.URL,
		FileName: image.FileName,
	})
}
