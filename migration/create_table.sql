
CREATE TABLE IF NOT EXISTS `test` (
                                      `id` int(10) UNSIGNED AUTO_INCREMENT,
                                      `name` VARCHAR(100) default 'unknown',
                                      `number` int(10) UNSIGNED,
                                      `passed` BOOL DEFAULT false,
                                      `created_at` DATETIME DEFAULT NOW(),
                                      `updated_at` DATETIME DEFAULT NOW(),
                                      PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS `test`;