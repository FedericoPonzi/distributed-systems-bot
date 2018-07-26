package main

import (
	"database/sql"
	"log"
	"time"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository() (mysqlRepository *MysqlRepository) {

	return &MysqlRepository{connectDb()}
}

type FeedRss struct {
	id            int
	twitterHandle string
	url           string
	name          string
	category      int
}

type Shortlink struct {
	id       int
	url      string
	creation time.Time
	clicks   int
}

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", config.getDbConnectionString())
	log.Println(config.getDbConnectionString())
	if err != nil {
		log.Fatal("Impossible open db: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Problem with connection to the database:", err)
	}
	return db
}

func (repo *MysqlRepository) close() {
	repo.db.Close()
}

func (repo *MysqlRepository) getAllFeedRss() (toRet []*FeedRss) {
	rows, err := repo.db.Query("select id, name, twitterHandle, url from feed_rss")
	fatalIfErr(err)
	defer rows.Close()

	for rows.Next() {
		feed := new(FeedRss)
		err := rows.Scan(&feed.id, &feed.name, &feed.twitterHandle, &feed.url)
		if err != nil {
			log.Fatal(err)
		}
		toRet = append(toRet, feed)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return toRet
}

func (repo *MysqlRepository) findFeedRssIdByUrl(url string) (id int) {
	log.Println("Looking for id, from url: ", url)
	err := repo.db.QueryRow("select id from feed_rss where url = ?", url).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
func (repo *MysqlRepository) addFeedRssVisited(feed_id int, updated *time.Time) {
	stmt, err := repo.db.Prepare("INSERT INTO feed_rss_visited(feed_id,updated) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(feed_id, updated)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
}

func (repo *MysqlRepository) getLastFeedRssUpdatedByFeedId(id int) (t time.Time) {
	log.Println("Looking for last update for visited: ", id)
	err := repo.db.QueryRow("select updated from feed_rss_visited where feed_id = ? order by id desc limit 1;", id).Scan(&t)
	if err != nil {
		log.Println("getLastFeedRssUpdatedByFeedId: error ", err)
		t = time.Now() // Because probably this is a new feed rss source.
	}

	return t
}
func (repo MysqlRepository) addShortlink(id string, url string) {
	stmt, err := repo.db.Prepare("INSERT INTO shortlink(uuid, url) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(id, url)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
}
