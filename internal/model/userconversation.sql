CREATE TABLE `user_conversation`  (
                                      `id` bigint NOT NULL AUTO_INCREMENT,
                                      `user_id` bigint NOT NULL COMMENT '用户id',
                                      `conversation_id` bigint NOT NULL COMMENT '会话记录id',
                                      `first_time` datetime NOT NULL,
                                      `last_time` datetime NOT NULL,
                                      PRIMARY KEY (`id`) USING BTREE,
                                      INDEX `idx_user_id` (`user_id`) USING BTREE,
                                      INDEX `idx_conversation_id` (`conversation_id`) USING BTREE,
                                      INDEX `idx_user_conversation` (`user_id`,`conversation_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
