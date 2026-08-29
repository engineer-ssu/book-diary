package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/sijms/go-ora/v2" // Pure Go Oracle Driver
)

// Oracle DB 커넥션 풀을 열어주는 헬퍼 함수
func getDBConnection() (*sql.DB, error) {
	// Vercel Environment Variable에서 접속 정보 로드
	// 예: oracle://user:password@hostname:1521/service_name
	connStr := os.Getenv("ORACLE_DATABASE_URL")
	return sql.Open("oracle", connStr)
}

// Vercel Entrypoint (반드시 Handler여야 함)
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getDiariesHandler(w, r)
	case http.MethodPost:
		createDiaryHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

// [GET] /api/diaries - 일기 목록 조회
func getDiariesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, `{"error":"DB connection failed"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 독서 상태 필터링 (?status=READING)
	statusFilter := r.URL.Query().Get("status")

	query := `
		SELECT 
			d.diary_id, d.status, d.rating, d.read_pages, d.content, d.is_private, d.created_at,
			b.isbn, b.title, b.author, b.cover_image, b.total_pages
		FROM diaries d
		JOIN books b ON d.isbn = b.isbn
	`
	var rows *sql.Rows
	if statusFilter != "" {
		query += " WHERE d.status = :1 ORDER BY d.created_at DESC"
		rows, err = db.Query(query, statusFilter)
	} else {
		query += " ORDER BY d.created_at DESC"
		rows, err = db.Query(query)
	}

	if err != nil {
		http.Error(w, `{"error":"Query execution failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	diaries := make([]DiaryRes, 0)
	for rows.Next() {
		var d DiaryRes
		var isPrivateNum int
		err := rows.Scan(
			&d.DiaryID, &d.Status, &d.Rating, &d.ReadPages, &d.Content, &isPrivateNum, &d.CreatedAt,
			&d.Book.ISBN, &d.Book.Title, &d.Book.Author, &d.Book.CoverImage, &d.Book.TotalPages,
		)
		if err != nil {
			continue
		}
		d.IsPrivate = (isPrivateNum == 1)
		d.Tags = []string{} // 필요시 별도 JOIN 처리
		diaries = append(diaries, d)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(diaries)
}

// [POST] /api/diaries - 일기 및 도서 등록
func createDiaryHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDiaryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	db, err := getDBConnection()
	if err != nil {
		http.Error(w, `{"error":"DB connection failed"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, `{"error":"Transaction start failed"}`, http.StatusInternalServerError)
		return
	}

	// 1. 책 정보가 업로드되지 않았다면 MERGE(UPSERT) 문으로 저장
	upsertBookSQL := `
		MERGE INTO books b
		USING (SELECT :1 AS isbn FROM dual) src
		ON (b.isbn = src.isbn)
		WHEN NOT MATCHED THEN
			INSERT (isbn, title, author, publisher, cover_image, total_pages)
			VALUES (:1, :2, :3, :4, :5, :6)
	`
	_, err = tx.Exec(upsertBookSQL, req.Book.ISBN, req.Book.Title, req.Book.Author, req.Book.Publisher, req.Book.CoverImage, req.Book.TotalPages)
	if err != nil {
		tx.Rollback()
		http.Error(w, `{"error":"Failed to upsert book"}`, http.StatusInternalServerError)
		return
	}

	// 2. 일기(Diary) 데이터 저장 (임시로 user_id = 1 고정)
	isPrivateVal := 0
	if req.IsPrivate {
		isPrivateVal = 1
	}

	insertDiarySQL := `
		INSERT INTO diaries (user_id, isbn, status, rating, read_pages, content, is_private)
		VALUES (1, :1, :2, :3, :4, :5, :6)
		RETURNING diary_id INTO :7
	`
	var newDiaryID int64
	_, err = tx.Exec(insertDiarySQL, req.Book.ISBN, req.Status, req.Rating, req.ReadPages, req.Content, isPrivateVal, sql.Out{Dest: &newDiaryID})
	if err != nil {
		tx.Rollback()
		http.Error(w, `{"error":"Failed to create diary"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"Transaction commit failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "일기가 성공적으로 저장되었습니다.",
		"data": map[string]interface{}{
			"diaryId": newDiaryID,
		},
	})
}