package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/lius-new/rssagg/internal/database"
)

// startScraping: 抓取内容
// params:{db}:
// params:{concurrency}: 多少协程抓取内容
// params:{timeBetweenRequest}: 请求的间隔时间
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration ", concurrency, timeBetweenRequest)

	// ticker 提供类似定时器的功能
	// for ; ; <-ticker.C, ticker.C的类型是chan
	// 通过遍历chan, 达到定时器的效果.for ;; 可以在程序执行的时候启动, 如果是for range 则需要先等timeBetweenRequest
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		// context.Background() 类似http.Request.Context
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		// sync.WaitGroup 来实现并发
		// feeds.len == concurrency , 所以并发数和concurrency有关
		// 每次循环feeds 都往wg添加1个任务, 然后在循环外部等待任务执行结束
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetchd(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feeds as fetched: ", err)
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
	}

	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" { // 如果item.Description不为空则设置为item.Description, 和Valid 为true, 其默认为description.String=""和description.valid=false
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") { // 重复贴跳过
				continue
			}
			log.Println("failed to create post: ", err)
		}

		log.Println("Found post", item.Title)
	}
	log.Printf("Feed %s collected , %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
