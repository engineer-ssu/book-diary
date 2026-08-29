package handler

import "time"

// 책 정보 Struct
type Book struct {
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Publisher  string `json:"publisher,omitempty"`
	CoverImage string `json:"coverImage,omitempty"`
	TotalPages int    `json:"totalPages"`
}

// 일기 등록 Request Body
type CreateDiaryReq struct {
	Book       Book     `json:"book"`
	Status     string   `json:"status"`
	Rating     float64  `json:"rating"`
	ReadPages  int      `json:"readPages"`
	StartDate  *string  `json:"startDate"` // YYYY-MM-DD
	EndDate    *string  `json:"endDate"`   // YYYY-MM-DD
	Content    string   `json:"content"`
	IsPrivate  bool     `json:"isPrivate"`
	Tags       []string `json:"tags"`
}

// 일기 Response Struct
type DiaryRes struct {
	DiaryID   int64     `json:"diaryId"`
	Status    string    `json:"status"`
	Rating    float64   `json:"rating"`
	ReadPages int       `json:"readPages"`
	StartDate *string   `json:"startDate"`
	EndDate   *string   `json:"endDate"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"isPrivate"`
	CreatedAt time.Time `json:"createdAt"`
	Book      Book      `json:"book"`
	Tags      []string  `json:"tags"`
}