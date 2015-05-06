# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.23)
# Database: yorm_test
# Generation Time: 2015-05-06 16:21:30 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table golang_word
# ------------------------------------------------------------

DROP TABLE IF EXISTS `golang_word`;

CREATE TABLE `golang_word` (
  `aid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(32) NOT NULL DEFAULT '',
  `rate` float DEFAULT NULL,
  PRIMARY KEY (`aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `golang_word` WRITE;
/*!40000 ALTER TABLE `golang_word` DISABLE KEYS */;

INSERT INTO `golang_word` (`aid`, `word`, `rate`)
VALUES
	(1,'go',1.02),
	(2,'import',0.99),
	(3,'main',0.88);

/*!40000 ALTER TABLE `golang_word` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table program_language
# ------------------------------------------------------------

DROP TABLE IF EXISTS `program_language`;

CREATE TABLE `program_language` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT NULL,
  `rank_month` date DEFAULT NULL,
  `position` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `program_language` WRITE;
/*!40000 ALTER TABLE `program_language` DISABLE KEYS */;

INSERT INTO `program_language` (`id`, `name`, `rank_month`, `position`, `created`)
VALUES
	(1,'Java','2015-04-30',1,'2015-04-30 16:39:00'),
	(2,'C','2015-04-30',2,'2015-04-30 16:39:00'),
	(3,'C++','2015-04-30',3,'2015-04-30 16:39:00'),
	(4,'Go','2015-04-30',42,'2015-04-30 16:39:00'),
	(103,'PHP','2015-05-07',12,'2015-05-07 00:03:46');

/*!40000 ALTER TABLE `program_language` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
