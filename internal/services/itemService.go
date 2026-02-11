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

func buildCommonQuery(campus, category, status string) *gorm.DB {
	query := database.DB.Model(&models.Item{})

	if campus != "" && strings.ToLower(campus) != "all" {
		query = query.Where("campus = ?", campus)
	}
	if category != "" && strings.ToLower(category) != "all" {
		query = query.Where("category = ?", category)
	}
	if status != "" && strings.ToLower(status) != "all" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status != ? AND status != ?", models.StatusCancelled, models.StatusArchived)
	}

	return query
}

func GetAllItem(campus, category string, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	query := buildCommonQuery(campus, category, string(models.StatusApproved))

	result := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetRecord(campus, category string, lostOrFound int, pageNum, pageSize int) ([]models.Item, error) {
	var record []models.Item
	query := buildCommonQuery(campus, category, string(models.StatusApproved))

	result := query.Where("type = ?", getTargetType(lostOrFound)).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func GetAllLostAndFoundTotalPageNum(campus, category string) (*int64, error) {
	var pageNum int64
	query := buildCommonQuery(campus, category, string(models.StatusApproved))

	result := query.Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetTotalPageNum(campus, category string, lostOrFound int) (*int64, error) {
	var pageNum int64
	query := buildCommonQuery(campus, category, string(models.StatusApproved))

	result := query.Where("type = ?", getTargetType(lostOrFound)).
		Count(&pageNum)
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
	result := database.DB.Model(models.Item{}).Select("*").
		Omit("id", "type", "created_at", "publisher").
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
	result := database.DB.
		Where("status = ?", models.StatusPending).Preload("Item").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order("created_at desc").Find(&claims)
	if result.Error != nil {
		return nil, result.Error
	}
	return claims, nil
}

func GetPendingClaimTotalPageNumByAdmin() (*int64, error) {
	var pageNum int64
	result := database.DB.Model(&models.Claim{}).
		Where("status = ?", models.StatusPending).
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
