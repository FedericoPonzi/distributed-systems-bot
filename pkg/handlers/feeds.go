package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	app "github.com/FedericoPonzi/distributed-systems-bot/pkg/main"
	"github.com/mmcdole/gofeed"
)

type FeedRssWrapper struct {
	id            int //from FeedRssEntity
	twitterHandle string
	updated       *time.Time
	items         []*gofeed.Item
}

type FeedHandler struct {
	repo           *app.MysqlRepository
	twitterHandler *TwitterHandler
}

type FeedUpdatedWorkUnit struct {
	lastUpdated time.Time       //Last updated from database
	feed        *FeedRssWrapper //Feed just fetched.
}

func NewFeedHandler(repo *app.MysqlRepository, twitterHandler *TwitterHandler) *FeedHandler {
	return &FeedHandler{repo, twitterHandler}
}
func (handler FeedHandler) Main() {
	/** Download all feeds. **/
	feeds := handler.fetchAllRss()

	/** Init short link service **/
	/** Get the feedsItem to publish from the feed (compare updated value of last time, with feeds/items published time) **/
	feedItemsToPublish := handler.getUpdatedItems(feeds)

	/** Save fetched rss results **/
	handler.saveLastFetched(feeds)
	if len(feedItemsToPublish) > 0 {
		fmt.Println("There are ", len(feedItemsToPublish), " items to publish.")
	} else {
		fmt.Println("There are no item to publish.")
	}
	/** If there are items to publish, do it. **/
	for _, item := range feedItemsToPublish {
		log.Println("Done iterating on items to publish. ", item.Title)
		handler.twitterHandler.PublishLinkWithTitle(item.Title, item.Link)
	}
	//handler.scheduleTweets(feedItemsToPublish)
}

/**
Single worker, that will look at feeds, and add to results all the FeedItem that should be published (=> they are new => published after lastUpdated)
*/
func (handler FeedHandler) getUpdatedFeedsItemWorker(jobs <-chan FeedUpdatedWorkUnit, results chan<- gofeed.Item, wg *sync.WaitGroup) {
	for j := range jobs {
		lastUpdated := j.lastUpdated

		feed := j.feed

		if feed.updated == nil {
			log.Println(feed.id, ": Updated nill. Feed:", feed)
		}
		fmt.Println("Last updated:", lastUpdated, " feed updated:", feed.updated)
		if lastUpdated.Before(*feed.updated) {
			log.Println("(", j.feed.id, ") has updates! Last tweeted post was from ", lastUpdated, " now is ", feed.updated)
			for _, feedItem := range feed.items {
				// Arxiv case. This feed is updated daily so the items haven't these fields.
				if feedItem.PublishedParsed == nil && feedItem.UpdatedParsed == nil {
					results <- *feedItem
				}
				if feedItem.PublishedParsed != nil && lastUpdated.Before(*feedItem.PublishedParsed) &&
					!strings.Contains(feedItem.Title, "Sponsored") {
					fmt.Println("This item is in queue for post:", feedItem.Link)
					results <- *feedItem
				}
			}
		} else {
			log.Println("(", j.feed.id, ") wasn't updated: ", lastUpdated, " is after: ", feed.updated)
		}
	}
	wg.Done()
}

/**
It will find, inside at the fetched feeds, what feed items need to be published. Uses `getUpdatedFeedsItemWorker`
*/
func (handler FeedHandler) getUpdatedItems(feedsFetched []*FeedRssWrapper) (feedItemsToPublish []gofeed.Item) {

	log.Println("Generating list of feed items to publish..")
	jobs := make(chan FeedUpdatedWorkUnit, 100)
	results := make(chan gofeed.Item, 100)
	wg := new(sync.WaitGroup)

	// Spawn workers
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go handler.getUpdatedFeedsItemWorker(jobs, results, wg)
	}
	// Send work
	for _, feed := range feedsFetched {
		lastUpdated := handler.repo.GetLastFeedRssUpdatedByFeedId(feed.id)
		jobs <- FeedUpdatedWorkUnit{lastUpdated: lastUpdated, feed: feed}
	}
	log.Println("done sending jobs.")
	close(jobs)

	//Wait for results
	go func() {
		wg.Wait()
		log.Println("I'm done waiting.")
		close(results)
	}()

	//Fetch results
	for res := range results {
		feedItemsToPublish = append(feedItemsToPublish, res)
	}
	log.Println("done computing updated items")

	return feedItemsToPublish
}

/**
  Downloads all the feeds rss.
  Gets the url from database, run a gorotuine for every url - maybe a fixed pool size would be better?
  Also, maybe run in stream fashion ;)
*/
func (handler FeedHandler) fetchAllRss() (feedsFetched []*FeedRssWrapper) {
	feedsRss := handler.repo.GetAllFeedRss()

	c := make(chan *FeedRssWrapper, len(feedsRss))
	wg := new(sync.WaitGroup)
	wg.Add(len(feedsRss))

	for i := range feedsRss {
		feed := feedsRss[i]
		go func() {
			defer wg.Done()
			handler.fetchSingleRss(feed, c)
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	for feed := range c {
		feedsFetched = append(feedsFetched, feed)
	}

	log.Println("All work done.")
	return feedsFetched

}

/**
A worker to fetch a single rss feed.
*/
func (handler FeedHandler) fetchSingleRss(rss *app.FeedRss, c chan *FeedRssWrapper) {

	log.Println("I'm gonna fetch: ", rss.Url())
	url := rss.Url()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println("Error: While fetching feed:" + rss.Url() + ", got error:" + err.Error())
		return
	}

	var updated *time.Time

	if feed.UpdatedParsed != nil {
		updated = feed.UpdatedParsed
		fmt.Println(feed.Title + " updated:" + feed.Published)
	} else {
		// It may happen that there is no "updated" field. In this case, get the last post published date:
		if len(feed.Items) > 0 {
			log.Println("Updated of [" + feed.Title + "] is nil. Last article published: " + feed.Items[0].Published)
			updated = feed.Items[0].PublishedParsed
		} else {
			log.Println("Updated of [" + feed.Title + "] is nil. Apparently, feed.Items is empty. Please investigate.")
			return
		}
	}
	log.Println("Feed fetched for: " + rss.Url())
	toRet := FeedRssWrapper{id: rss.Id(), twitterHandle: rss.TwitterHandle(), updated: updated, items: feed.Items}
	c <- &toRet
}

/**
Store the last time these feeds have been fetched + they're `updated` value.
*/
func (handler FeedHandler) saveLastFetched(feeds []*FeedRssWrapper) {
	for _, f := range feeds {
		handler.repo.AddFeedRssVisited(f.id, f.updated)
	}

}
func (handler FeedHandler) scheduleTweets(items []gofeed.Item) {

	for _, item := range items {
		fmt.Println(item)
		//handler.twitterHandler.PublishLinkWithTitle(item.Title, item.Link)
	}
}
