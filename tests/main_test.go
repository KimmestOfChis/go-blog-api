package main

import (
    "blogapi/database"
    "bytes"
    "database/sql"
    "encoding/json"
    "github.com/DATA-DOG/go-sqlmock"
    "net/http"
    "net/http/httptest"
    "testing"
)

var (
    mockDB   *sql.DB
    mock     sqlmock.Sqlmock
    err      error
)

func setup() {
    mockDB, mock, err = sqlmock.New()
    if err != nil {
        panic(err)
    }
    // Mock the Connect function to return our mock database connection
    database.Connect = func() (*sql.DB, error) { return mockDB, nil }
}


func TestGetPost(t *testing.T) {
	setup()

	postID := "some-id"
	mock.ExpectQuery("^SELECT id, title, content FROM posts WHERE id = \\$1$").
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content"}).
			AddRow(postID, "Test Title", "Test Content"))

	req, err := http.NewRequest("GET", "/posts/"+postID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPost)

	handler.ServeHTTP(rr, req)

	expected := `{"ID":"some-id","Title":"Test Title","Content":"Test Content"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAllPosts(t *testing.T) {
	setup()

	mock.ExpectQuery("^SELECT id, title, content FROM posts$").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content"}).
			AddRow("id1", "Title1", "Content1").
			AddRow("id2", "Title2", "Content2"))

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllPosts)

	handler.ServeHTTP(rr, req)

	expected := `[{"ID":"id1","Title":"Title1","Content":"Content1"},{"ID":"id2","Title":"Title2","Content":"Content2"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreatePost(t *testing.T) {
	setup()

	post := &models.Post{
		Title:   "New Post",
		Content: "This is the content",
	}
	postBytes, _ := json.Marshal(post)
	mock.ExpectExec("^INSERT INTO posts").
		WithArgs(post.Title, post.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(postBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createPost)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}


