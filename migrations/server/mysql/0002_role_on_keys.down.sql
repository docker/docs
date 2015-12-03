ALTER TABLE `timestamp_keys` DROP KEY `gun_role`, DROP COLUMN `role`, ADD UNIQUE KEY `gun` (`gun`);
