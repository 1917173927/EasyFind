package services

import (
	"errors"
	"os"
	"strings"

	"easyfind/internal/config"
	"easyfind/internal/models"
	"easyfind/pkg/database"

	"gorm.io/gorm"
)

func getTargetType(lostOrFound int) models.ItemType {
	if lostOrFound == 1 {
		return models.TypeLost
	}
	return models.TypeFound
}

func buildCommonQuery(campus, category, status, location string, days int, hasBounty *bool) *gorm.DB {
	// 使用 Preload("Publisher") 预加载发布者信息
	query := database.DB.Model(&models.Item{}).Preload("Publisher", func(db *gorm.DB) *gorm.DB {
		// 这里可以根据需要排除敏感字段，例如密码
		return db.Select("id, username, name, nickname, avatar, phone, role, is_active, first_login")
	})

	if campus != "" && strings.ToLower(campus) != "all" {
		query = query.Where("campus = ?", campus)
	}
	if category != "" && strings.ToLower(category) != "all" {
		query = query.Where("category = ?", category)
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if days > 0 {
		query = query.Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days)
	}
	if hasBounty != nil {
		if *hasBounty {
			query = query.Where("bounty > 0")
		} else {
			query = query.Where("bounty = 0")
		}
	}

	if status != "" && strings.ToLower(status) != "all" {
		query = query.Where("status = ?", status)
	} else {
		// 默认只显示已审核通过的，除非明确指定了其他状态，或者后台管理接口可能会覆盖这个逻辑
		// 这里改为：如果 status 为空，默认 status = 'approved'。
		// 但为了兼容原来的逻辑（IsBounty修改前是 buildCommonQuery(campus, category, string(models.StatusApproved))）
		// 我们需要小心。调用者一般会传入 status。
		// 如果调用者传空字符串，以前是 "!= Cancelled AND != Archived"。
		// 现在的需求是"按物品状态筛选"。
		// 假如用户没传 status，默认应该展示 approved。
		// 如果用户传了 status，就按 status 查。
		// 之前的调用者传了 `string(models.StatusApproved)`。
		// 所以如果 status 是 "approved"，就是 status = 'approved'。

		// 保持原有逻辑：如果 status 不为空，就用 status。
		// 如果 status 为空，保留原来的排除逻辑？
		// 原来 logic: if status != "" -> Where("status = ?", status) else -> Where("status != ? ...")
		query = query.Where("status != ? AND status != ?", models.StatusCancelled, models.StatusArchived)
	}

	return query
}

func GetAllItem(campus, category, location string, days int, status string, hasBounty *bool, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	// 如果 status 为空，默认只查 approved
	if status == "" {
		status = string(models.StatusApproved)
	}
	query := buildCommonQuery(campus, category, status, location, days, hasBounty)

	result := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetRecord(campus, category string, lostOrFound int, location string, days int, status string, hasBounty *bool, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	if status == "" {
		status = string(models.StatusApproved)
	}
	query := buildCommonQuery(campus, category, status, location, days, hasBounty)

	if lostOrFound != 0 {
		query = query.Where("type = ?", getTargetType(lostOrFound))
	}

	result := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetAllLostAndFoundTotalPageNum(campus, category, location string, days int, status string, hasBounty *bool) (*int64, error) {
	var pageNum int64
	if status == "" {
		status = string(models.StatusApproved)
	}
	query := buildCommonQuery(campus, category, status, location, days, hasBounty)

	result := query.Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetTotalPageNum(campus, category string, lostOrFound int, location string, days int, status string, hasBounty *bool) (*int64, error) {
	var pageNum int64
	if status == "" {
		status = string(models.StatusApproved)
	}
	query := buildCommonQuery(campus, category, status, location, days, hasBounty)

	if lostOrFound != 0 {
		query = query.Where("type = ?", getTargetType(lostOrFound))
	}

	result := query.Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetRecordByAdmin(campus, category string, lostOrFound int, status string, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	query := database.DB.Model(&models.Item{})

	if campus != "" && strings.ToLower(campus) != "all" {
		query = query.Where("campus = ?", campus)
	}
	if category != "" && strings.ToLower(category) != "all" {
		query = query.Where("category = ?", category)
	}

	query = query.Where("type = ?", getTargetType(lostOrFound))

	if status != "" && strings.ToLower(status) != "all" {
		query = query.Where("status = ?", status)
	}

	result := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetTotalPageNumByAdmin(campus, category string, lostOrFound int, status string) (*int64, error) {
	var pageNum int64
	query := database.DB.Model(&models.Item{})

	if campus != "" && strings.ToLower(campus) != "all" {
		query = query.Where("campus = ?", campus)
	}
	if category != "" && strings.ToLower(category) != "all" {
		query = query.Where("category = ?", category)
	}

	query = query.Where("type = ?", getTargetType(lostOrFound))

	if status != "" && strings.ToLower(status) != "all" {
		query = query.Where("status = ?", status)
	}

	result := query.Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetPendingRecordByAdmin(lostOrFound, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	result := database.DB.Where(models.Item{
		Type: getTargetType(lostOrFound),
	}).Where("status = ?", models.StatusPending).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetPendingTotalPageNumByAdmin(lostOrFound int) (*int64, error) {
	var pageNum int64
	result := database.DB.Model(models.Item{}).Where(models.Item{
		Type: getTargetType(lostOrFound),
	}).Where("status = ?", models.StatusPending).
		Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetAllMyRecord(publisherID uint, status models.ItemStatus, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	query := database.DB.Model(&models.Item{}).Where(models.Item{
		PublisherID: publisherID,
	})

	if status != "" && strings.ToLower(string(status)) != "all" {
		query = query.Where("status = ?", status)
	}
	result := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)

	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetMyTotalPageNum(publisherID uint, status models.ItemStatus) (*int64, error) {
	var pageNum int64
	query := database.DB.Model(&models.Item{}).Where(models.Item{
		PublisherID: publisherID,
	})

	if status != "" && strings.ToLower(string(status)) != "all" {
		query = query.Where("status = ?", status)
	}
	result := query.Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetCategoryList() ([]models.LostCategory, error) {
	var categories []models.LostCategory
	result := database.DB.Where(models.LostCategory{}).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

func GetRecordById(id uint) (models.Item, error) {
	var record models.Item
	result := database.DB.First(&record, id)
	if result.Error != nil {
		return models.Item{}, result.Error
	}
	return record, nil
}

func CreateRecord(record models.Item) error {
	result := database.DB.Create(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func RemoveImg(record models.Item, img1 string, img2 string, img3 string, img4 string) {
	if record.Img1 != "" && record.Img1 != img1 && record.Img1 != img2 && record.Img1 != img3 && record.Img1 != img4 {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img1, config.GetWebpUrlKey()))
	}
	if record.Img2 != "" && record.Img2 != img1 && record.Img2 != img2 && record.Img2 != img3 && record.Img2 != img4 {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img2, config.GetWebpUrlKey()))
	}
	if record.Img3 != "" && record.Img3 != img1 && record.Img3 != img2 && record.Img3 != img3 && record.Img3 != img4 {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img3, config.GetWebpUrlKey()))
	}
	if record.Img4 != "" && record.Img4 != img1 && record.Img4 != img2 && record.Img4 != img3 && record.Img4 != img4 {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img4, config.GetWebpUrlKey()))
	}
}

func UpdateRecord(id int, record models.Item) error {
	result := database.DB.Model(&models.Item{}).
		Omit("id", "type", "created_at", "publisher", "publisher_id", "status").
		Where("id = ?", id).Updates(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteRecord(id uint) error {
	record, err := GetRecordById(id)
	if record.Status == models.StatusApproved {
		return errors.New("approved record cannot be deleted")
	}
	if err != nil {
		return err
	}
	if record.Img1 != "" {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img1, config.GetWebpUrlKey()))
	}
	if record.Img2 != "" {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img2, config.GetWebpUrlKey()))
	}
	if record.Img3 != "" {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img3, config.GetWebpUrlKey()))
	}
	if record.Img4 != "" {
		_ = os.Remove("./img/" + strings.TrimPrefix(record.Img4, config.GetWebpUrlKey()))
	}
	return database.DB.Unscoped().Delete(&models.Item{}, id).Error
}

func CancelRecord(id uint) error {
	record, err := GetRecordById(id)
	if record.Status == models.StatusPending {
		return errors.New("pending record cannot be cancelled")
	}
	if err != nil {
		return err
	}
	result := database.DB.Model(models.Item{}).Where("id = ?", id).
		Update("status", models.StatusCancelled)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ApproveRecord(id uint) error {
	result := database.DB.Model(models.Item{}).Where("id = ?", id).
		Update("status", models.StatusApproved)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func RejectRecord(id uint, rejectReason string) error {
	result := database.DB.Model(models.Item{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        models.StatusRejected,
			"reject_reason": rejectReason,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ArchiveRecord(id uint, processMethod string) error {
	result := database.DB.Model(&models.Item{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         models.StatusArchived,
			"process_method": processMethod,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreatClaim(claim models.Claim) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&claim).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Item{}).Where("id = ?", claim.ItemID).
			Update("status", models.StatusMatched).Error; err != nil {
			return err
		}

		return nil
	})
}

func GetMyClaim(ClaimantID uint, pageNum, pageSize int) ([]models.Claim, error) {
	var claims []models.Claim
	result := database.DB.Where("claimant_id = ?", ClaimantID).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&claims)
	if result.Error != nil {
		return nil, result.Error
	}
	return claims, nil
}

func GetMyClaimTotalPageNum(ClaimantID uint) (*int64, error) {
	var pageNum int64
	result := database.DB.Model(&models.Claim{}).Where("claimant_id = ?", ClaimantID).
		Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetPendingClaimByAdmin(pageNum, pageSize int) ([]models.Claim, error) {
	var claims []models.Claim
	result := database.DB.Model(&models.Claim{}).
		Where("status = ?", "Pending"). // 尝试匹配大写 Pending (GORM 默认值可能是大写)
		Preload("Item").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&claims)

	// 如果查不到大写的，再试试小写的 (兼容性处理)
	if len(claims) == 0 {
		database.DB.Model(&models.Claim{}).
			Where("status = ?", "pending").
			Preload("Item").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Order("created_at desc").Find(&claims)
	}

	if result.Error != nil {
		return nil, result.Error
	}
	return claims, nil
}

func GetPendingClaimTotalPageNumByAdmin() (*int64, error) {
	var pageNum int64
	result := database.DB.Model(&models.Claim{}).
		Where("status = ?", string(models.StatusPending)).
		Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetClaimByID(id uint) (models.Claim, error) {
	var claim models.Claim
	result := database.DB.First(&claim, id)
	if result.Error != nil {
		return models.Claim{}, result.Error
	}
	return claim, nil
}

func ApproveClaim(id uint) error {
	result := database.DB.Model(models.Claim{}).Where("id = ?", id).
		Update("status", models.StatusApproved)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func RejectClaim(id uint) error {
	result := database.DB.Model(models.Claim{}).Where("id = ?", id).
		Update("status", models.StatusRejected)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ConfirmClaim(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var claim models.Claim
		if err := tx.First(&claim, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&claim).Update("status", models.StatusArchived).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Item{}).Where("id = ?", claim.ItemID).
			Update("status", models.StatusClaimed).Error; err != nil {
			return err
		}

		return nil
	})
}
