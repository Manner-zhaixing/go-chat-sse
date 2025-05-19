CREATE TABLE `user`  (
                         `id` bigint NOT NULL AUTO_INCREMENT,
                         `username` varchar(255) NOT NULL,
                         `password` varchar(255) NOT NULL,
                         `conversation_nums` BIGINT NOT NULL COMMENT '每个用户对应的会话记录数量' ,
                         `register_time` datetime NOT NULL,
                         `last_login_time` datetime NOT NULL,
                         PRIMARY KEY (`id`) USING BTREE,
                         UNIQUE INDEX idx_username (username)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic AUTO_INCREMENT = 100;
