-- This migrates initial-notaryserver.sql to tables that are needed for GORM

ALTER TABLE `tuf_files`
ADD COLUMN `created_at` timestamp NULL DEFAULT NULL AFTER `id`,
ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL AFTER `created_at`,
ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL AFTER `updated_at`,
MODIFY `id` int(10) unsigned AUTO_INCREMENT;

ALTER TABLE `timestamp_keys`
ADD COLUMN `id` int(10) unsigned AUTO_INCREMENT FIRST,
ADD COLUMN `created_at` timestamp NULL DEFAULT NULL AFTER `id`,
ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL AFTER `created_at`,
ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL AFTER `updated_at`,
DROP PRIMARY KEY,
ADD PRIMARY KEY (`id`),
ADD UNIQUE (`gun`);

ALTER TABLE `timestamp_keys` ADD COLUMN `role` VARCHAR(255) NOT NULL, DROP KEY `gun`, ADD UNIQUE KEY `gun_role` (`gun`, `role`);

UPDATE `timestamp_keys` SET `role`="timestamp";
