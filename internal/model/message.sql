CREATE TABLE `message`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `user_id` bigint NOT NULL COMMENT '会话记录属于用户的userid',
                            `conversation_id` BIGINT NOT NULL COMMENT '会话记录id',
                            `model_id` int NOT NULL COMMENT '模型id 0-deepseek-v3,1-deepseek-r1',
                            `from_id` bigint NOT NULL COMMENT '消息发送的用户id',
                            `to_id` bigint NOT NULL COMMENT '消息接收的用户id',
                            `content` VARCHAR(255) not NULL COMMENT '消息内容',
                            `cur_time` datetime NOT NULL COMMENT '消息发送时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
