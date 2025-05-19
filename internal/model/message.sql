CREATE TABLE `message`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `message_id` BIGINT NOT NULL comment '消息id，业务生成存入',
                            `user_id` bigint NOT NULL COMMENT '会话记录属于用户的userid',
                            `conversation_id` BIGINT NOT NULL COMMENT '会话记录id',
                            `model_id` int NOT NULL COMMENT '模型id 0-deepseek-v3,1-deepseek-r1',
                            `from_id` bigint NOT NULL COMMENT '消息发送的用户id',
                            `to_id` bigint NOT NULL COMMENT '消息接收的用户id',
                            `content` VARCHAR(255) not NULL COMMENT '消息内容',
                            `done` TINYINT not NULL comment '流消息是否停止了，0-没停止，1-停止',
                            `cur_time` datetime NOT NULL COMMENT '消息发送时间',
                            PRIMARY KEY (`id`) USING BTREE,
                            UNIQUE KEY `idx_message_id` (`message_id`) using BTREE,
                            INDEX `idx_conversation_id` (`conversation_id`) using BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
