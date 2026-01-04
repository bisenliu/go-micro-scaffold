-- Create "common_schemas" table
CREATE TABLE `common_schemas` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (
  `id` char(36) NOT NULL COMMENT "用户ID",
  `name` varchar(50) NOT NULL COMMENT "用户名",
  `open_id` varchar(255) NOT NULL COMMENT "open_id",
  `password` varchar(100) NOT NULL COMMENT "密码",
  `phone_number` varchar(255) NOT NULL DEFAULT "" COMMENT "手机号",
  `gender` bigint NOT NULL COMMENT "性别",
  `created_at` timestamp NOT NULL COMMENT "创建时间",
  `updated_at` timestamp NOT NULL COMMENT "更新时间",
  PRIMARY KEY (`id`),
  INDEX `user_created_at` (`created_at`),
  UNIQUE INDEX `user_open_id` (`open_id`),
  UNIQUE INDEX `user_phone_number` (`phone_number`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
