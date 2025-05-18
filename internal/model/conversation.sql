CREATE TABLE `conversation`  (
                                 `id` bigint NOT NULL AUTO_INCREMENT,
                                 `user_id` BIGINT NOT NULL comment '大会话属于用户的userid',
                                 `message_nums` BIGINT NOT NULL comment '大会话的消息数量',
                                 `first_time` datetime NOT NULL comment '大会话的创建时间',
                                 `last_time` datetime NOT NULL comment '大会话的最后一次更新时间',
                                 PRIMARY KEY (`id`) USING BTREE

) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;


