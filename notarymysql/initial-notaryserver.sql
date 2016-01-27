DROP TABLE IF EXISTS `tuf_files`;
CREATE TABLE `tuf_files` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`gun` varchar(255) NOT NULL,
	`role` varchar(255) NOT NULL,
	`version` int(11) NOT NULL,
	`sha256` char(64) DEFAULT NULL,
	`data` longblob NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `gun` (`gun`,`role`,`version`),
	INDEX `sha256` (`sha256`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `timestamp_keys`;
CREATE TABLE `timestamp_keys` (
	`gun` varchar(255) NOT NULL,
	`cipher` varchar(50) NOT NULL,
	`public` blob NOT NULL,
	PRIMARY KEY (`gun`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
