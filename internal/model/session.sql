CREATE TABLE `session`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `session_id` BIGINT NOT NULL comment '一问一答对应的sessionid，业务生成',
                            `user_id` bigint NOT NULL COMMENT '会话记录属于用户的userid',
                            `conversation_id` bigint NOT NULL COMMENT '会话记录id',
                            `message_id` bigint NOT NULL COMMENT '消息id',
                            `res_message_id` BIGINT NOT NULL comment '流消息',
                            `cur_time` datetime NOT NULL COMMENT '时间',
                            PRIMARY KEY (`id`) USING BTREE,
                            UNIQUE KEY `index_session_id` (`session_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
