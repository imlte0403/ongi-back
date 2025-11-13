package main

import (
	"fmt"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
)

func main() {
	// 데이터베이스 연결
	database.Connect()

	// 기존 club_members 삭제
	database.DB.Exec("DELETE FROM club_members")
	database.DB.Exec("UPDATE clubs SET member_count = 0")

	fmt.Println("Deleted all club members")

	// 그룹 매칭 실행
	err := services.MatchUsersToClubs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Matching completed!")

	// 결과 확인: 사용자 76과 같은 클럽에 있는 유사 사용자 확인
	var user76Members []models.ClubMember
	database.DB.Preload("Club").Where("user_id = ?", 76).Find(&user76Members)

	fmt.Printf("\n사용자 76이 가입한 클럽: %d개\n", len(user76Members))

	for _, member := range user76Members {
		fmt.Printf("\n클럽: %s (ID: %d)\n", member.Club.Name, member.ClubID)

		// 이 클럽의 모든 멤버 확인
		var clubMembers []models.ClubMember
		database.DB.Preload("User").Where("club_id = ?", member.ClubID).Find(&clubMembers)

		fmt.Printf("  총 멤버 수: %d명\n", len(clubMembers))
		fmt.Println("  멤버 목록:")
		for _, cm := range clubMembers {
			// 사용자 프로필 가져오기
			var profile models.UserProfile
			database.DB.Where("user_id = ?", cm.UserID).First(&profile)

			fmt.Printf("    - ID: %d, 이름: %s, 프로필: (사교성:%.0f, 활동성:%.0f, 유연성:%.0f)\n",
				cm.UserID, cm.User.Name, profile.SocialityScore, profile.ActivityScore, profile.FlexibilityScore)
		}
	}
}
