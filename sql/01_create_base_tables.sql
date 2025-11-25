-- ===========================================
-- ISCTF 数据库初始化 - 基础表创建
-- ===========================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1. 参赛学校表
DROP TABLE IF EXISTS `dalictf_school`;
CREATE TABLE `dalictf_school` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '学校主键ID',
  `school_name` VARCHAR(255) NOT NULL COMMENT '参赛学校名称',
  `school_admin` BIGINT(20) DEFAULT NULL COMMENT '院校负责人ID',
  `user_count` INT(11) NOT NULL DEFAULT 0 COMMENT '学校参赛人数',
  `status` ENUM('active', 'suspended') NOT NULL DEFAULT 'active' COMMENT '学校状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_school_name` (`school_name`),
  KEY `idx_school_admin` (`school_admin`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参赛学校表';

-- 2. 题目类型分类表
DROP TABLE IF EXISTS `dalictf_challenge_category`;
CREATE TABLE `dalictf_challenge_category` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '分类主键ID',
  `direction` VARCHAR(50) NOT NULL COMMENT '类型方向',
  `name_zh` VARCHAR(50) NOT NULL COMMENT '类型中文名称',
  `name_en` VARCHAR(50) NOT NULL COMMENT '类型英文名称',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '类型描述',
  `icon` VARCHAR(100) DEFAULT NULL COMMENT '图标标识',
  `color` VARCHAR(20) DEFAULT NULL COMMENT '主题颜色',
  `sort_order` INT(11) NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'active' COMMENT '分类状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_direction` (`direction`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目类型分类表';

-- 3. 系统配置表
DROP TABLE IF EXISTS `dalictf_config`;
CREATE TABLE `dalictf_config` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '配置主键ID',
  `config_key` VARCHAR(100) NOT NULL COMMENT '配置键名',
  `config_value` TEXT DEFAULT NULL COMMENT '配置键值',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置说明',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

SET FOREIGN_KEY_CHECKS = 1;
