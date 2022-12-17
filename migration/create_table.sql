
CREATE TABLE IF NOT EXISTS `test` (
                                      `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                      `name` VARCHAR(100) NOT NULL default 'unknown',
                                      `number` int(10) UNSIGNED,
                                      `passed` BOOL NOT NULL DEFAULT false,
                                      `created_at` DATETIME NOT NULL DEFAULT CURRENT_DATE,
                                      `updated_at` DATETIME NOT NULL DEFAULT CURRENT_DATE,
                                      PRIMARY KEY (`id`)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `test`;