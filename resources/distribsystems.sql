-- phpMyAdmin SQL Dump
-- version 4.5.4.1deb2ubuntu2
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Jul 26, 2018 at 08:26 PM
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

-- --------------------------------------------------------

--
-- Table structure for table `feed_category`
--

CREATE TABLE `feed_category` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `feed_category`
--

INSERT INTO `feed_category` (`id`, `name`) VALUES
(1, 'Personal Blog'),
(2, 'Corporate News Feeds');

-- --------------------------------------------------------

--
-- Table structure for table `feed_rss`
--

CREATE TABLE `feed_rss` (
  `id` int(11) NOT NULL,
  `name` varchar(300) NOT NULL,
  `twitterHandle` varchar(120) NOT NULL,
  `url` varchar(150) NOT NULL,
  `category_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `feed_rss`
--

INSERT INTO `feed_rss` (`id`, `name`, `twitterHandle`, `url`, `category_id`) VALUES
(1, 'The Morning Paper', 'adriancolyer', 'https://blog.acolyer.org/feed/', 1),
(2, 'GCP Data & ML', 'gcpdataml', 'https://cloud.google.com/blog/big-data/feed.xml', 2),
(3, 'The Paper Trail', 'henryr', 'http://www.the-paper-trail.org/index.xml', 1);

-- --------------------------------------------------------

--
-- Table structure for table `feed_rss_visited`
--

CREATE TABLE `feed_rss_visited` (
  `id` bigint(20) NOT NULL,
  `feed_id` int(11) NOT NULL,
  `visited` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `shortlink`
--

CREATE TABLE `shortlink` (
  `uuid` varchar(50) NOT NULL,
  `url` varchar(250) NOT NULL,
  `creation` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `clicks` int(11) DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tweet`
--

CREATE TABLE `tweet` (
  `id` bigint(20) NOT NULL,
  `tweet` varchar(330) NOT NULL,
  `posted` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `published` tinyint(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `feed_category`
--
ALTER TABLE `feed_category`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

--
-- Indexes for table `feed_rss`
--
ALTER TABLE `feed_rss`
  ADD UNIQUE KEY `id_2` (`id`),
  ADD UNIQUE KEY `url` (`url`),
  ADD KEY `id` (`id`);

--
-- Indexes for table `feed_rss_visited`
--
ALTER TABLE `feed_rss_visited`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id_2` (`id`),
  ADD KEY `id` (`id`);

--
-- Indexes for table `shortlink`
--
ALTER TABLE `shortlink`
  ADD PRIMARY KEY (`uuid`),
  ADD UNIQUE KEY `uuid` (`uuid`),
  ADD UNIQUE KEY `url` (`url`);

--
-- Indexes for table `tweet`
--
ALTER TABLE `tweet`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `feed_category`
--
ALTER TABLE `feed_category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `feed_rss`
--
ALTER TABLE `feed_rss`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT for table `feed_rss_visited`
--
ALTER TABLE `feed_rss_visited`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=67;
--
-- AUTO_INCREMENT for table `tweet`
--
ALTER TABLE `tweet`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
