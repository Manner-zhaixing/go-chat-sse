CREATE TABLE `session`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `user_id` bigint NOT NULL COMMENT '会话记录属于用户的userid',
                            `conversation_id` bigint NOT NULL COMMENT '会话记录id',
                            `message_id` bigint NOT NULL COMMENT '消息id',
                            `cur_time` datetime NOT NULL COMMENT '时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
