-- phpMyAdmin SQL Dump
-- version 4.7.4
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1:3306
-- Generation Time: Mar 18, 2018 at 07:39 AM
-- Server version: 5.7.19
-- PHP Version: 5.6.31

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

SET FOREIGN_KEY_CHECKS = 0;
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `naomi`
--

-- --------------------------------------------------------

--
-- Table structure for table `authentication_tokens`
--

DROP TABLE IF EXISTS `authentication_tokens`;
CREATE TABLE IF NOT EXISTS `authentication_tokens` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `token` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` int(10) UNSIGNED NOT NULL,
  `additional` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `expire_at` datetime NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `authentication_tokens_user_id_foreign` (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
CREATE TABLE IF NOT EXISTS `comments` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `topic_id` int(10) UNSIGNED NOT NULL,
  `comment` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `comments_user_id_foreign` (`user_id`),
  KEY `comments_topic_id_foreign` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `downloads`
--

DROP TABLE IF EXISTS `downloads`;
CREATE TABLE IF NOT EXISTS `downloads` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `failed_jobs`
--

DROP TABLE IF EXISTS `failed_jobs`;
CREATE TABLE IF NOT EXISTS `failed_jobs` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `connection` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `queue` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `exception` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `failed_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `forums`
--

DROP TABLE IF EXISTS `forums`;
CREATE TABLE IF NOT EXISTS `forums` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `forums_user_id_foreign` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `friend_requests`
--

DROP TABLE IF EXISTS `friend_requests`;
CREATE TABLE IF NOT EXISTS `friend_requests` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sender` int(10) UNSIGNED NOT NULL,
  `receiver` int(10) UNSIGNED NOT NULL,
  `status` enum('pending','accepted','declined') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `friend_requests_sender_foreign` (`sender`),
  KEY `friend_requests_receiver_foreign` (`receiver`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `games`
--

DROP TABLE IF EXISTS `games`;
CREATE TABLE IF NOT EXISTS `games` (
  `gid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `game_ip` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  `game_port` int(11) NOT NULL,
  `game_version` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status_join` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status_mapname` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `players_connected` int(11) NOT NULL,
  `players_joining` int(11) NOT NULL,
  `players_max` int(11) NOT NULL DEFAULT '32',
  `team_1` int(11) NOT NULL DEFAULT '0',
  `team_2` int(11) NOT NULL DEFAULT '0',
  `team_distribution` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` date DEFAULT NULL,
  PRIMARY KEY (`gid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_heroes`
--

DROP TABLE IF EXISTS `game_heroes`;
CREATE TABLE IF NOT EXISTS `game_heroes` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `heroName` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `online` tinyint(1) NOT NULL DEFAULT '0',
  `ip_address` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_heroes_heroname_unique` (`heroName`),
  KEY `game_heroes_user_id_foreign` (`user_id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `game_heroes`
--

INSERT INTO `game_heroes` (`id`, `user_id`, `heroName`, `online`, `ip_address`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 1, 'MargeSimpson', 0, '127.0.0.1', '2017-12-10 02:40:10', '2017-12-10 02:40:10', NULL),
(2, 2, '|ccc|xSyn', 0, '127.0.0.1', '2017-12-10 02:41:47', '2017-12-10 02:41:47', NULL),
(3, 2, 'sos', 0, '127.0.0.1', '2017-12-10 02:42:01', '2017-12-10 02:42:01', NULL),
(4, 2, 'sus', 0, '127.0.0.1', '2017-12-10 02:42:13', '2017-12-10 02:42:13', NULL),
(5, 3, 'asdfasdf', 0, '127.0.0.1', '2017-12-10 02:59:33', '2017-12-10 02:59:33', NULL),
(6, 5, 'ghhh', 0, '127.0.0.1', '2018-01-30 03:36:43', '2018-01-30 03:36:43', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `game_player_regions`
--

DROP TABLE IF EXISTS `game_player_regions`;
CREATE TABLE IF NOT EXISTS `game_player_regions` (
  `userid` int(11) NOT NULL,
  `region` varchar(2) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_player_server_preferences`
--

DROP TABLE IF EXISTS `game_player_server_preferences`;
CREATE TABLE IF NOT EXISTS `game_player_server_preferences` (
  `userid` int(11) NOT NULL,
  `gid` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`userid`,`gid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_servers`
--

DROP TABLE IF EXISTS `game_servers`;
CREATE TABLE IF NOT EXISTS `game_servers` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `servername` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `secretKey` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_servers_servername_unique` (`servername`),
  UNIQUE KEY `game_servers_secretkey_unique` (`secretKey`),
  KEY `game_servers_user_id_foreign` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `game_servers`
--

INSERT INTO `game_servers` (`id`, `user_id`, `servername`, `secretKey`, `created_at`, `updated_at`) VALUES
(1, 1, 'MargeSimpson', 'MargeSimpson', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `game_server_client`
--

DROP TABLE IF EXISTS `game_server_client`;
CREATE TABLE IF NOT EXISTS `game_server_client` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `community_name` varchar(255) DEFAULT NULL,
  `ip_address` varchar(50) NOT NULL,
  `client_version` varchar(50) NOT NULL,
  `port` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=118 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `game_server_client`
--

INSERT INTO `game_server_client` (`id`, `name`, `community_name`, `ip_address`, `client_version`, `port`) VALUES
(1, '\"[iad]EA Battlefield Heroes Server(192.168.56.1:18569)\"', 'DICE', '192.168.0.10', '1.58.429030.0', 18569);
-- --------------------------------------------------------

--
-- Table structure for table `game_server_player_stats`
--

DROP TABLE IF EXISTS `game_server_player_stats`;
CREATE TABLE IF NOT EXISTS `game_server_player_stats` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `gid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `pid` int(10) UNSIGNED NOT NULL,
  `statsKey` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `statsValue` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_server_player_stats_gid_pid_statskey_unique` (`gid`,`pid`,`statsKey`),
  KEY `game_server_player_stats_pid_foreign` (`pid`)
) ENGINE=MyISAM AUTO_INCREMENT=1174 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_server_regions`
--

DROP TABLE IF EXISTS `game_server_regions`;
CREATE TABLE IF NOT EXISTS `game_server_regions` (
  `game_ip` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  `region` varchar(2) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country` varchar(2) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `weight` int(11) DEFAULT NULL,
  PRIMARY KEY (`game_ip`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_server_stats`
--

DROP TABLE IF EXISTS `game_server_stats`;
CREATE TABLE IF NOT EXISTS `game_server_stats` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `gid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `statsKey` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `statsValue` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_server_stats_gid_statskey_unique` (`gid`,`statsKey`)
) ENGINE=MyISAM AUTO_INCREMENT=9027 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `game_stats`
--

DROP TABLE IF EXISTS `game_stats`;
CREATE TABLE IF NOT EXISTS `game_stats` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `heroID` int(10) UNSIGNED NOT NULL,
  `statsKey` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `statsValue` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_stats_user_id_heroid_statskey_unique` (`user_id`,`heroID`,`statsKey`)
) ENGINE=MyISAM AUTO_INCREMENT=119 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `game_stats`
--

INSERT INTO `game_stats` (`id`, `user_id`, `heroID`, `statsKey`, `statsValue`, `created_at`, `updated_at`) VALUES
(1, 2, 1, 'level', '7', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(2, 2, 1, 'elo', '1000', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(3, 2, 1, 'c_team', '2', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(4, 2, 1, 'c_kit', '0', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(5, 2, 1, 'c_skc', '9', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(6, 2, 1, 'c_hrc', '4', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(7, 2, 1, 'c_hrs', '87', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(8, 2, 1, 'c_ft', '109', '2017-12-10 02:40:10', '2017-12-10 02:40:10'),
(9, 2, 2, 'level', '1', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(10, 2, 2, 'elo', '1000', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(11, 2, 2, 'c_team', '1', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(12, 2, 2, 'c_kit', '2', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(13, 2, 2, 'c_skc', '4', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(14, 2, 2, 'c_hrc', '2', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(15, 2, 2, 'c_hrs', '121', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(16, 2, 2, 'c_ft', '0', '2017-12-10 02:41:47', '2017-12-10 02:41:47'),
(17, 2, 3, 'level', '1', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(18, 2, 3, 'elo', '1000', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(19, 2, 3, 'c_team', '2', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(20, 2, 3, 'c_kit', '2', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(21, 2, 3, 'c_skc', '3', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(22, 2, 3, 'c_hrc', '3', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(23, 2, 3, 'c_hrs', '83', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(24, 2, 3, 'c_ft', '109', '2017-12-10 02:42:01', '2017-12-10 02:42:01'),
(25, 2, 4, 'level', '1', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(26, 2, 4, 'elo', '1000', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(27, 2, 4, 'c_team', '1', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(28, 2, 4, 'c_kit', '0', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(29, 2, 4, 'c_skc', '2', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(30, 2, 4, 'c_hrc', '2', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(31, 2, 4, 'c_hrs', '121', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(32, 2, 4, 'c_ft', '132', '2017-12-10 02:42:13', '2017-12-10 02:42:13'),
(33, 3, 5, 'level', '1', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(34, 3, 5, 'elo', '1000', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(35, 3, 5, 'c_team', '2', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(36, 3, 5, 'c_kit', '1', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(37, 3, 5, 'c_skc', '4', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(38, 3, 5, 'c_hrc', '2', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(39, 3, 5, 'c_hrs', '82', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(40, 3, 5, 'c_ft', '108', '2017-12-10 02:59:33', '2017-12-10 02:59:33'),
(41, 3, 5, 'c_apr', '4;935;936', NULL, NULL),
(42, 3, 5, 'c_emo', '5000;5007;5016;0;0;0;0;0;0', NULL, NULL),
(43, 3, 5, 'c_eqp', '3011;3009;2075;3156;2001;0;0;0;0;0', NULL, NULL),
(44, 3, 5, 'c_ltp', '9292.0000', NULL, NULL),
(45, 3, 5, 'c_wmid0', '0.0000', NULL, NULL),
(46, 3, 3, 'c_ltm', '9292.0000', NULL, NULL),
(47, 3, 3, 'c_slm', '0.0000', NULL, NULL),
(48, 3, 3, 'c_tut', '1.0000', NULL, NULL),
(49, 2, 1, 'c_apr', '4;5;201', NULL, NULL),
(50, 2, 1, 'c_emo', '5000;5001;5002;5003;5004;5005;5006;5007;5008', NULL, NULL),
(51, 2, 1, 'c_eqp', '0;0;0;0;2136;0;0;3001;0;0', NULL, NULL),
(52, 2, 1, 'c_ltp', '9289.0000', NULL, NULL),
(53, 2, 1, 'c_wmid0', '6000.0000', NULL, NULL),
(54, 2, 2, 'c_ltm', '9330.0000', NULL, NULL),
(55, 2, 2, 'c_slm', '7.0000', NULL, NULL),
(56, 2, 2, 'c_tut', '2.0000', NULL, NULL),
(57, 2, 3, 'c_apr', '4;36', NULL, NULL),
(58, 2, 4, 'c_apr', '979;981', NULL, NULL),
(59, 2, 2, 'c_apr', '9;73', NULL, NULL),
(60, 2, 2, 'c_emo', '5000;0;0;0;0;0;0;5007;0', NULL, NULL),
(61, 2, 2, 'c_eqp', '0;0;2026;0;0;0;0;0;0;0', NULL, NULL),
(62, 2, 2, 'c_ltp', '9329.0000', NULL, NULL),
(63, 2, 2, 'c_wmid0', '0.0000', NULL, NULL),
(64, 2, 2, 'mid0', '6000.0000', NULL, NULL),
(65, 2, 1, 'ds', '0.0000', NULL, NULL),
(66, 2, 1, 'ks', '0.0000', NULL, NULL),
(67, 2, 1, 'm_ct0', '0.7901', NULL, NULL),
(68, 2, 1, 'fc_los0', '1.0000', NULL, NULL),
(69, 2, 1, 'ft_los2', '1.0000', NULL, NULL),
(70, 2, 1, 'los', '1.0000', NULL, NULL),
(71, 2, 1, 'm_los0', '1.0000', NULL, NULL),
(72, 2, 2, 'ct', '8.0492', NULL, NULL),
(73, 2, 2, 'ds', '0.0000', NULL, NULL),
(74, 2, 2, 'fc_los2', '2.0000', NULL, NULL),
(75, 2, 2, 'fi', '4.0000', NULL, NULL),
(76, 2, 2, 'ft_los1', '2.0000', NULL, NULL),
(77, 2, 2, 'ks', '0.0000', NULL, NULL),
(78, 2, 2, 'los', '2.0000', NULL, NULL),
(79, 2, 2, 'm_ct0', '30.6764', NULL, NULL),
(80, 2, 2, 'm_los0', '2.0000', NULL, NULL),
(81, 2, 2, 'sw3191', '4.0000', NULL, NULL),
(82, 2, 2, 'tv5', '8.0000', NULL, NULL),
(83, 2, 2, 'tw3191', '3.0000', NULL, NULL),
(84, 2, 2, 'tw3200', '4.0000', NULL, NULL),
(85, 2, 4, 'c_emo', '5000;5007;5016;0;0;0;0;0;0', NULL, NULL),
(86, 2, 4, 'c_eqp', '3014;3002;0;2141;0;0;0;0;0;0', NULL, NULL),
(87, 2, 4, 'c_ltp', '9329.0000', NULL, NULL),
(88, 2, 4, 'c_wmid0', '0.0000', NULL, NULL),
(89, 2, 4, 'mid0', '6000.0000', NULL, NULL),
(90, 2, 3, 'c_emo', '5000;5007;5016;0;5078;0;0;0;0', NULL, NULL),
(91, 2, 3, 'c_eqp', '3003;3009;0;2026;0;0;0;0;0;0', NULL, NULL),
(92, 2, 3, 'c_ltp', '9330.0000', NULL, NULL),
(93, 2, 3, 'c_wmid0', '0.0000', NULL, NULL),
(94, 2, 3, 'ds', '0.0000', NULL, NULL),
(95, 2, 3, 'ks', '0.0000', NULL, NULL),
(96, 2, 3, 'mid0', '6000.0000', NULL, NULL),
(97, 2, 1, 'c_slm', '25284892.0000', NULL, NULL),
(98, 0, 2, 'c_slm', '6.0000', NULL, NULL),
(99, 0, 1, 'c_ltp', '9277.0000', NULL, NULL),
(100, 5, 6, 'level', '1', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(101, 5, 6, 'elo', '1000', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(102, 5, 6, 'c_team', '2', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(103, 5, 6, 'c_kit', '0', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(104, 5, 6, 'c_skc', '6', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(105, 5, 6, 'c_hrc', '5', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(106, 5, 6, 'c_hrs', '86', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(107, 5, 6, 'c_ft', '103', '2018-01-30 03:36:43', '2018-01-30 03:36:43'),
(108, 2, 4, 'ds', '0.0000', NULL, NULL),
(109, 2, 4, 'ks', '0.0000', NULL, NULL),
(110, 2, 4, 'm_ct0', '11.4642', NULL, NULL),
(111, 2, 3, 'm_ct0', '21.6262', NULL, NULL),
(112, 0, 2, 'c_ltm', '9283.0000', NULL, NULL),
(113, 3, 5, 'mid0', '-6000', NULL, NULL),
(114, 2, 2, 'c_wmid1', '0.0000', NULL, NULL),
(115, 2, 2, 'mid1', '0.0000', NULL, NULL),
(116, 2, 3, 'c_wmid1', '0.0000', NULL, NULL),
(117, 2, 3, 'edm', '75.0000', NULL, NULL),
(118, 2, 3, 'mid1', '0.0000', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `jobs`
--

DROP TABLE IF EXISTS `jobs`;
CREATE TABLE IF NOT EXISTS `jobs` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `queue` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `attempts` tinyint(3) UNSIGNED NOT NULL,
  `reserved_at` int(10) UNSIGNED DEFAULT NULL,
  `available_at` int(10) UNSIGNED NOT NULL,
  `created_at` int(10) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  KEY `jobs_queue_reserved_at_index` (`queue`,`reserved_at`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
CREATE TABLE IF NOT EXISTS `migrations` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `migration` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `migrations`
--

INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES
(1, '2014_10_12_000000_create_users_table', 1),
(2, '2014_10_12_100000_create_password_resets_table', 1),
(3, '2017_06_30_231158_create_forums_table', 1),
(4, '2017_06_30_231224_create_topics_table', 1),
(5, '2017_06_30_231302_create_comments_table', 1),
(6, '2017_07_01_114721_create_user_signatures_table', 1),
(7, '2017_07_03_090226_create_authentication_tokens_table', 1),
(8, '2017_07_04_003106_create_downloads_table', 1),
(9, '2017_07_08_212121_create_user_friends_table', 1),
(10, '2017_07_08_214941_create_friend_requests_table', 1),
(11, '2017_07_11_221348_create_roles_table', 1),
(12, '2017_07_11_221359_create_permissions_table', 1),
(13, '2017_07_11_221409_create_permission_role_table', 1),
(14, '2017_07_11_221425_create_role_user_table', 1),
(15, '2017_07_14_145823_create_news_table', 1),
(16, '2017_07_18_000000_update_users_table', 1),
(17, '2017_07_18_161536_create_jobs_table', 1),
(18, '2017_07_18_161634_create_failed_jobs_table', 1),
(19, '2017_07_20_003702_create_user_discord_table', 1),
(20, '2017_07_21_000000_update_comments_table', 1),
(21, '2017_07_21_000000_update_topics_table', 1),
(23, '2017_07_22_172358_update_user_table', 1),
(24, '2017_07_22_174548_update_user_table_for_game', 1),
(25, '2017_07_22_181120_create_game_server_table', 1),
(26, '2017_07_22_190026_create_game_heroes_table', 1),
(27, '2017_07_22_194145_create_game_stats_table', 1),
(29, '2017_07_28_061915_update_user_discords_table', 1),
(30, '2017_07_29_202605_create_game_server_stats', 1),
(31, '2017_07_29_224036_create_game_server_player_stats', 1),
(32, '2017_07_31_010137_update_game_server_player_stats', 2),
(33, '2017_08_01_030034_create_games', 2),
(34, '2017_08_02_062636_create_game_server_regions', 2),
(35, '2017_08_03_052656_update_game_server_regions', 2),
(36, '2017_08_04_022239_create_game_player_regions', 2),
(37, '2017_08_09_160449_create_game_player_server_preferences', 2);

-- --------------------------------------------------------

--
-- Table structure for table `news`
--

DROP TABLE IF EXISTS `news`;
CREATE TABLE IF NOT EXISTS `news` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `text` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` date NOT NULL,
  `user_id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `news_user_id_foreign` (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `password_resets`
--

DROP TABLE IF EXISTS `password_resets`;
CREATE TABLE IF NOT EXISTS `password_resets` (
  `email` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  KEY `password_resets_email_index` (`email`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `password_resets`
--

INSERT INTO `password_resets` (`email`, `token`, `created_at`) VALUES
('plsdonthackme@gmail.com', '$2y$10$STTxd6g2mVuqVxLXOyHGQOj97pF7w7LT.9/F.lP3qgd37AKL6PnJy', '2017-12-10 02:56:49');

-- --------------------------------------------------------

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `slug` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `permissions`
--

INSERT INTO `permissions` (`id`, `slug`, `description`) VALUES
(1, 'game.createhero', NULL),
(2, 'game.unlimitedheroes', NULL),
(3, 'game.login', NULL),
(4, 'game.matchmake', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `permission_role`
--

DROP TABLE IF EXISTS `permission_role`;
CREATE TABLE IF NOT EXISTS `permission_role` (
  `permission_id` int(10) UNSIGNED NOT NULL,
  `role_id` int(10) UNSIGNED NOT NULL,
  KEY `permission_role_permission_id_foreign` (`permission_id`),
  KEY `permission_role_role_id_foreign` (`role_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `permission_role`
--

INSERT INTO `permission_role` (`permission_id`, `role_id`) VALUES
(1, 1),
(2, 1),
(3, 1),
(4, 1);

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
CREATE TABLE IF NOT EXISTS `roles` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`id`, `title`, `slug`) VALUES
(1, 'administrator', 'administrator\r\n');

-- --------------------------------------------------------

--
-- Table structure for table `role_user`
--

DROP TABLE IF EXISTS `role_user`;
CREATE TABLE IF NOT EXISTS `role_user` (
  `user_id` int(10) UNSIGNED NOT NULL,
  `role_id` int(10) UNSIGNED NOT NULL,
  `expire_at` datetime DEFAULT NULL,
  KEY `role_user_user_id_foreign` (`user_id`),
  KEY `role_user_role_id_foreign` (`role_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `role_user`
--

INSERT INTO `role_user` (`user_id`, `role_id`, `expire_at`) VALUES
(2, 1, NULL),
(3, 1, NULL),
(5, 1, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `topics`
--

DROP TABLE IF EXISTS `topics`;
CREATE TABLE IF NOT EXISTS `topics` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` int(10) UNSIGNED NOT NULL,
  `forum_id` int(10) UNSIGNED NOT NULL,
  `text` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_comment` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `topics_user_id_foreign` (`user_id`),
  KEY `topics_forum_id_foreign` (`forum_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `birthday` date NOT NULL,
  `language` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'enUS',
  `country` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `remember_token` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `notifications` mediumtext COLLATE utf8mb4_unicode_ci,
  `ip_address` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_token` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_username_unique` (`username`),
  UNIQUE KEY `users_email_unique` (`email`),
  UNIQUE KEY `users_game_token_unique` (`game_token`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `username`, `email`, `birthday`, `language`, `country`, `password`, `remember_token`, `created_at`, `updated_at`, `notifications`, `ip_address`, `game_token`) VALUES
(1, 'MargeSimpson', 'game-server@gmail.com', '2017-12-06', 'enUS', 'Brazil', 'MargeSimpson', NULL, NULL, NULL, NULL, NULL, 'MargeSimpson'),
(2, 'lok', 'eeee@gmail.com', '1996-07-31', 'enUS', '', 'Kiop09!qas', NULL, '2017-12-10 02:37:43', '2017-12-10 02:37:43', NULL, NULL, '1234'),
(3, 'Kiop', '1234@gmail.com', '1997-07-31', 'enUS', '', '$2y$10$yqc94mokmEtovgK1oXjoOeAN2/14iacs9BKDu.qokd8NWN4GoeGM2', NULL, '2017-12-10 02:59:21', '2017-12-10 02:59:21', NULL, NULL, '123'),
(4, 'syn', 'syn@gmail.com', '2018-01-27', 'fr', 'FR', '$2y$10$ba9AYKTutmdxQgxVRbX5iOPkogUJTpaspu93COBZG.PiCtOULMfYy', NULL, '2018-01-27 07:19:31', '2018-01-27 07:19:31', NULL, NULL, 'iRzD8EvSJF9fboG'),
(5, 'caio', '22232@gmail.com', '1997-07-31', 'enUS', '', '$2y$10$M9RiC3XlhebivlXRPXeccupEZ2aajmfS4UMwiPICLeNkhsCYdI5Qe', NULL, '2018-01-30 03:11:10', '2018-01-30 04:52:54', '{\"news\":false}', '127.0.0.1', '6CUMvXUu6HNBajzR3wEtLwzfF4EgIrCJ');

-- --------------------------------------------------------

--
-- Table structure for table `user_discords`
--

DROP TABLE IF EXISTS `user_discords`;
CREATE TABLE IF NOT EXISTS `user_discords` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `discord_id` bigint(20) UNSIGNED NOT NULL,
  `discord_name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `discord_email` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `discord_discriminator` int(10) UNSIGNED DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_discords_user_id_foreign` (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `user_friends`
--

DROP TABLE IF EXISTS `user_friends`;
CREATE TABLE IF NOT EXISTS `user_friends` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `friend_id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_friends_user_id_foreign` (`user_id`),
  KEY `user_friends_friend_id_foreign` (`friend_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `user_signatures`
--

DROP TABLE IF EXISTS `user_signatures`;
CREATE TABLE IF NOT EXISTS `user_signatures` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL,
  `image` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_signatures_user_id_foreign` (`user_id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `user_signatures`
--

INSERT INTO `user_signatures` (`id`, `user_id`, `image`, `created_at`, `updated_at`) VALUES
(1, 5, '/images/signatures/PBR1iVHjCSbiDEzvnFXJshOErarNq9bJ2SDXLLmm4n7yedtE.jpg', '2018-01-30 04:56:23', '2018-01-30 04:56:23');

--
-- Constraints for dumped tables
--

--
-- Constraints for table `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `comments_topic_id_foreign` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`),
  ADD CONSTRAINT `comments_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `forums`
--
ALTER TABLE `forums`
  ADD CONSTRAINT `forums_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
SET FOREIGN_KEY_CHECKS = 1;
