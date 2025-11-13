package main

import (
	"encoding/json"
	"fmt"
	"log"
	"ongi-back/database"
	"ongi-back/models"

	"github.com/joho/godotenv"
)

type PreferredScores struct {
	Sociality   float64 `json:"sociality"`
	Activity    float64 `json:"activity"`
	Intimacy    float64 `json:"intimacy"`
	Immersion   float64 `json:"immersion"`
	Flexibility float64 `json:"flexibility"`
}

func seedClubs() {
	clubs := []struct {
		Name             string
		Description      string
		Category         string
		Vibe             string
		MeetingFrequency string
		Location         string
		MemberCount      int
		MaxMembers       int
		Tags             []string
		PreferredScores  PreferredScores
	}{
		// 조용하고 깊은 대화 선호 클럽
		{
			Name:             "조용한 아침 책모임",
			Description:      "매주 토요일 오전, 조용한 카페에서 깊은 대화를 나누는 책모임입니다",
			Category:         "책모임",
			Vibe:             "cozy",
			MeetingFrequency: "주 1회",
			Location:         "강남",
			MemberCount:      5,
			MaxMembers:       8,
			Tags:             []string{"조용한", "깊은대화", "아침형인간", "카페"},
			PreferredScores:  PreferredScores{Sociality: 30, Activity: 40, Intimacy: 80, Immersion: 90, Flexibility: 40},
		},
		{
			Name:             "심야 철학 토론방",
			Description:      "격주 금요일 밤, 인생과 철학에 대해 깊이 있게 이야기하는 모임",
			Category:         "토론",
			Vibe:             "deep",
			MeetingFrequency: "격주",
			Location:         "홍대",
			MemberCount:      4,
			MaxMembers:       6,
			Tags:             []string{"밤형인간", "철학", "깊은사색", "소규모"},
			PreferredScores:  PreferredScores{Sociality: 20, Activity: 30, Intimacy: 85, Immersion: 95, Flexibility: 30},
		},
		{
			Name:             "작은 글쓰기 모임",
			Description:      "조용한 도서관에서 각자 글을 쓰고 피드백을 나누는 모임",
			Category:         "글쓰기",
			Vibe:             "cozy",
			MeetingFrequency: "월 2회",
			Location:         "신촌",
			MemberCount:      6,
			MaxMembers:       8,
			Tags:             []string{"글쓰기", "창작", "피드백", "집중"},
			PreferredScores:  PreferredScores{Sociality: 25, Activity: 35, Intimacy: 75, Immersion: 90, Flexibility: 50},
		},

		// 사교적이고 활발한 클럽
		{
			Name:             "주말 등산 러버스",
			Description:      "매주 일요일 다양한 산을 오르며 건강과 친목을 다지는 모임",
			Category:         "운동",
			Vibe:             "energetic",
			MeetingFrequency: "주 1회",
			Location:         "수도권",
			MemberCount:      12,
			MaxMembers:       20,
			Tags:             []string{"등산", "운동", "친목", "활발한"},
			PreferredScores:  PreferredScores{Sociality: 85, Activity: 90, Intimacy: 60, Immersion: 50, Flexibility: 70},
		},
		{
			Name:             "테니스 동호회 ACE",
			Description:      "주 2회 테니스를 치고 회식도 자주 하는 활기찬 동호회",
			Category:         "운동",
			Vibe:             "energetic",
			MeetingFrequency: "주 2회",
			Location:         "강남",
			MemberCount:      16,
			MaxMembers:       24,
			Tags:             []string{"테니스", "운동", "친목", "회식"},
			PreferredScores:  PreferredScores{Sociality: 90, Activity: 85, Intimacy: 55, Immersion: 45, Flexibility: 75},
		},
		{
			Name:             "번개 맛집 탐방대",
			Description:      "갑자기 번개 모임으로 맛집을 탐방하는 유연한 모임",
			Category:         "맛집",
			Vibe:             "casual",
			MeetingFrequency: "비정기",
			Location:         "서울 전역",
			MemberCount:      20,
			MaxMembers:       30,
			Tags:             []string{"맛집", "번개", "자유로운", "다양한"},
			PreferredScores:  PreferredScores{Sociality: 95, Activity: 80, Intimacy: 40, Immersion: 35, Flexibility: 95},
		},

		// 중간 밸런스 클럽
		{
			Name:             "주말 브런치 클럽",
			Description:      "토요일 아침 브런치를 먹으며 편하게 이야기 나누는 모임",
			Category:         "브런치",
			Vibe:             "casual",
			MeetingFrequency: "주 1회",
			Location:         "연남동",
			MemberCount:      8,
			MaxMembers:       12,
			Tags:             []string{"브런치", "주말", "편한", "카페"},
			PreferredScores:  PreferredScores{Sociality: 65, Activity: 55, Intimacy: 60, Immersion: 50, Flexibility: 60},
		},
		{
			Name:             "영화 감상 모임",
			Description:      "격주로 영화를 보고 가볍게 감상을 나누는 모임",
			Category:         "문화",
			Vibe:             "casual",
			MeetingFrequency: "격주",
			Location:         "강남",
			MemberCount:      10,
			MaxMembers:       15,
			Tags:             []string{"영화", "문화", "감상", "여유"},
			PreferredScores:  PreferredScores{Sociality: 60, Activity: 50, Intimacy: 55, Immersion: 70, Flexibility: 55},
		},
		{
			Name:             "보드게임 나이트",
			Description:      "금요일 밤 보드게임 카페에서 다양한 게임을 즐기는 모임",
			Category:         "게임",
			Vibe:             "casual",
			MeetingFrequency: "주 1회",
			Location:         "홍대",
			MemberCount:      12,
			MaxMembers:       16,
			Tags:             []string{"보드게임", "밤", "실내", "전략"},
			PreferredScores:  PreferredScores{Sociality: 70, Activity: 45, Intimacy: 50, Immersion: 75, Flexibility: 50},
		},

		// 활동성 높은 클럽
		{
			Name:             "새벽 러닝 크루",
			Description:      "매일 새벽 6시 한강에서 러닝하는 열정적인 모임",
			Category:         "운동",
			Vibe:             "energetic",
			MeetingFrequency: "주 5회",
			Location:         "한강",
			MemberCount:      18,
			MaxMembers:       25,
			Tags:             []string{"러닝", "새벽", "루틴", "열정"},
			PreferredScores:  PreferredScores{Sociality: 75, Activity: 95, Intimacy: 45, Immersion: 40, Flexibility: 30},
		},
		{
			Name:             "클라이밍 동호회",
			Description:      "주 2-3회 클라이밍 센터에서 실력을 키우는 모임",
			Category:         "운동",
			Vibe:             "energetic",
			MeetingFrequency: "주 2회",
			Location:         "강남/홍대",
			MemberCount:      14,
			MaxMembers:       20,
			Tags:             []string{"클라이밍", "운동", "도전", "실내"},
			PreferredScores:  PreferredScores{Sociality: 70, Activity: 90, Intimacy: 55, Immersion: 60, Flexibility: 65},
		},
		{
			Name:             "요가 힐링 클럽",
			Description:      "주 2회 요가로 몸과 마음을 힐링하는 모임",
			Category:         "운동",
			Vibe:             "chill",
			MeetingFrequency: "주 2회",
			Location:         "강남",
			MemberCount:      10,
			MaxMembers:       15,
			Tags:             []string{"요가", "힐링", "명상", "건강"},
			PreferredScores:  PreferredScores{Sociality: 55, Activity: 70, Intimacy: 65, Immersion: 80, Flexibility: 60},
		},

		// 친밀도 높은 소규모 클럽
		{
			Name:             "소규모 독서 토론",
			Description:      "4-6명이 한 달에 한 권씩 책을 읽고 깊이 있게 토론",
			Category:         "책모임",
			Vibe:             "deep",
			MeetingFrequency: "월 1회",
			Location:         "강남",
			MemberCount:      5,
			MaxMembers:       6,
			Tags:             []string{"독서", "토론", "소규모", "진지함"},
			PreferredScores:  PreferredScores{Sociality: 30, Activity: 35, Intimacy: 90, Immersion: 95, Flexibility: 40},
		},
		{
			Name:             "친한 친구들의 보드게임",
			Description:      "매주 같은 멤버들이 모여 편하게 게임하는 모임",
			Category:         "게임",
			Vibe:             "cozy",
			MeetingFrequency: "주 1회",
			Location:         "신촌",
			MemberCount:      6,
			MaxMembers:       8,
			Tags:             []string{"보드게임", "친목", "소규모", "편안함"},
			PreferredScores:  PreferredScores{Sociality: 40, Activity: 50, Intimacy: 85, Immersion: 65, Flexibility: 45},
		},
		{
			Name:             "소수 정예 코딩 스터디",
			Description:      "4명이서 매주 알고리즘 문제를 풀고 리뷰하는 모임",
			Category:         "학습",
			Vibe:             "deep",
			MeetingFrequency: "주 1회",
			Location:         "온라인/오프라인",
			MemberCount:      4,
			MaxMembers:       5,
			Tags:             []string{"코딩", "학습", "알고리즘", "집중"},
			PreferredScores:  PreferredScores{Sociality: 25, Activity: 40, Intimacy: 80, Immersion: 95, Flexibility: 35},
		},

		// 유연성 높은 클럽
		{
			Name:             "자유로운 산책 모임",
			Description:      "원하는 사람만 참여하는 느슨한 산책 모임",
			Category:         "산책",
			Vibe:             "chill",
			MeetingFrequency: "비정기",
			Location:         "서울 각지",
			MemberCount:      15,
			MaxMembers:       25,
			Tags:             []string{"산책", "자유", "유연함", "힐링"},
			PreferredScores:  PreferredScores{Sociality: 60, Activity: 55, Intimacy: 45, Immersion: 40, Flexibility: 90},
		},
		{
			Name:             "즉흥 사진 촬영대",
			Description:      "날씨 좋은 날 번개로 모여서 사진 찍으러 다니는 모임",
			Category:         "취미",
			Vibe:             "casual",
			MeetingFrequency: "비정기",
			Location:         "서울 전역",
			MemberCount:      10,
			MaxMembers:       15,
			Tags:             []string{"사진", "번개", "자유", "창작"},
			PreferredScores:  PreferredScores{Sociality: 65, Activity: 70, Intimacy: 50, Immersion: 60, Flexibility: 95},
		},

		// 몰입도 높은 클럽
		{
			Name:             "심층 역사 스터디",
			Description:      "한국사를 깊이 있게 공부하고 토론하는 학습 모임",
			Category:         "학습",
			Vibe:             "deep",
			MeetingFrequency: "주 1회",
			Location:         "강남",
			MemberCount:      7,
			MaxMembers:       10,
			Tags:             []string{"역사", "학습", "토론", "깊이"},
			PreferredScores:  PreferredScores{Sociality: 40, Activity: 35, Intimacy: 70, Immersion: 95, Flexibility: 40},
		},
		{
			Name:             "악기 연주 마스터반",
			Description:      "기타/피아노를 진지하게 연습하고 합주하는 모임",
			Category:         "음악",
			Vibe:             "deep",
			MeetingFrequency: "주 2회",
			Location:         "홍대",
			MemberCount:      8,
			MaxMembers:       12,
			Tags:             []string{"악기", "연주", "연습", "집중"},
			PreferredScores:  PreferredScores{Sociality: 50, Activity: 60, Intimacy: 65, Immersion: 90, Flexibility: 45},
		},
		{
			Name:             "프로그래밍 부트캠프",
			Description:      "주말 8시간 동안 집중해서 프로젝트를 만드는 모임",
			Category:         "학습",
			Vibe:             "deep",
			MeetingFrequency: "주 1회",
			Location:         "강남",
			MemberCount:      6,
			MaxMembers:       10,
			Tags:             []string{"프로그래밍", "프로젝트", "몰입", "학습"},
			PreferredScores:  PreferredScores{Sociality: 45, Activity: 50, Intimacy: 60, Immersion: 95, Flexibility: 35},
		},

		// 추가 다양한 클럽
		{
			Name:             "주말 카페 투어",
			Description:      "매주 새로운 카페를 방문하며 커피를 즐기는 모임",
			Category:         "카페",
			Vibe:             "casual",
			MeetingFrequency: "주 1회",
			Location:         "서울 각지",
			MemberCount:      10,
			MaxMembers:       15,
			Tags:             []string{"카페", "커피", "탐방", "여유"},
			PreferredScores:  PreferredScores{Sociality: 70, Activity: 60, Intimacy: 55, Immersion: 50, Flexibility: 70},
		},
		{
			Name:             "펫프렌즈 산책모임",
			Description:      "반려견과 함께 주말 한강 산책하는 모임",
			Category:         "반려동물",
			Vibe:             "chill",
			MeetingFrequency: "주 1회",
			Location:         "한강",
			MemberCount:      12,
			MaxMembers:       20,
			Tags:             []string{"반려견", "산책", "한강", "힐링"},
			PreferredScores:  PreferredScores{Sociality: 75, Activity: 65, Intimacy: 60, Immersion: 45, Flexibility: 65},
		},
		{
			Name:             "저녁 명상 모임",
			Description:      "평일 저녁 명상과 요가로 하루를 마무리하는 모임",
			Category:         "명상",
			Vibe:             "chill",
			MeetingFrequency: "주 3회",
			Location:         "강남",
			MemberCount:      8,
			MaxMembers:       12,
			Tags:             []string{"명상", "요가", "힐링", "평화"},
			PreferredScores:  PreferredScores{Sociality: 40, Activity: 55, Intimacy: 70, Immersion: 85, Flexibility: 50},
		},
		{
			Name:             "주식투자 스터디",
			Description:      "주말마다 시장을 분석하고 투자 전략을 공유하는 모임",
			Category:         "재테크",
			Vibe:             "deep",
			MeetingFrequency: "주 1회",
			Location:         "강남",
			MemberCount:      9,
			MaxMembers:       12,
			Tags:             []string{"주식", "투자", "재테크", "분석"},
			PreferredScores:  PreferredScores{Sociality: 55, Activity: 45, Intimacy: 65, Immersion: 85, Flexibility: 50},
		},
		{
			Name:             "오픈 마이크 공연 모임",
			Description:      "한 달에 한 번 작은 공연장에서 자유롭게 공연하는 모임",
			Category:         "공연",
			Vibe:             "energetic",
			MeetingFrequency: "월 1회",
			Location:         "홍대",
			MemberCount:      15,
			MaxMembers:       25,
			Tags:             []string{"공연", "음악", "자유", "열정"},
			PreferredScores:  PreferredScores{Sociality: 80, Activity: 75, Intimacy: 50, Immersion: 70, Flexibility: 80},
		},
	}

	for _, clubData := range clubs {
		// Tags를 JSON으로 변환
		tagsJSON, _ := json.Marshal(clubData.Tags)

		// PreferredScores를 JSON으로 변환
		scoresJSON, _ := json.Marshal(clubData.PreferredScores)

		club := models.Club{
			Name:             clubData.Name,
			Description:      clubData.Description,
			Category:         clubData.Category,
			Vibe:             clubData.Vibe,
			MeetingFrequency: clubData.MeetingFrequency,
			Location:         clubData.Location,
			MemberCount:      clubData.MemberCount,
			MaxMembers:       clubData.MaxMembers,
			Tags:             string(tagsJSON),
			PreferredScores:  string(scoresJSON),
		}

		err := database.DB.Create(&club).Error
		if err != nil {
			log.Printf("Failed to create club %s: %v", clubData.Name, err)
		} else {
			fmt.Printf("Created club: %s\n", clubData.Name)
		}
	}
}

func main() {
	// .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = database.Migrate()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Starting club seeding...")
	seedClubs()
	fmt.Println("Club seeding completed!")
}
