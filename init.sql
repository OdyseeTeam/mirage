-- this init file is not guaranteed to be up to date with godycdn schema. check the gody-cdn repo to ensure you have the right table
use mirage;
CREATE TABLE `object`
(
    `id`               bigint unsigned                  NOT NULL AUTO_INCREMENT,
    `hash`             char(64) COLLATE utf8_unicode_ci NOT NULL,
    `is_stored`        tinyint(1)                       NOT NULL DEFAULT '0',
    `length`           bigint unsigned                           DEFAULT NULL,
    `last_accessed_at` timestamp                        NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `id` (`id`),
    UNIQUE KEY `hash_idx` (`hash`),
    KEY `last_accessed_idx` (`last_accessed_at`),
    KEY `is_stored_idx` (`is_stored`)
);
CREATE TABLE `metadata`
(
    `id`             int(11)      NOT NULL AUTO_INCREMENT,
    `original_url`   text         NOT NULL,
    `godycdn_hash`   varchar(64)  NOT NULL,
    `checksum`       varchar(64)  NOT NULL,
    `original_size`  int(11)      NOT NULL,
    `otpimized_size` int(11)      NOT NULL,
    `original_mime`  varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `metadata_godycdn_hash_index` (`godycdn_hash`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;