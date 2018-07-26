package main

import (
	"github.com/mmcdole/gofeed"
	"sync"
	"fmt"
	"log"
	"time"
)

type FeedRssWrapper struct {
	id int //from FeedRssEntity
	twitterHandle string
	updated *time.Time
	items []* gofeed.Item
}

type FeedHandler struct {
	repo *MysqlRepository
	twitterHandler *TwitterHandler
}

func NewFeedHandler (repo *MysqlRepository, twitterHandler *TwitterHandler) *FeedHandler {
	return &FeedHandler{repo, twitterHandler}
}
func (handler FeedHandler) main() {
	feeds := handler.fetchAllRss()
	shortlinkService := NewShortLinkService(handler.repo)

	handler.saveLastFetched(feeds)
	feedItemsToPublish := handler.getUpdatedItems(feeds)
	if len(feedItemsToPublish) > 0 {
		log.Println("Done iterating on items to publish.", feedItemsToPublish[0].Title)
		link := "https://ds.fponzi.me/" + shortlinkService.generateShortlink(feedItemsToPublish[0])
		handler.twitterHandler.publishLinkWithTitle("This is a test. Sorry for the inconvenience." + feedItemsToPublish[0].Title, link)
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
		if lastUpdated.Before(*feed.updated) {
			for _, feedItem := range feed.items {

				if (feedItem.UpdatedParsed != nil && lastUpdated.Before(*feed.updated)) ||
					(feedItem.PublishedParsed != nil && lastUpdated.Before(*feedItem.PublishedParsed)){
					fmt.Println("This item is in queue for post:", feedItem.Link)
					results <- *feedItem
				}

			}
		}else {
			log.Println("It was updated after: ", lastUpdated, " is after: ", feed.updated)
		}
	}
	wg.Done()
}
type FeedUpdatedWorkUnit struct{
	lastUpdated time.Time
	feed *FeedRssWrapper
}

/**
	It will find, inside at the fetched feeds, what feed items need to be published. It will use `getUpdatedFeedsItemWorker`
 */
func (handler FeedHandler) getUpdatedItems(feedsFetched[] *FeedRssWrapper) (feedItemsToPublish []gofeed.Item ){

	log.Println("Generating list of feed items to publish..")
	jobs := make(chan FeedUpdatedWorkUnit, 100)
	results := make(chan gofeed.Item, 100)
	wg := new(sync.WaitGroup)

	for w := 1; w <= 3; w++ {
		go handler.getUpdatedFeedsItemWorker(jobs, results, wg)
	}

	for _, feed := range feedsFetched {
		wg.Add(1)
		lastUpdated := handler.repo.getLastFeedRssUpdatedByFeedId(feed.id)
		jobs <- FeedUpdatedWorkUnit{lastUpdated: lastUpdated, feed:feed}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		feedItemsToPublish = append(feedItemsToPublish, res)
	}

	return feedItemsToPublish
}

/**
   Downloads all the feeds rss.
   Gets the url from database, run a gorotuine for every url - maybe a fixed pool size would be better?
 */
func (handler FeedHandler) fetchAllRss() (feedsFetched [] *FeedRssWrapper) {
	feedsRss := handler.repo.getAllFeedRss()

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

	wg.Wait()
	for i := 0; i < len(feedsRss); i++{
		feedsFetched = append(feedsFetched, <-c)
	}
	close(c)

	log.Println("All work done.")
	return feedsFetched

}
/**
	A worker to fetch a single rss feed.
 */
func (handler FeedHandler) fetchSingleRss(rss *FeedRss, c chan *FeedRssWrapper) {
	log.Println("I'm gonna fetch: ", rss.url)
	url := rss.url //"https://blog.acolyer.org/feed/"
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url) // TODO handle err. Add retry
	log.Println("Has parse errors? ", err)
	//fatalIfErr(err)
	if feed == nil{ // TODO. check err status.
		return
	}
	fmt.Println("Title: ", feed.Title)
	//fmt.Println(feed.UpdatedParsed)
	updated := feed.PublishedParsed
	if feed.UpdatedParsed != nil {
		updated = feed.UpdatedParsed
	}
	toRet := FeedRssWrapper{id: rss.id, twitterHandle:rss.twitterHandle, updated:updated, items:feed.Items}
	c<- &toRet
}
/**
	Store the last time these feeds have been fetched + they're `updated` value.
 */
func (handler FeedHandler) saveLastFetched(feeds []*FeedRssWrapper) {
	for _, f := range feeds {
		handler.repo.addFeedRssVisited(f.id, f.updated)
	}

}
func (handler FeedHandler) scheduleTweets(items []gofeed.Item) {
	for _, item := range items {
		handler.twitterHandler.publishLinkWithTitle(item.Title, item.Link)
	}
}