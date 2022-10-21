package hn_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"hn_feed/config"
	"hn_feed/db"
	"hn_feed/db/models"
	"net/http"
	"sync"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0"

func GetNewPosts(feedType string) (map[int]models.Post, error) {
	var postsIds []int
	r, err := http.Get(fmt.Sprintf("%s/%s.json", baseUrl, feedType))
	if err != nil {
		return nil, errors.New("can't get new posts")
	}

	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&postsIds)

	var wg sync.WaitGroup
	var postsLock sync.RWMutex
	posts := make(map[int]models.Post, config.Config.MaxPosts)
	for i := 0; i < config.Config.MaxPosts; i++ {
		wg.Add(1)
		go func(postCount int) {
			defer wg.Done()
			r, err := http.Get(fmt.Sprintf("%s/item/%d.json", baseUrl, postsIds[postCount]))
			if err != nil {
				return
			}

			var post models.Post
			json.NewDecoder(r.Body).Decode(&post)
			if post.Title == "" {
				return
			}

			if post.Url == "" {
				post.Url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", post.PostId)
			}

			postsLock.Lock()
			db.DB.FirstOrCreate(&post, "post_id = ?", post.PostId) // Filling the ID field.
			posts[postCount] = post
			postsLock.Unlock()
		}(i)
	}
	wg.Wait()
	return posts, nil
}