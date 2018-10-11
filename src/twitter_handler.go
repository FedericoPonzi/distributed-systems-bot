package main

import (
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TwitterHandler struct {
	repo *MysqlRepository
	bot *twitter.Client
	dryRun bool
}

type Tweet struct {
	id int
	tweet string
	posted time.Time
	published int8
}

func NewTwitterHandler(repo * MysqlRepository, config TwitterConfig) *TwitterHandler {

	oauthConf := oauth1.NewConfig(config.Consumerkey, config.ConumerSecret)
	token := oauth1.NewToken(config.Token, config.TokenSecret)
	// http.Client will automatically authorize Requests
	httpClient := oauthConf.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	log.Println("Twitter bot initialized!")
	return &TwitterHandler{bot: client, repo: repo,
		dryRun:config.DryRun }
}

func (twitterHandler *TwitterHandler) retweetAndLike(t *twitter.Tweet) {
	tweet, resp, err := twitterHandler.bot.Statuses.Retweet(t.ID, nil)
	if err != nil {
		errStr := fmt.Sprint("Error twitting message:", tweet, resp)
		log.Println(errStr)
	}else {
		log.Println("Succesfully retweeted!")
	}
	favourite := twitter.FavoriteCreateParams{t.ID}
	twitterHandler.bot.Favorites.Create(&favourite)
}

func (twitterHandler *TwitterHandler) getFriendList() (toRet []twitter.User) {
	/*var cursor int64 = -1
	var friendParams twitter.FriendListParams
	for cursor != 0 {
		friendParams = twitter.FriendListParams{UserID: twitterHandler.DistribSystemsId, Cursor: cursor}
		list, _, err := twitterHandler.bot.Friends.List(&friendParams)
		log.Println("Sto prendendo la lista.. ", len(list.Users))
		log.Println(err)
		//fatalIfErr(err)
		for _, user := range list.Users {
			toRet = append(toRet, user)
		}
		cursor = list.NextCursor
	}
	*/
	return toRet
}

func (twitterHandler * TwitterHandler) runStreaming(){
	//https://developer.twitter.com/en/docs/tweets/filter-realtime/guides/basic-stream-parameters
	//https://developer.twitter.com/en/docs/tweets/filter-realtime/api-reference/post-statuses-filter
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.RetweetedStatus == nil {
			//twitterHandler.retweetAndLike(tweet)

			fmt.Println("Tweet text:", tweet.Text)

		}
	}

	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println("Direct", dm.SenderID)
	}

	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")
	//user, _, err := twitterHandler.bot.Users.Lookup(&twitter.UserLookupParams{ScreenName: []string{"distribsystems"}})
	//fmt.Println(user[0].ID)

	//fmt.Printf("Number of friends: %d \n ", len(twitterHandler.getFriendList()))
	// FILTER
	filterParams := &twitter.StreamFilterParams{
		StallWarnings: twitter.Bool(false),
		Track:         []string{"scalability", "distributed systems", "microservices"},
	}

	stream, err := twitterHandler.bot.Streams.Filter(filterParams)
	//twitterHandler.bot.Timelines.UserTimeline(twitter.UserTimelineParams{})
	//twitterHandler.bot.Streams.User(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}

// Publish a shortlink.
func (handler *TwitterHandler) PublishShortLink(shortlink string) (err error) {
	shortlinkService := NewShortLinkService(handler.repo)
	title, err := shortlinkService.FetchTitle(shortlink)
	if err != nil {
		log.Println("Error fetching the title: " + err.Error())
		return err
	}
	return handler.PublishLinkWithTitle(title, shortlink)
}

// Publish a Link - that will be shortlinked TODO
func (handler *TwitterHandler) PublishLink(link string) (err error) {
	shortlinkService := NewShortLinkService(handler.repo)
	shortlink, title, err := shortlinkService.GenerateShortlink(link)
	if err != nil {
		log.Println("Error generating the shortlink!! " + err.Error())
		return err
	}
	return handler.PublishShortLinkWithTitle(title, shortlink)
}

// Publish Title and a Link  that will be shortlinked
func (handler *TwitterHandler) PublishLinkWithTitle(title string, link string) (err error) {
	shortlinkService := NewShortLinkService(handler.repo)
	shortlink := shortlinkService.GenerateShortlinkWithTitle(link, title)
	return handler.PublishShortLinkWithTitle(title, shortlink)
}

// Publish a Title with a shortlink
func (handler *TwitterHandler) PublishShortLinkWithTitle(title string, link string) (err error) {
	/**
			Links uses 23 chars
			So we have max 257 chars for tweet text. 1 chars for new line, so 256.
			if the length of title is greater then 256 chars, then we need to append "..." so 256-3.
		 */
	if len(title) > 256 {
		title = title[:256-4] + "..." //4 and not 3, because arrays start from 0.
	}
	tweet := title + "\n" + link
	return handler.PublishTweet(tweet)
}

// Publish a text tweet.
func (handler *TwitterHandler) PublishTweet(tweet string) (err error){
	if !handler.dryRun {
		res, resp, err := handler.bot.Statuses.Update(tweet, nil)
		fmt.Println(res, resp, err)
		return err
		//TODO: save tweet in db!

	}
	return nil
}