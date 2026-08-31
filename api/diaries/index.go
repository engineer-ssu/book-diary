package handler

import (
  "encoding/json"
  "net/http"
  "time"
)

// --- Struct 정의 ---

type User struct {
  UserID        int64  `json:"userId"`
  Nickname      string `json:"nickname"`
  Email         string `json:"email,omitempty"`
  ProfileImgUrl string `json:"profileImgUrl,omitempty"`
}

type Book struct {
  ISBN       string `json:"isbn"`
  Title      string `json:"title"`
  Author     string `json:"author"`
  Publisher  string `json:"publisher,omitempty"`
  CoverImage string `json:"coverImage,omitempty"`
  TotalPages int    `json:"totalPages"`
}

type CreateDiaryReq struct {
  UserID    int64    `json:"userId"`
  Book      Book     `json:"book"`
  Status    string   `json:"status"`
  Rating    float64  `json:"rating"`
  ReadPages int      `json:"readPages"`
  StartDate *string  `json:"startDate"`
  EndDate   *string  `json:"endDate"`
  Content   string   `json:"content"`
  IsPrivate bool     `json:"isPrivate"`
  Tags      []string `json:"tags"`
}

type DiaryRes struct {
  DiaryID   int64     `json:"diaryId"`
  User      User      `json:"user"`
  Status    string    `json:"status"`
  Rating    float64   `json:"rating"`
  ReadPages int       `json:"readPages"`
  StartDate *string   `json:"startDate,omitempty"`
  EndDate   *string   `json:"endDate,omitempty"`
  Content   string    `json:"content"`
  IsPrivate bool      `json:"isPrivate"`
  CreatedAt time.Time `json:"createdAt"`
  Book      Book      `json:"book"`
  Tags      []string  `json:"tags"`
  LikeCount int       `json:"likeCount"`
}

// --- 메모리 인메모리 Mock 데이터 ---

var mockStartDate = "2026-08-01"
var mockEndDate = "2026-08-15"

var mockUsers = map[int64]User{
  1: {
    UserID:        1,
    Nickname:      "북워머",
    Email:         "bookworm@example.com",
    ProfileImgUrl: "https://api.dicebear.com/7.x/bottts/svg?seed=bookworm",
  },
  2: {
    UserID:        2,
    Nickname:      "독서왕한강팬",
    Email:         "reader@example.com",
    ProfileImgUrl: "https://api.dicebear.com/7.x/bottts/svg?seed=reader",
  },
}

// Mock 일기 데이터를 생성하여 반환하는 헬퍼 함수
func getMockDiaries() []DiaryRes {
  now := time.Now()
  return []DiaryRes{
    {
      DiaryID:   1,
      User:      mockUsers[1],
      Status:    "READING",
      Rating:    4.5,
      ReadPages: 320,
      StartDate: &mockStartDate,
      Content:   "Go 언어로 서버리스 함수를 작성하면서 가독성 높은 코드 설계의 필요성을 절감하고 읽기 시작했다. 3장의 함수 작성 원칙 부분이 특히 유용함.",
      IsPrivate: false,
      CreatedAt: now.AddDate(0, 0, -2),
      Book: Book{
        ISBN:       "9791162540169",
        Title:      "클린 코드 (Clean Code)",
        Author:     "로버트 C. 마틴",
        Publisher:  "인사이트",
        CoverImage: "https://images.unsplash.com/photo-1532012197267-da84d127e765?w=300",
        TotalPages: 584,
      },
      Tags:      []string{"개발", "리팩토링"},
      LikeCount: 12,
    },
    {
      DiaryID:   2,
      User:      mockUsers[2],
      Status:    "COMPLETED",
      Rating:    5.0,
      ReadPages: 216,
      StartDate: &mockStartDate,
      EndDate:   &mockEndDate,
      Content:   "문체가 너무 아름답고 묵직하다. 가슴이 먹먹해지는 경험이었다.",
      IsPrivate: false,
      CreatedAt: now.AddDate(0, 0, -5),
      Book: Book{
        ISBN:       "9788936434267",
        Title:      "소년이 온다",
        Author:     "한강",
        Publisher:  "창비",
        CoverImage: "https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=300",
        TotalPages: 216,
      },
      Tags:      []string{"소설", "문학"},
      LikeCount: 45,
    },
  }
}

// --- Vercel Entrypoint ---

func Handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

  if r.Method == http.MethodOptions {
    w.WriteHeader(http.StatusOK)
    return
  }

  switch r.Method {
  case http.MethodGet:
    getDiariesSocialMock(w, r)
  case http.MethodPost:
    createDiarySocialMock(w, r)
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
    json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
  }
}

// [GET] /api/diaries - 소셜 피드 일기 목록 반환
func getDiariesSocialMock(w http.ResponseWriter, r *http.Request) {
  userIdQuery := r.URL.Query().Get("userId")
  statusFilter := r.URL.Query().Get("status")

  filtered := make([]DiaryRes, 0)
  mockDiaries := getMockDiaries()

  for _, diary := range mockDiaries {
    if diary.IsPrivate {
      continue
    }

    if userIdQuery != "" && diary.User.UserID != 1 {
      continue
    }

    if statusFilter != "" && diary.Status != statusFilter {
      continue
    }

    filtered = append(filtered, diary)
  }

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(filtered)
}

// [POST] /api/diaries - 소셜 피드용 일기 작성
func createDiarySocialMock(w http.ResponseWriter, r *http.Request) {
  var req CreateDiaryReq
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
    return
  }

  author, exists := mockUsers[req.UserID]
  if !exists {
    author = mockUsers[1]
  }

  newID := time.Now().Unix()

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "message": "[Mock] 소셜 독서 일기가 등록되었습니다.",
    "data": map[string]interface{}{
      "diaryId":   newID,
      "author":    author.Nickname,
      "bookTitle": req.Book.Title,
      "isPrivate": req.IsPrivate,
      "createdAt": time.Now().Format(time.RFC3339),
    },
  })
}