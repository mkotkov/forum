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

    filter := r.FormValue("filters")
    topicFilter := r.FormValue("topic")
    selectedTopic, err := strconv.Atoi(topicFilter)
    if err != nil {
        selectedTopic = -1
    }

    var posts []repository.Posts

    // Переписываем часть кода для корректного отображения всех постов при отсутствии выбранного топика
    if selectedTopic == 0 {
        posts, err = queryData(a.repo.GetDB(), filter, -1)
    } else {
        posts, err = queryData(a.repo.GetDB(), filter, selectedTopic)
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    mostRecentPost, err := queryMostRecentPost(a.repo.GetDB())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    topics, err := a.repo.GetAllTopics(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Печатаем в консоль информацию о загруженных постах для отладки
    fmt.Printf("Loaded Posts: %+v\n", posts)

    data := struct {
        Posts          []repository.Posts
        MostRecentPost *repository.Posts
        Topics         []repository.Topic
        TimeZone       *time.Location
        Filter         string
        SelectedTopic  int
    }{
        Posts:          posts,
        MostRecentPost: mostRecentPost,
        Topics:         topics,
        TimeZone:       time.UTC,
        Filter:         filter,
        SelectedTopic:  selectedTopic,
    }

    if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}





func queryData(db *sql.DB, filter string, topicID int) ([]repository.Posts, error) {
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

