package services

import (
	"math"
	"ongi-back/database"
	"ongi-back/models"
	"sort"
)

type UserSimilarity struct {
	User       models.User `json:"user"`
	Similarity float64     `json:"similarity"`
}

// 유클리드 거리 기반 유사도 계산
func calculateSimilarity(profile1, profile2 *models.UserProfile) float64 {
	diff1 := profile1.SocialityScore - profile2.SocialityScore
	diff2 := profile1.ActivityScore - profile2.ActivityScore
	diff3 := profile1.IntimacyScore - profile2.IntimacyScore
	diff4 := profile1.ImmersionScore - profile2.ImmersionScore
	diff5 := profile1.FlexibilityScore - profile2.FlexibilityScore

	distance := math.Sqrt(
		diff1*diff1 + diff2*diff2 + diff3*diff3 + diff4*diff4 + diff5*diff5,
	)

	// 거리를 유사도로 변환 (0-100)
	maxDistance := math.Sqrt(5 * 100 * 100) // 최대 거리
	similarity := (1 - (distance / maxDistance)) * 100

	return similarity
}

func GetSimilarUsers(userID uint, limit int) ([]UserSimilarity, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var allProfiles []models.UserProfile
	err = database.DB.Preload("User").
		Where("user_id != ?", userID).
		Find(&allProfiles).Error
	if err != nil {
		return nil, err
	}

	// 유사도 계산
	similarities := []UserSimilarity{}
	for _, profile := range allProfiles {
		similarity := calculateSimilarity(&userProfile, &profile)
		similarities = append(similarities, UserSimilarity{
			User:       profile.User,
			Similarity: similarity,
		})
	}

	// 70% 이상 유사도만 필터링
	var filteredSimilarities []UserSimilarity
	for _, sim := range similarities {
		if sim.Similarity >= 70.0 {
			filteredSimilarities = append(filteredSimilarities, sim)
		}
	}

	// 유사도 높은 순으로 정렬
	sort.Slice(filteredSimilarities, func(i, j int) bool {
		return filteredSimilarities[i].Similarity > filteredSimilarities[j].Similarity
	})

	// 상위 N명 반환
	maxResults := limit
	if len(filteredSimilarities) < maxResults {
		maxResults = len(filteredSimilarities)
	}

	return filteredSimilarities[:maxResults], nil
}

func GetRecommendedClubs(userID uint, limit int) ([]models.Club, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var clubs []models.Club
	query := database.DB.Preload("Members")

	// 사교성이 높은 사람에게는 멤버가 많은 클럽 추천
	if userProfile.SocialityScore >= 70 {
		query = query.Order("member_count DESC")
	} else {
		// 친밀도가 높은 사람에게는 적당한 규모의 클럽 추천
		query = query.Order("member_count ASC")
	}

	err = query.Limit(limit).Find(&clubs).Error
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

func GetClubsWithSimilarMembers(userID uint, limit int) ([]models.Club, error) {
	// 유사한 사용자들이 많이 가입한 클럽 찾기
	similarUsers, err := GetSimilarUsers(userID, 20)
	if err != nil {
		return nil, err
	}

	if len(similarUsers) == 0 {
		return GetRecommendedClubs(userID, limit)
	}

	var userIDs []uint
	for _, userSim := range similarUsers {
		userIDs = append(userIDs, userSim.User.ID)
	}

	type ClubCount struct {
		ClubID uint
		Count  int64
	}

	var clubCounts []ClubCount
	err = database.DB.Model(&models.ClubMember{}).
		Select("club_id, COUNT(*) as count").
		Where("user_id IN ?", userIDs).
		Group("club_id").
		Order("count DESC").
		Limit(limit).
		Scan(&clubCounts).Error

	if err != nil {
		return nil, err
	}

	var clubIDs []uint
	for _, cc := range clubCounts {
		clubIDs = append(clubIDs, cc.ClubID)
	}

	var clubs []models.Club
	if len(clubIDs) > 0 {
		err = database.DB.Where("id IN ?", clubIDs).Find(&clubs).Error
		if err != nil {
			return nil, err
		}
	}

	return clubs, nil
}

func GetRecommendedMeetings(userID uint, limit int) ([]models.Meeting, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var meetings []models.Meeting
	query := database.DB.Preload("Club")

	// 활동성이 높은 사람에게는 다양한 모임 추천
	if userProfile.ActivityScore >= 70 {
		query = query.Order("scheduled_at ASC")
	} else {
		query = query.Order("max_members ASC")
	}

	err = query.Limit(limit).Find(&meetings).Error
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

type UserGroup struct {
	Users      []models.User
	AvgProfile models.UserProfile
}

// 비슷한 성향의 사용자들을 그룹화하고 적합한 클럽에 매칭
func MatchUsersToClubs() error {
	// 1. 프로필이 있는 모든 사용자 가져오기
	var profiles []models.UserProfile
	err := database.DB.Preload("User").Find(&profiles).Error
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		return nil
	}

	// 2. 사용자들을 유사도 기반으로 그룹화
	groups := groupSimilarUsers(profiles, 70.0) // 70% 이상 유사도

	// 3. 각 그룹에 적합한 클럽 찾기
	var clubs []models.Club
	err = database.DB.Preload("Members").Find(&clubs).Error
	if err != nil {
		return err
	}

	// 4. 각 그룹을 클럽에 매칭
	for _, group := range groups {
		if len(group.Users) == 0 {
			continue
		}

		// 그룹의 평균 성향과 가장 잘 맞는 클럽 찾기
		bestClub := findBestClubForGroup(group, clubs)
		if bestClub == nil {
			continue
		}

		// 그룹의 모든 사용자를 해당 클럽에 추가
		addedCount := 0
		for _, user := range group.Users {
			// 이미 가입했는지 확인
			var existingMember models.ClubMember
			err := database.DB.Where("club_id = ? AND user_id = ?", bestClub.ID, user.ID).
				First(&existingMember).Error

			if err != nil { // 가입하지 않은 경우만 추가
				member := models.ClubMember{
					ClubID: bestClub.ID,
					UserID: user.ID,
				}
				database.DB.Create(&member)
				addedCount++
			}
		}

		// 실제 추가된 멤버 수만큼 증가
		if addedCount > 0 {
			database.DB.Model(&models.Club{}).
				Where("id = ?", bestClub.ID).
				Update("member_count", database.DB.Raw("member_count + ?", addedCount))
		}
	}

	return nil
}

// 유사한 사용자들을 그룹화
func groupSimilarUsers(profiles []models.UserProfile, threshold float64) []UserGroup {
	var groups []UserGroup
	used := make(map[uint]bool)

	for i, profile1 := range profiles {
		if used[profile1.UserID] {
			continue
		}

		// 새 그룹 생성
		group := UserGroup{
			Users: []models.User{profile1.User},
		}
		used[profile1.UserID] = true

		// 유사한 사용자들 찾기
		for j, profile2 := range profiles {
			if i == j || used[profile2.UserID] {
				continue
			}

			similarity := calculateSimilarity(&profile1, &profile2)
			if similarity >= threshold {
				group.Users = append(group.Users, profile2.User)
				used[profile2.UserID] = true
			}
		}

		// 그룹의 평균 프로필 계산
		group.AvgProfile = calculateGroupAverage(profiles, group.Users)
		groups = append(groups, group)
	}

	return groups
}

// 그룹의 평균 성향 계산
func calculateGroupAverage(allProfiles []models.UserProfile, groupUsers []models.User) models.UserProfile {
	var sumSociality, sumActivity, sumIntimacy, sumImmersion, sumFlexibility float64
	count := 0

	userIDMap := make(map[uint]bool)
	for _, user := range groupUsers {
		userIDMap[user.ID] = true
	}

	for _, profile := range allProfiles {
		if userIDMap[profile.UserID] {
			sumSociality += profile.SocialityScore
			sumActivity += profile.ActivityScore
			sumIntimacy += profile.IntimacyScore
			sumImmersion += profile.ImmersionScore
			sumFlexibility += profile.FlexibilityScore
			count++
		}
	}

	if count == 0 {
		return models.UserProfile{}
	}

	return models.UserProfile{
		SocialityScore:   sumSociality / float64(count),
		ActivityScore:    sumActivity / float64(count),
		IntimacyScore:    sumIntimacy / float64(count),
		ImmersionScore:   sumImmersion / float64(count),
		FlexibilityScore: sumFlexibility / float64(count),
	}
}

// 그룹에 가장 적합한 클럽 찾기
func findBestClubForGroup(group UserGroup, clubs []models.Club) *models.Club {
	if len(clubs) == 0 {
		return nil
	}

	var bestClub *models.Club
	bestScore := -1.0

	for i := range clubs {
		club := &clubs[i]

		// 멤버 수 체크 (그룹 전체가 들어갈 수 있는지)
		if club.MaxMembers > 0 && club.MemberCount+len(group.Users) > club.MaxMembers {
			continue
		}

		// 클럽의 선호 성향과 그룹 평균 성향의 유사도 계산
		score := calculateClubMatchScore(group.AvgProfile, club)

		if score > bestScore {
			bestScore = score
			bestClub = club
		}
	}

	return bestClub
}

// 클럽과 그룹의 매칭 점수 계산
func calculateClubMatchScore(avgProfile models.UserProfile, club *models.Club) float64 {
	// Vibe에 따른 성향 매칭
	score := 0.0

	switch club.Vibe {
	case "energetic":
		score += avgProfile.SocialityScore * 0.3
		score += avgProfile.ActivityScore * 0.4
		score += avgProfile.FlexibilityScore * 0.3
	case "cozy":
		score += avgProfile.IntimacyScore * 0.4
		score += avgProfile.ImmersionScore * 0.3
		score += (100 - avgProfile.ActivityScore) * 0.3
	case "deep":
		score += avgProfile.ImmersionScore * 0.5
		score += avgProfile.IntimacyScore * 0.3
		score += avgProfile.FlexibilityScore * 0.2
	case "casual":
		score += avgProfile.FlexibilityScore * 0.4
		score += avgProfile.SocialityScore * 0.3
		score += avgProfile.ActivityScore * 0.3
	case "chill":
		score += (100 - avgProfile.ActivityScore) * 0.4
		score += avgProfile.FlexibilityScore * 0.3
		score += avgProfile.IntimacyScore * 0.3
	default:
		// 균형잡힌 점수
		score = (avgProfile.SocialityScore + avgProfile.ActivityScore +
			avgProfile.IntimacyScore + avgProfile.ImmersionScore +
			avgProfile.FlexibilityScore) / 5.0
	}

	return score
}
