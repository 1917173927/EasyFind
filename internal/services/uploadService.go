package services

import (
	"crypto/md5"
	"easyfind/internal/models"
	"easyfind/pkg/database"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type UploadService struct{}

var UploadServiceApp = new(UploadService)

// UploadImage 处理图片上传 (去重)
func (s *UploadService) UploadImage(file *multipart.FileHeader) (*models.Image, error) {
	// 1. 打开文件读取内容
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 2. 计算 MD5 (去重)
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, err
	}
	md5String := hex.EncodeToString(hash.Sum(nil))

	// 3. 检查数据库是否存在
	var existingImage models.Image
	err = database.DB.Where("hash = ?", md5String).First(&existingImage).Error
	if err == nil {
		// 存在则直接返回
		return &existingImage, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. 不存在，保存文件
	// 重置文件指针
	if _, err := src.Seek(0, 0); err != nil {
		return nil, err
	}

	// 确保目录存在
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 生成新文件名: md5_timestamp.ext
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s_%d%s", md5String, time.Now().Unix(), ext)
	dstPath := filepath.Join(uploadDir, newFileName)

	// 创建目标文件
	out, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, err
	}

	// 5. 写入数据库
	// URL 路径取决于 Static 服务配置，这里假设是 /uploads/images/xxx
	// windows 下 filepath.Join 会用反斜杠，URL 需要正斜杠
	url := "/" + filepath.ToSlash(dstPath)

	newImage := models.Image{
		FileName: file.Filename,
		Hash:     md5String,
		URL:      url,     // 前端访问路径
		Path:     dstPath, // 物理路径
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}

	if err := database.DB.Create(&newImage).Error; err != nil {
		// 如果入库失败，是否删除文件？通常建议删除，但这里暂不处理
		return nil, err
	}

	return &newImage, nil
}
