package pkg

import (
	"database/sql"
	"log"
	"time"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(config *Config) (mysqlRepository *MysqlRepository) {
	return &MysqlRepository{connectDb(config.GetDbConnectionString())}
}

type FeedRss struct {
	id            int
	twitterHandle string
	url           string
	name          string
	category      int
}

func (feed FeedRss) Url() string {
	return feed.url
}

func (feed FeedRss) Id() int {
	return feed.id
}

func (feed FeedRss) TwitterHandle() string {
	return feed.twitterHandle
}

type Shortlink struct {
	id       int
	url      string
	creation time.Time
	clicks   int
}

func connectDb(connection string) *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal("Impossible open db: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Problem with connection to the database:", err)
	}
	return db
}

func (repo *MysqlRepository) Close() {
	repo.db.Close()
}

func (repo *MysqlRepository) GetAllFeedRss() (toRet []*FeedRss) {
	rows, err := repo.db.Query("select id, name, twitterHandle, url from feed_rss")
	if err != nil {
		log.Fatal("Error running getallfeedrss query: " + err.Error())
	}

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

func (repo *MysqlRepository) FindFeedRssIdByUrl(url string) (id int) {
	log.Println("Looking for id, from url: ", url)
	err := repo.db.QueryRow("select id from feed_rss where url = ?", url).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func (repo *MysqlRepository) AddFeedRssVisited(feed_id int, updated *time.Time) {
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

func (repo *MysqlRepository) GetLastFeedRssUpdatedByFeedId(id int) (t time.Time) {
	log.Println("Looking for last update for visited: ", id)
	err := repo.db.QueryRow("select updated from feed_rss_visited where feed_id = ? order by id desc limit 1;", id).Scan(&t)
	if err != nil {
		log.Println("getLastFeedRssUpdatedByFeedId: error ", err)
		t = time.Now() // Because probably this is a new feed rss source.
	}

	return t
}

func (repo MysqlRepository) AddShortlink(id string, url string) error {
	stmt, err := repo.db.Prepare("INSERT INTO shortlink(uuid, url) VALUES(?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id, url)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
