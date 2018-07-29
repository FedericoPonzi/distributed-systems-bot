-- phpMyAdmin SQL Dump
-- version 4.5.4.1deb2ubuntu2
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Jul 29, 2018 at 08:41 PM
-- Server version: 5.7.22-0ubuntu0.16.04.1-log
-- PHP Version: 7.0.30-0ubuntu0.16.04.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `distribsystems`
--

--
-- Dumping data for table `feed_category`
--

INSERT INTO `feed_category` (`id`, `name`) VALUES
(1, 'Personal Blog'),
(2, 'Corporate News Feeds'),
(3, 'Engineering Blog');

--
-- Dumping data for table `feed_rss`
--

INSERT INTO `feed_rss` (`id`, `name`, `twitterHandle`, `url`, `category_id`) VALUES
(1, 'The Morning Paper', 'adriancolyer', 'https://blog.acolyer.org/feed/', 1),
(2, 'GCP Data & ML', 'gcpdataml', 'https://cloud.google.com/blog/big-data/feed.xml', 2),
(3, 'The Paper Trail', 'henryr', 'http://www.the-paper-trail.org/index.xml', 1),
(6, 'HighScalability', 'highscal', 'https://feeds.feedburner.com/HighScalability', 1),
(7, 'James Hamilton\'s Prespectives', '', 'https://perspectives.mvdirona.com/feed/', 1),
(8, 'Twitter Engineering', 'twitter', 'https://blog.twitter.com/api/blog.rss?name=engineering', 3),
(9, 'Hortonworks', 'hortonworks', 'http://hortonworks.com/blog/feed', 2),
(10, 'Denis Rystsov', 'rystsov', 'http://rystsov.info/feed.xml', 1);

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
