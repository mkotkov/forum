package application

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/internal/repository"
)

func (a *App) MainPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(
        "public/html/header.html",
        "public/html/footer.html",
        "public/html/forum-card.html",
        "public/html/rec-post.html",
        "public/html/index.html",
        "public/html/create.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Fetch all posts for the main content
    posts, err := queryData(a.repo.GetDB())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch the most recent post for the "Recently created" section
    mostRecentPost, err := queryMostRecentPost(a.repo.GetDB())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch all topics
	topics, err := a.repo.GetAllTopics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Include the most recent post and topics in the data for rendering the template
	data := struct {
		Posts          []repository.Posts
		MostRecentPost *repository.Posts
		Topics         []repository.Topic
		TimeZone       *time.Location
	}{
		Posts:          posts,
		MostRecentPost: mostRecentPost,
		Topics:         topics,  // Убедитесь, что это поле заполнено
		TimeZone:       time.UTC,
	}

    // Render the template
    err = tmpl.ExecuteTemplate(w, "index", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}



func queryData(db *sql.DB) ([]repository.Posts, error) {
    rows, err := db.Query(repository.SQLSelectAllPosts)
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

        // Обработка NULL значения для topic_id
        if post.TopicID != -1 {
			post.Topic = getTopicName(db, post.TopicID)
		} else {
			post.Topic = "No Topic"
		}

        posts = append(posts, post)
    }

    return posts, nil
}

func getTopicName(db *sql.DB, topicID int) string {
	var topicName string
	err := db.QueryRow("SELECT COALESCE(name, 'Unknown Topic') FROM topics WHERE id = ?", topicID).Scan(&topicName)
	if err != nil {
		fmt.Printf("Error retrieving topic name for topic ID %d: %v\n", topicID, err)
		return "Unknown Topic"
	}
	return topicName
}




func queryMostRecentPost(db *sql.DB) (*repository.Posts, error) {
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

