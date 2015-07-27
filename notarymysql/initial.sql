DROP TABLE IF EXISTS `tuf_files`;
CREATE TABLE `tuf_files` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`gun` varchar(255) NOT NULL,
	`role` varchar(255) NOT NULL,
	`version` int(11) NOT NULL,
	`data` longblob NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `gun` (`gun`,`role`,`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `timestamp_keys`;
CREATE TABLE `timestamp_keys` (
	`gun` varchar(255) NOT NULL,
	`cipher` varchar(50) NOT NULL,
	`public` blob NOT NULL,
	PRIMARY KEY (`gun`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `private_keys`;
CREATE TABLE `private_keys` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`created_at` datetime NOT NULL,
	`updated_at` datetime NOT NULL,
	`deleted_at` datetime DEFAULT NULL,
	`key_id`  varchar(255) NOT NULL,
	`encryption_alg`  varchar(255) NOT NULL,
	`keywrap_alg`  varchar(255) NOT NULL,
	`algorithm`  varchar(50) NOT NULL,
	`passphrase_alias`  varchar(50) NOT NULL,
	`public`  blob NOT NULL,
	`private`  blob NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE (`key_id`),
	UNIQUE (`key_id`,`encryption`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
