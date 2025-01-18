# Host: 127.0.0.1  (Version: 5.7.26)
# Date: 2025-01-18 21:09:44
# Generator: MySQL-Front 5.3  (Build 4.234)

/*!40101 SET NAMES utf8 */;

#
# Structure for table "gbl_book"
#

DROP TABLE IF EXISTS `gbl_book`;
CREATE TABLE `gbl_book` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `book_name` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  `vision` varchar(255) NOT NULL DEFAULT 'true',
  `hash` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

#
# Data for table "gbl_book"
#

/*!40000 ALTER TABLE `gbl_book` DISABLE KEYS */;
INSERT INTO `gbl_book` VALUES (1,'资本论','text','true','54f01d0b9a4dcdeabe98b54877eed5173d033398926412ccf787a58786c1b221'),(2,'共产党宣言','text','true','ac0644da80b6102085b0e6a39a235e73bb54e99b93cb6997e04131a5f2b6c19e'),(3,'家庭、私有制和国家的起源','text','true','e36eb7acf1cb58fdde3e600c566c8975ecc6cb302d34cbc22834b9b66e922b40'),(4,'资本','text','true','2e5149644da1f5b55bf025fa39f6bf0c339e96a08b53344355b793ca00e0e682');
/*!40000 ALTER TABLE `gbl_book` ENABLE KEYS */;

#
# Structure for table "gbl_chapter"
#

DROP TABLE IF EXISTS `gbl_chapter`;
CREATE TABLE `gbl_chapter` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `book_id` varchar(255) DEFAULT NULL,
  `hash` varchar(255) NOT NULL DEFAULT '',
  `file_list` longtext NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "gbl_chapter"
#

/*!40000 ALTER TABLE `gbl_chapter` DISABLE KEYS */;
/*!40000 ALTER TABLE `gbl_chapter` ENABLE KEYS */;

#
# Structure for table "gbl_config"
#

DROP TABLE IF EXISTS `gbl_config`;
CREATE TABLE `gbl_config` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(255) NOT NULL DEFAULT '',
  `value` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

#
# Data for table "gbl_config"
#

/*!40000 ALTER TABLE `gbl_config` DISABLE KEYS */;
INSERT INTO `gbl_config` VALUES (1,'guest-allow','true');
/*!40000 ALTER TABLE `gbl_config` ENABLE KEYS */;

#
# Structure for table "gbl_storage"
#

DROP TABLE IF EXISTS `gbl_storage`;
CREATE TABLE `gbl_storage` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `chapter_id` int(11) NOT NULL DEFAULT '0',
  `type` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "gbl_storage"
#

/*!40000 ALTER TABLE `gbl_storage` DISABLE KEYS */;
/*!40000 ALTER TABLE `gbl_storage` ENABLE KEYS */;

#
# Structure for table "gbl_token"
#

DROP TABLE IF EXISTS `gbl_token`;
CREATE TABLE `gbl_token` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `token` varchar(255) NOT NULL DEFAULT '',
  `dietime` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

#
# Data for table "gbl_token"
#

/*!40000 ALTER TABLE `gbl_token` DISABLE KEYS */;
INSERT INTO `gbl_token` VALUES (1,1,'3a11e70b44ef8eac411de1523811a93be255b1bb7eedc036f61925c750ade47b',1730993191);
/*!40000 ALTER TABLE `gbl_token` ENABLE KEYS */;

#
# Structure for table "gbl_user"
#

DROP TABLE IF EXISTS `gbl_user`;
CREATE TABLE `gbl_user` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `admin` varchar(255) NOT NULL DEFAULT 'false',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

#
# Data for table "gbl_user"
#

/*!40000 ALTER TABLE `gbl_user` DISABLE KEYS */;
INSERT INTO `gbl_user` VALUES (1,'admin','8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918','true'),(2,'noadmin','fbe03dc7f00d059debe445169f331bba6d217008c91a6e98678556eef11ed85a','false');
/*!40000 ALTER TABLE `gbl_user` ENABLE KEYS */;
