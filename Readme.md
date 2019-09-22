## Distributed Systems bot
@distribsystems

This bot has the following aims:

 1. Tweet links from feeds of interesting websites. [Done]
 1. Retweet and like "interesting" tweets
 3. Schedule posts based on a publication timetable (with exceptions - like "publish now"), made from a Telegram bot interface.
 4. Save every tweet in the database, to do some analysis - based on interactions. [Done]
 5. A web based interface to let users review the sources of news.
 6. Use a shortlink system while publishing tweets, to track clicks and understand 
        which feeds are more interesting for the users. [Done]

## 1. Feed RSS
The feed rss is now powered by the feeds_handler module.


## 3. Telegram bot + post at
Why telegram? Because I like the option to schedule a post by just sending a message on Telegram.
I already use it to send links through telegram, but it's not possible to schedule tweets (yet).

---

## Deployment
A systemd service for the telegram bot + crontab for recurring fetches.
Example of crontab setup:
 * `run.sh`:
```bash
#!/bin/bash
/home/bots/distributed-systems-bot/bin/distributed-systems-bot --fetch-rss --config /home/bots/distributed-systems-bot/bin/config.yaml 2>&1 > /home/bots/distributed-systems-bot/logs/crontab-fetch-rss.out
```
 * example crontab:
```
0 * * * *       /home/bots/distributed-systems-bot/run.sh
```
