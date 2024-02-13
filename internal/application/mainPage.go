package application

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"forum/internal/repository"
)

func (a *App) MainPage(w http.ResponseWriter, r *http.Request, message string, isAuthorized bool) {
	userID, err := repository.GetUserIDFromToken(r, a.repo)
	if err != nil {
		a.UnregPage(w, r, "<div class="+"error"+"><p>You have an incorrect userID!</p></div>", false)
		return
	}

	userLogin, err := a.repo.GetUserLoginByID(r.Context(), userID)
	if err != nil {
		a.UnregPage(w, r, "<div class="+"error"+"><p>You have an incorrect username!</p></div>", false)
		return
	}

	tmpl, err := template.ParseFiles(
		"public/html/unreg-header.html",
		"public/html/header.html",
		"public/html/footer.html",
		"public/html/forum-card.html",
		"public/html/login-form.html",
		"public/html/signup.html",
		"public/html/rec-post.html",
		"public/html/index.html",
		"public/html/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := r.FormValue("filters")
	topicFilter := r.FormValue("topic")
	selectedTopic, err := strconv.Atoi(topicFilter)
	if err != nil {
		selectedTopic = -1
	}

	var posts []repository.Posts

	if selectedTopic == 0 {
		posts, err = QueryData(a.repo.GetDB(), filter, -1)
	} else {
		posts, err = QueryData(a.repo.GetDB(), filter, selectedTopic)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mostRecentPost, err := QueryMostRecentPost(a.repo.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topics, err := a.repo.GetAllTopics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Loaded Posts: %+v\n", posts)

	data := struct {
		Posts          []repository.Posts
		MostRecentPost *repository.Posts
		Topics         []repository.Topic
		TimeZone       *time.Location
		Filter         string
		SelectedTopic  int
		Message        string
		IsAuthorized   bool
		NameUser       string
	}{
		Posts:          posts,
		MostRecentPost: mostRecentPost,
		Topics:         topics,
		TimeZone:       time.UTC,
		Filter:         filter,
		SelectedTopic:  selectedTopic,
		Message:        message,
		IsAuthorized:   isAuthorized,
		NameUser:       userLogin,
	}

	if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func QueryData(db *sql.DB, filter string, topicID int) ([]repository.Posts, error) {
	var query string
	var args []interface{}
	if topicID == -1 {
		// If "All" is selected, fetch all posts
		switch filter {
		case "oldest":
			query = repository.SQLSelectAllPosts + " ORDER BY post_date ASC"
		case "most_likes":
			query = repository.SQLSelectAllPosts + " ORDER BY like_count DESC"
		case "most_dislikes":
			query = repository.SQLSelectAllPosts + " ORDER BY dislike_count DESC"
		case "most_recent", "":
			query = repository.SQLSelectAllPosts + " ORDER BY post_date DESC"
		default:
			return nil, fmt.Errorf("invalid filter value")
		}
	} else {
		// Otherwise, fetch posts for the selected topic
		switch filter {
		case "oldest":
			query = repository.SQLSelectAllPosts + " WHERE COALESCE(topic_id, -1) = ? ORDER BY post_date ASC"
		case "most_likes":
			query = repository.SQLSelectAllPosts + " WHERE COALESCE(topic_id, -1) = ? ORDER BY like_count DESC"
		case "most_dislikes":
			query = repository.SQLSelectAllPosts + " WHERE COALESCE(topic_id, -1) = ? ORDER BY dislike_count DESC"
		case "most_recent", "":
			query = repository.SQLSelectAllPosts + " WHERE COALESCE(topic_id, -1) = ? ORDER BY post_date DESC"
		default:
			return nil, fmt.Errorf("invalid filter value")
		}
		args = append(args, topicID)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []repository.Posts

	for rows.Next() {
		var post repository.Posts
		err := rows.Scan(
			&post.Id,
			&post.Author,
			&post.PostDate,
			&post.Title,
			&post.FullText,
			&post.Slug,
			&post.LikeCount,
			&post.DislikeCount,
			&post.TopicID,
		)
		if err != nil {
			return nil, err
		}

		if post.TopicID != -1 {
			post.Topic = GetTopicName(db, post.TopicID)
		} else {
			post.Topic = "No Topic"
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetTopicName(db *sql.DB, topicID int) string {
	var topicName string
	err := db.QueryRow("SELECT COALESCE(name, 'Unknown Topic') FROM topics WHERE id = ?", topicID).Scan(&topicName)
	if err != nil {
		fmt.Printf("Error retrieving topic name for topic ID %d: %v\n", topicID, err)
		return "Unknown Topic"
	}
	return topicName
}

func QueryMostRecentPost(db *sql.DB) (*repository.Posts, error) {
	var post repository.Posts

	// Check if there are any posts in the database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return nil, err
	}

	// If there are no posts, return an empty result
	if count == 0 {
		return nil, nil
	}

	// Otherwise, proceed with the original query
	row := db.QueryRow(repository.SQLSelectMostRecentPost)
	err = row.Scan(&post.Id, &post.Author, &post.PostDate, &post.Title, &post.FullText, &post.Slug, &post.LikeCount, &post.DislikeCount, &post.Topic)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
