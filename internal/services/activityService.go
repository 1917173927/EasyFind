package services

import (
	"easyfind/internal/models"
	"easyfind/pkg/database"
	"sort"
)

// GetPostActivities 获取帖子动态
func GetPostActivities(userID uint, page, size int) ([]models.PostActivity, int64, error) {
	var activities []models.PostActivity
	var totalActivities []models.PostActivity

	// 1. 获取用户发布的所有物品 Item 状态
	var items []models.Item
	if err := database.DB.Where("publisher_id = ?", userID).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	for _, item := range items {
		// 类型1: 状态更新 (Approved, Pending, Matched, Claimed)
		if item.Status != models.StatusRejected {
			act := models.PostActivity{
				Type:       "status_update",
				Title:      "你的帖子状态已更新",
				Time:       item.UpdatedAt, // 使用最后更新时间
				ItemID:     item.ID,
				PeerUserID: 0,
				ItemName:   item.Title,
				LossTime:   item.Time,
				Location:   item.Location,
				Img:        item.Img1,
				Status:     string(item.Status),
			}
			totalActivities = append(totalActivities, act)
		}

		// 类型2: 被驳回 (Rejected)
		if item.Status == models.StatusRejected {
			act := models.PostActivity{
				Type:         "rejected",
				Title:        "你的帖子被驳回!",
				Time:         item.UpdatedAt,
				ItemID:       item.ID,
				PeerUserID:   0,
				ItemName:     item.Title,
				LossTime:     item.Time,
				Location:     item.Location,
				Img:          item.Img1,
				Status:       string(item.Status),
				RejectReason: item.RejectReason,
			}
			totalActivities = append(totalActivities, act)
		}
	}

	// 2. 获取用户发布物品的认领记录 Claim (Someone Claimed)
	// 先找到所有 ItemID
	var itemIDs []uint
	for _, item := range items {
		itemIDs = append(itemIDs, item.ID)
	}

	if len(itemIDs) > 0 {
		var claims []models.Claim
		// 查询关联这些 ItemID 的 Claim
		if err := database.DB.Where("item_id IN ?", itemIDs).Preload("Claimant").Preload("Item").Find(&claims).Error; err != nil {
			return nil, 0, err
		}

		for _, claim := range claims {
			// 类型3: 有人认领 (Claim Receive)
			act := models.PostActivity{
				Type:         "claim_received",
				Title:        "你的帖子有人认领",
				Time:         claim.CreatedAt, // 认领创建时间
				ItemID:       claim.ItemID,
				PeerUserID:   claim.ClaimantID,
				ItemName:     claim.Item.Title,
				LossTime:     claim.Item.Time,
				Location:     claim.Item.Location,
				Img:          claim.Item.Img1,
				Status:       string(claim.Item.Status),
				ClaimID:      claim.ID,
				ClaimantName: claim.Claimant.Name, // Assuming Account has Name
			}
			totalActivities = append(totalActivities, act)
		}
	}

	// 3. 排序 (按时间倒序)
	sort.Slice(totalActivities, func(i, j int) bool {
		return totalActivities[i].Time.After(totalActivities[j].Time)
	})

	// 4. 分页
	total := int64(len(totalActivities))
	start := (page - 1) * size
	end := start + size

	if start >= len(totalActivities) {
		return []models.PostActivity{}, total, nil
	}
	if end > len(totalActivities) {
		end = len(totalActivities)
	}

	activities = totalActivities[start:end]

	return activities, total, nil
}
