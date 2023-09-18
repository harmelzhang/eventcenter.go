-- MySQL

CREATE DATABASE IF NOT EXISTS `eventcenter` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `eventcenter`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `topic`;
CREATE TABLE `topic`  (
    `id` CHAR(36) NOT NULL COMMENT '主键',
    `name` VARCHAR(100) NOT NULL COMMENT '主题名称',
    `create_time` DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主题信息';

DROP TABLE IF EXISTS `endpoint`;
CREATE TABLE `endpoint`  (
    `id` CHAR(36) NOT NULL COMMENT '主键',
    `server_name` VARCHAR(50) NOT NULL COMMENT '服务名称',
    `is_micro` TINYINT NOT NULL COMMENT '是否微服务',
    `topic_id` CHAR(36) NOT NULL COMMENT '监听主题',
    `type` VARCHAR(50) NOT NULL COMMENT '事件类型',
    `protocol` VARCHAR(10) NOT NULL COMMENT '处理协议（HTTP、TCP、RPC）',
    `endpoint` VARCHAR(255) NOT NULL COMMENT '终端地址',
    `register_time` DATETIME NOT NULL DEFAULT NOW() COMMENT '注册时间',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_endpoint_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topic` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='终端注册信息';

DROP TABLE IF EXISTS `event`;
CREATE TABLE `event`  (
    `id` CHAR(36) NOT NULL COMMENT '主键',
    `source` VARCHAR(100) NOT NULL COMMENT '事件源',
    `topic_id` CHAR(36) NOT NULL COMMENT '所属主题',
    `type` VARCHAR(100) NOT NULL COMMENT '事件类型',
    `data` TEXT NOT NULL COMMENT '上下文数据',
    `create_time` DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
    `cloudevents` TEXT NOT NULL COMMENT '完整CloudEvents规范数据',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_event_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topic` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件信息';
