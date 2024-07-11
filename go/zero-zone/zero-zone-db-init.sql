/*
 Navicat Premium Data Transfer

 Source Server         : docker-13307
 Source Server Type    : MySQL
 Source Server Version : 50744 (5.7.44)
 Source Host           : localhost:13307
 Source Schema         : zero_zone

 Target Server Type    : MySQL
 Target Server Version : 50744 (5.7.44)
 File Encoding         : 65001

 Date: 11/07/2024 13:47:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for demo_curd
-- ----------------------------
DROP TABLE IF EXISTS `demo_curd`;
CREATE TABLE `demo_curd` (
  `firm_id` int(11) NOT NULL AUTO_INCREMENT,
  `firm_name` varchar(255) DEFAULT NULL,
  `firm_alias` varchar(255) DEFAULT NULL,
  `firm_code` varchar(255) DEFAULT NULL,
  `firm_desc` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`firm_id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of demo_curd
-- ----------------------------
BEGIN;
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, '123', '123', '23', '213213', '2024-05-27 17:56:23', '2024-05-27 17:56:23', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, '221', 'sdfdsf', '333', '', '2024-05-27 17:59:46', '2024-05-27 17:59:46', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, '213', '123', '231', '123', '2024-05-27 18:01:00', '2024-05-27 18:01:00', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (24, '324234', '234234', '234234', '234324', '2024-06-11 16:38:15', '2024-06-11 16:38:15', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (25, '432423', '423432', '42432', '423423432', '2024-06-11 16:38:19', '2024-06-11 16:38:19', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (26, '23432', '432432', '4324', '234324', '2024-06-11 16:38:23', '2024-06-11 16:38:23', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (27, 'erwer', 'werwe', '13123', 'wrwrwerewr', '2024-06-11 16:38:34', '2024-06-11 16:38:34', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (28, '12321', '3213', '1231232', '1321321321312321', '2024-06-11 16:38:39', '2024-06-11 16:38:39', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (29, '2321', '321321', '21321321', '312321321321', '2024-06-11 16:38:45', '2024-06-11 16:38:45', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (30, 'ere', '213123', '213213', '12312312321321', '2024-06-11 16:38:55', '2024-06-11 16:38:55', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (31, '123123123123', '3213213213', '213213213', '123123123', '2024-06-11 16:39:05', '2024-06-11 16:39:05', NULL);
INSERT INTO `demo_curd` (`firm_id`, `firm_name`, `firm_alias`, `firm_code`, `firm_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (32, '1111111', '11111', '1111111', '111111111', '2024-06-11 17:50:42', '2024-06-11 17:50:42', NULL);
COMMIT;

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `parent_id` int(11) unsigned NOT NULL COMMENT '父级id',
  `name` varchar(50) NOT NULL COMMENT '部门简称',
  `full_name` varchar(50) NOT NULL COMMENT '部门全称',
  `unique_key` varchar(50) NOT NULL COMMENT '唯一值',
  `type` tinyint(1) unsigned NOT NULL COMMENT '1=公司 2=子公司 3=部门',
  `status` tinyint(1) unsigned NOT NULL COMMENT '0=禁用 1=开启',
  `order_num` int(11) unsigned NOT NULL COMMENT '排序值',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='部门';

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
BEGIN;
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `full_name`, `unique_key`, `type`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (1, 0, '方舟', '方舟互联', 'arklnk', 1, 1, 0, '', '2022-08-17 02:09:17', '2022-08-22 02:13:54');
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `full_name`, `unique_key`, `type`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (2, 0, '思忆', '思忆技术', 'siyee', 1, 1, 0, '2121', '2022-08-19 06:40:10', '2024-04-26 01:40:23');
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `full_name`, `unique_key`, `type`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (3, 0, 'sdfsdffdsf', 'fullName', 'unique_k_1', 1, 1, 100, 'edddd', '2024-04-26 01:59:06', '2024-05-24 18:03:34');
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `full_name`, `unique_key`, `type`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (4, 0, 'name2', 'fullName2', 'unique_k_2', 2, 1, 100, 'edddd', '2024-04-26 02:05:42', '2024-04-26 02:05:42');
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `full_name`, `unique_key`, `type`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (5, 0, '123', '12312', '231', 1, 1, 100, '123123', '2024-05-24 18:02:01', '2024-05-24 18:02:01');
COMMIT;

-- ----------------------------
-- Table structure for sys_dictionary
-- ----------------------------
DROP TABLE IF EXISTS `sys_dictionary`;
CREATE TABLE `sys_dictionary` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '0=配置集 !0=父级id',
  `name` varchar(50) NOT NULL COMMENT '名称',
  `type` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '1文本 2数字 3数组 4单选 5多选 6下拉 7日期 8时间 9单图 10多图 11单文件 12多文件',
  `unique_key` varchar(50) NOT NULL COMMENT '唯一值',
  `value` varchar(2048) NOT NULL DEFAULT '' COMMENT '配置值',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=禁用 1=开启',
  `order_num` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='系统参数';

-- ----------------------------
-- Records of sys_dictionary
-- ----------------------------
BEGIN;
INSERT INTO `sys_dictionary` (`id`, `parent_id`, `name`, `type`, `unique_key`, `value`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (1, 0, '系统配置', 0, 'sys', '', 1, 0, '', '2022-08-22 10:03:58', '2024-06-13 15:11:33');
INSERT INTO `sys_dictionary` (`id`, `parent_id`, `name`, `type`, `unique_key`, `value`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (2, 1, '默认密码', 1, 'sys_pwd', '123456', 1, 0, '新建用户默认密码', '2022-08-22 10:03:58', '2022-08-28 08:41:39');
INSERT INTO `sys_dictionary` (`id`, `parent_id`, `name`, `type`, `unique_key`, `value`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (3, 1, '更新个人密码', 1, 'sys_ch_pwd', '', 1, 0, '', '2022-08-25 03:18:47', '2024-05-31 18:01:12');
INSERT INTO `sys_dictionary` (`id`, `parent_id`, `name`, `type`, `unique_key`, `value`, `status`, `order_num`, `remark`, `create_time`, `update_time`) VALUES (4, 1, '更新个人资料', 1, 'sys_userinfo', '', 1, 0, '', '2022-08-25 03:28:36', '2024-05-31 18:01:15');
COMMIT;

-- ----------------------------
-- Table structure for sys_job
-- ----------------------------
DROP TABLE IF EXISTS `sys_job`;
CREATE TABLE `sys_job` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(50) NOT NULL COMMENT '岗位名称',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=禁用 1=开启 ',
  `order_num` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '开启时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='工作岗位';

-- ----------------------------
-- Records of sys_job
-- ----------------------------
BEGIN;
INSERT INTO `sys_job` (`id`, `name`, `status`, `order_num`, `create_time`, `update_time`) VALUES (1, '前端', 1, 0, '2022-08-17 03:15:56', '2022-08-17 05:27:26');
INSERT INTO `sys_job` (`id`, `name`, `status`, `order_num`, `create_time`, `update_time`) VALUES (2, '后端', 1, 0, '2022-08-17 03:15:56', '2022-08-17 05:32:50');
INSERT INTO `sys_job` (`id`, `name`, `status`, `order_num`, `create_time`, `update_time`) VALUES (3, '设计', 1, 0, '2022-08-17 03:15:56', '2022-08-17 05:32:55');
COMMIT;

-- ----------------------------
-- Table structure for sys_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_log`;
CREATE TABLE `sys_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `user_id` int(11) unsigned NOT NULL COMMENT '操作账号',
  `ip` varchar(100) NOT NULL COMMENT 'ip',
  `uri` varchar(200) NOT NULL COMMENT '请求路径',
  `type` tinyint(1) unsigned NOT NULL COMMENT '1=登录日志 2=操作日志',
  `request` varchar(2048) NOT NULL DEFAULT '' COMMENT '请求数据',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=失败 1=成功',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=78 DEFAULT CHARSET=utf8mb4 COMMENT='系统日志';

-- ----------------------------
-- Records of sys_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (1, 1, '', '/admin/user/login', 1, '', 1, '2024-04-26 01:31:25', '2024-04-26 01:31:25');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (2, 2, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:28:14', '2024-04-28 01:28:14');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (3, 1, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:33:28', '2024-04-28 01:33:28');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (4, 2, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:34:40', '2024-04-28 01:34:40');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (5, 1, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:35:04', '2024-04-28 01:35:04');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (6, 1, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:36:16', '2024-04-28 01:36:16');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (7, 2, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:36:41', '2024-04-28 01:36:41');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (8, 1, '', '/admin/user/login', 1, '', 1, '2024-04-28 01:37:25', '2024-04-28 01:37:25');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (9, 1, '', '/admin/user/login', 1, '', 1, '2024-04-28 02:41:54', '2024-04-28 02:41:54');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (10, 2, '', '/admin/user/login', 1, '', 1, '2024-04-28 07:06:46', '2024-04-28 07:06:46');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (11, 1, '', '/admin/user/login', 1, '', 1, '2024-04-29 04:17:50', '2024-04-29 04:17:50');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (12, 1, '', '/admin/user/login', 1, '', 1, '2024-04-30 03:01:51', '2024-04-30 03:01:51');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (13, 1, '', '/admin/user/login', 1, '', 1, '2024-04-30 07:05:27', '2024-04-30 07:05:27');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (14, 1, '', '/admin/user/login', 1, '', 1, '2024-05-22 17:21:23', '2024-05-22 17:21:23');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (15, 1, '', '/admin/user/login', 1, '', 1, '2024-05-23 09:27:45', '2024-05-23 09:27:45');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (16, 1, '', '/admin/user/login', 1, '', 1, '2024-05-23 10:33:53', '2024-05-23 10:33:53');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (17, 1, '', '/admin/user/login', 1, '', 1, '2024-05-23 10:55:39', '2024-05-23 10:55:39');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (18, 1, '', '/admin/user/login', 1, '', 1, '2024-05-24 09:49:37', '2024-05-24 09:49:37');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (19, 1, '', '/admin/user/login', 1, '', 1, '2024-05-27 18:05:18', '2024-05-27 18:05:18');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (20, 1, '', '/admin/user/login', 1, '', 1, '2024-05-31 17:51:19', '2024-05-31 17:51:19');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (21, 1, '', '/admin/user/login', 1, '', 1, '2024-05-31 17:55:47', '2024-05-31 17:55:47');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (22, 7, '', '/admin/user/login', 1, '', 1, '2024-05-31 17:59:01', '2024-05-31 17:59:01');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (23, 7, '', '/admin/user/login', 1, '', 1, '2024-05-31 18:03:30', '2024-05-31 18:03:30');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (24, 1, '', '/admin/user/login', 1, '', 1, '2024-05-31 18:07:06', '2024-05-31 18:07:06');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (25, 1, '', '/admin/user/login', 1, '', 1, '2024-06-03 09:20:29', '2024-06-03 09:20:29');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (26, 10, '', '/admin/user/login', 1, '', 1, '2024-06-03 09:24:40', '2024-06-03 09:24:40');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (27, 10, '', '/admin/user/login', 1, '', 1, '2024-06-03 09:25:22', '2024-06-03 09:25:22');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (28, 1, '', '/admin/user/login', 1, '', 1, '2024-06-03 09:27:03', '2024-06-03 09:27:03');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (29, 1, '', '/admin/user/login', 1, '', 1, '2024-06-03 14:18:57', '2024-06-03 14:18:57');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (30, 1, '', '/admin/user/login', 1, '', 1, '2024-06-03 14:22:32', '2024-06-03 14:22:32');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (31, 1, '112.49.37.125:52783, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-03 14:27:38', '2024-06-03 14:27:38');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (32, 1, '', '/admin/user/login', 1, '', 1, '2024-06-03 14:43:00', '2024-06-03 14:43:00');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (33, 1, '112.49.37.125:55089, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-03 14:44:53', '2024-06-03 14:44:53');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (34, 1, '', '/admin/user/login', 1, '', 1, '2024-06-04 14:11:58', '2024-06-04 14:11:58');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (35, 1, '', '/admin/user/login', 1, '', 1, '2024-06-05 09:25:41', '2024-06-05 09:25:41');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (36, 1, '', '/admin/user/login', 1, '', 1, '2024-06-07 16:25:52', '2024-06-07 16:25:52');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (37, 1, '112.49.37.125:58296, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-12 15:11:04', '2024-06-12 15:11:04');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (38, 1, '112.49.37.125:63835, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-13 10:01:40', '2024-06-13 10:01:40');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (39, 12, '', '/admin/user/login', 1, '', 1, '2024-06-13 15:18:34', '2024-06-13 15:18:34');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (40, 15, '', '/admin/user/login', 1, '', 1, '2024-06-13 15:45:34', '2024-06-13 15:45:34');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (41, 12, '', '/admin/user/login', 1, '', 1, '2024-06-13 15:49:36', '2024-06-13 15:49:36');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (42, 12, '', '/admin/user/login', 1, '', 1, '2024-06-13 15:52:33', '2024-06-13 15:52:33');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (43, 12, '', '/admin/user/login', 1, '', 1, '2024-06-13 15:55:37', '2024-06-13 15:55:37');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (44, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 16:42:42', '2024-06-14 16:42:42');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (45, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 16:51:21', '2024-06-14 16:51:21');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (46, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:02:29', '2024-06-14 17:02:29');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (47, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:07:44', '2024-06-14 17:07:44');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (48, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:09:00', '2024-06-14 17:09:00');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (49, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:15:32', '2024-06-14 17:15:32');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (50, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:16:52', '2024-06-14 17:16:52');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (51, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:22:09', '2024-06-14 17:22:09');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (52, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:25:00', '2024-06-14 17:25:00');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (53, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:41:07', '2024-06-14 17:41:07');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (54, 1, '', '/admin/user/login', 1, '', 1, '2024-06-14 17:54:16', '2024-06-14 17:54:16');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (55, 1, '112.49.37.125:57769, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-20 15:48:17', '2024-06-20 15:48:17');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (56, 1, '112.49.37.125:52444, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-25 10:44:06', '2024-06-25 10:44:06');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (57, 1, '112.49.37.125:62533, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-25 11:26:14', '2024-06-25 11:26:14');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (58, 12, '', '/admin/user/login', 1, '', 1, '2024-06-25 11:33:59', '2024-06-25 11:33:59');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (59, 1, '', '/admin/user/login', 1, '', 1, '2024-06-25 11:44:47', '2024-06-25 11:44:47');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (60, 21, '', '/admin/user/login', 1, '', 1, '2024-06-25 11:45:16', '2024-06-25 11:45:16');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (61, 1, '', '/admin/user/login', 1, '', 1, '2024-06-25 11:45:52', '2024-06-25 11:45:52');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (62, 1, '112.49.37.125:50513, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-25 13:42:34', '2024-06-25 13:42:34');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (63, 22, '112.49.37.125:50580, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-25 13:43:42', '2024-06-25 13:43:42');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (64, 22, '112.49.37.125:50628, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-25 13:44:25', '2024-06-25 13:44:25');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (65, 23, '112.49.37.125:54233, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-26 10:56:48', '2024-06-26 10:56:48');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (66, 23, '112.49.37.125:54346, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-26 10:57:09', '2024-06-26 10:57:09');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (67, 1, '112.49.37.125:54390, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-26 10:57:41', '2024-06-26 10:57:41');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (68, 22, '112.49.37.125:52187, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-26 17:02:47', '2024-06-26 17:02:47');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (69, 1, '', '/admin/user/login', 1, '', 1, '2024-06-27 14:16:18', '2024-06-27 14:16:18');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (70, 1, '112.49.37.125:63539, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-06-28 16:41:10', '2024-06-28 16:41:10');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (71, 1, '112.49.37.125:63959, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-01 11:21:02', '2024-07-01 11:21:02');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (72, 1, '112.49.37.125:64098, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-01 11:23:09', '2024-07-01 11:23:09');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (73, 23, '112.49.37.125:64211, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-01 11:23:44', '2024-07-01 11:23:44');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (74, 23, '112.49.37.125:64254, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-01 11:24:16', '2024-07-01 11:24:16');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (75, 22, '112.49.37.125:55850, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-08 17:50:20', '2024-07-08 17:50:20');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (76, 1, '112.49.37.125:52223, 118.178.129.116', '/admin/user/login', 1, '', 1, '2024-07-09 14:20:44', '2024-07-09 14:20:44');
INSERT INTO `sys_log` (`id`, `user_id`, `ip`, `uri`, `type`, `request`, `status`, `create_time`, `update_time`) VALUES (77, 1, '', '/admin/user/login', 1, '', 1, '2024-07-11 05:44:05', '2024-07-11 05:44:05');
COMMIT;

-- ----------------------------
-- Table structure for sys_perm_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_perm_menu`;
CREATE TABLE `sys_perm_menu` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `name` varchar(50) NOT NULL COMMENT '名称',
  `router` varchar(1024) NOT NULL DEFAULT '' COMMENT '路由',
  `perms` varchar(1024) NOT NULL DEFAULT '' COMMENT '权限',
  `type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0=目录 1=菜单 2=权限',
  `icon` varchar(200) NOT NULL DEFAULT '' COMMENT '图标',
  `order_num` int(11) unsigned DEFAULT '0' COMMENT '排序值',
  `view_path` varchar(1024) NOT NULL DEFAULT '' COMMENT '页面路径',
  `is_show` tinyint(1) unsigned DEFAULT '1' COMMENT '0=隐藏 1=显示',
  `active_router` varchar(1024) NOT NULL DEFAULT '' COMMENT '当前激活的菜单',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `keep_alive` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8mb4 COMMENT='权限&菜单';

-- ----------------------------
-- Records of sys_perm_menu
-- ----------------------------
BEGIN;
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (1, 0, '系统管理', '/sys', '[]', 0, 'Setting', 0, '', 1, '', '2022-08-12 02:14:20', '2024-04-28 07:54:21', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (2, 1, '菜单管理', '/sys/menu', '[]', 1, 'Memo', 0, 'views/system/menu', 1, '', '2022-08-12 02:14:20', '2024-04-28 08:02:02', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (3, 2, '查询', '', '[\"sys/perm/menu/list\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (4, 2, '新增', '', '[\"sys/perm/menu/add\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (5, 2, '删除', '', '[\"sys/perm/menu/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (6, 2, '更新', '', '[\"sys/perm/menu/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (7, 1, '角色管理', '/sys/role', '[]', 1, 'Avatar', 0, 'views/system/role', 1, '', '2022-08-12 02:14:20', '2024-04-28 08:01:24', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (8, 7, '查询', '', '[\"sys/role/list\",\"sys/perm/menu/list\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (9, 7, '新增', '', '[\"sys/role/add\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (10, 7, '删除', '', '[\"sys/role/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (11, 7, '更新', '', '[\"sys/role/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (13, 1, '部门管理', '/sys/dept', '[]', 1, 'Stamp', 0, 'views/system/dept', 0, '', '2022-08-12 02:14:20', '2024-06-13 15:11:47', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (14, 13, '查询', '', '[\"sys/dept/list\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (15, 13, '新增', '', '[\"sys/dept/add\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (16, 13, '删除', '', '[\"sys/dept/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (17, 13, '更新', '', '[\"sys/dept/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (19, 18, '查询', '', '[\"sys/job/page\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (20, 18, '新增', '', '[\"sys/job/add\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (21, 18, '删除', '', '[\"sys/job/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (22, 18, '更新', '', '[\"sys/job/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (25, 24, '查询', '', '[\"sys/profession/page\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (26, 24, '新增', '', '[\"sys/profession/add\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (27, 24, '删除', '', '[\"sys/profession/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (28, 24, '更新', '', '[\"sys/profession/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (29, 1, '用户管理', '/sys/user', '[]', 1, 'User', 0, 'views/system/user', 1, '', '2022-08-12 02:14:20', '2024-06-17 09:55:45', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (30, 29, '查询', '', '[\"sys/user/page\",\"sys/dept/list\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-24 03:46:56', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (31, 29, '新增', '', '[\"sys/user/add\",\"sys/user/rdpj/info\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-24 03:17:19', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (32, 29, '删除', '', '[\"sys/user/delete\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (33, 29, '更新', '', '[\"sys/user/update\",\"sys/user/rdpj/info\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-24 03:08:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (34, 29, '改密', '', '[\"sys/user/password/update\"]', 2, '', 0, '', 1, '', '2022-08-12 02:14:20', '2022-08-25 04:51:46', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (38, 37, '查询', '', '[\"config/dict/list\",\"config/dict/data/page\"]', 2, '', 0, '', 1, '', '2022-08-22 03:42:07', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (39, 37, '新增', '', '[\"config/dict/add\"]', 2, '', 0, '', 1, '', '2022-08-22 03:42:07', '2022-08-23 09:37:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (40, 37, '删除', '', '[\"config/dict/delete\"]', 2, '', 0, '', 1, '', '2022-08-22 03:42:07', '2022-08-23 09:37:20', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (41, 37, '更新', '', '[\"config/dict/update\"]', 2, '', 0, '', 1, '', '2022-08-22 03:42:07', '2022-08-23 09:37:34', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (42, 0, '模块功能', '/feat', '[]', 0, 'Lollipop', 0, '', 1, '', '2022-08-23 04:47:23', '2024-05-28 17:29:07', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (43, 1, '登录日志', '/log/login', '[]', 1, 'Stopwatch', 0, 'views/log/login', 0, '', '2022-08-23 04:47:51', '2024-06-13 15:17:43', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (44, 43, '查询', '', '[\"log/login/page\"]', 2, '', 0, '', 1, '', '2022-08-22 03:42:07', '2022-08-23 09:35:53', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (45, 0, '文档中心', '/doc', '[]', 0, 'DocumentCopy', 0, '', 1, '', '2022-08-29 09:22:32', '2024-06-06 12:27:00', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (46, 45, 'DemoCurd', '/feat/demoCurd', '[]', 1, 'DocumentCopy', 0, 'views/feat/demo_curd', 1, '', '2022-08-29 09:29:49', '2024-05-31 11:45:38', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (50, 0, '首页', '/', '[]', 0, 'House', 10000, 'views/dashboard/index', 1, '', '2024-04-30 07:07:26', '2024-06-17 13:57:21', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (51, 45, '组件用例', '/helpSearch', '[]', 0, 'DocumentCopy', 0, 'views/help/search', 1, '', '2024-05-23 09:33:38', '2024-05-23 09:36:12', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (52, 1, '脚手架', '/sys/curd', '[]', 0, 'BrushFilled', 0, 'views/system/curd', 1, '', '2024-05-24 10:08:53', '2024-05-24 10:44:50', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (53, 42, 'featDemo', '/feat/thirdPartDevConf', '[]', 0, '', 0, 'views/feat/third_part_dev_conf', 1, '', '2024-05-24 10:45:24', '2024-07-11 05:47:05', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (54, 45, '表单验证', '/helpForm', '[]', 0, 'Check', 0, 'views/help/form', 1, '', '2024-05-27 18:02:44', '2024-05-28 17:27:14', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (57, 45, 'DemoGorm', '/feat/testGorm', '[]', 1, 'House', 100, 'views/feat/test_gorm', 1, '', '2024-05-31 18:02:51', '2024-05-31 18:13:29', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (58, 57, '查询', '', '[\"权限\"]', 2, 'House', 100, 'views/dashboard/index', 1, '', '2024-05-31 18:11:08', '2024-05-31 18:11:08', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (70, 45, 'testDel', '/abc', '[]', 0, '', 0, '', 0, '', '2024-06-13 15:48:50', '2024-06-13 15:49:48', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (71, 45, 'test2', '/test2', '[]', 1, 'House', 100, 'views/dashboard/index', 1, '', '2024-06-13 15:50:17', '2024-06-13 15:50:17', 0);
INSERT INTO `sys_perm_menu` (`id`, `parent_id`, `name`, `router`, `perms`, `type`, `icon`, `order_num`, `view_path`, `is_show`, `active_router`, `create_time`, `update_time`, `keep_alive`) VALUES (72, 1, 'Redis', '/feat/redis', '[]', 0, 'Link', 0, 'views/feat/redis', 1, '', '2024-06-14 15:32:41', '2024-06-14 15:54:33', 0);
COMMIT;

-- ----------------------------
-- Table structure for sys_profession
-- ----------------------------
DROP TABLE IF EXISTS `sys_profession`;
CREATE TABLE `sys_profession` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(50) NOT NULL COMMENT '职称',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=禁用 1=开启',
  `order_num` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='职称';

-- ----------------------------
-- Records of sys_profession
-- ----------------------------
BEGIN;
INSERT INTO `sys_profession` (`id`, `name`, `status`, `order_num`, `create_time`, `update_time`) VALUES (1, 'CEO', 1, 0, '2022-08-17 05:09:26', '2022-08-17 05:09:26');
INSERT INTO `sys_profession` (`id`, `name`, `status`, `order_num`, `create_time`, `update_time`) VALUES (2, 'CTO', 1, 0, '2022-08-17 05:09:26', '2022-08-17 05:09:26');
COMMIT;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `name` varchar(50) NOT NULL COMMENT '名称',
  `unique_key` varchar(50) NOT NULL COMMENT '唯一标识',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `perm_menu_ids` json NOT NULL COMMENT '权限集',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=禁用 1=开启',
  `order_num` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `perm_menu_ids_all` json NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='角色';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
BEGIN;
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (1, 0, '超级管理员', 'root', '超级管理员', '[]', 1, 0, '2022-08-19 02:38:19', '2024-06-05 09:28:14', '[]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (2, 0, '演示', 'demo', 'dsff vvv', '[50, 42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69]', 1, 0, '2022-08-23 13:13:05', '2024-06-13 15:20:23', '[50, 42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (3, 0, '11vvv', '22vvv', '3vvv', '[]', 1, 0, '2024-05-24 18:13:26', '2024-06-05 09:28:04', '[]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (4, 0, '22', '33', '44', '[]', 1, 0, '2024-05-24 18:13:36', '2024-06-05 09:28:08', '[]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (5, 0, '111', '2222', '3333', '[]', 1, 0, '2024-05-24 18:14:25', '2024-06-05 09:28:11', '[]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (6, 0, '12', '444', '44', '[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 29, 30, 31, 32, 33, 34, 43, 44, 52]', 1, 0, '2024-06-05 09:33:26', '2024-06-05 09:33:37', '[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 29, 30, 31, 32, 33, 34, 43, 44, 52]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (7, 0, '模块功能', 'feat', '', '[42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69]', 1, 0, '2024-06-05 09:36:42', '2024-06-13 15:54:35', '[42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (8, 0, '文档中心', 'doc', '', '[45, 57, 58, 71, 46, 51, 54, 70]', 1, 0, '2024-06-13 15:49:08', '2024-06-13 15:55:04', '[45, 57, 58, 71, 46, 51, 54, 70]');
INSERT INTO `sys_role` (`id`, `parent_id`, `name`, `unique_key`, `remark`, `perm_menu_ids`, `status`, `order_num`, `create_time`, `update_time`, `perm_menu_ids_all`) VALUES (9, 0, '全部权限', '12', '', '[50, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 29, 30, 31, 32, 33, 34, 43, 44, 52, 72, 42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69, 45, 57, 58, 71, 46, 51, 54, 70]', 1, 0, '2024-06-25 11:16:04', '2024-06-25 11:16:04', '[50, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 29, 30, 31, 32, 33, 34, 43, 44, 52, 72, 42, 53, 55, 56, 59, 62, 63, 64, 65, 66, 67, 68, 69, 45, 57, 58, 71, 46, 51, 54, 70]');
COMMIT;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `account` varchar(50) NOT NULL COMMENT '账号',
  `password` char(32) NOT NULL COMMENT '密码',
  `username` varchar(50) NOT NULL COMMENT '姓名',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(400) NOT NULL DEFAULT '' COMMENT '头像',
  `gender` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0=保密 1=女 2=男',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮件',
  `mobile` char(11) NOT NULL DEFAULT '' COMMENT '手机号',
  `profession_id` int(11) unsigned NOT NULL COMMENT '职称',
  `job_id` int(11) unsigned NOT NULL COMMENT '岗位',
  `dept_id` int(11) unsigned NOT NULL COMMENT '部门',
  `role_ids` json NOT NULL COMMENT '角色集',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '0=禁用 1=开启',
  `order_num` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`,`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COMMENT='用户';

-- ----------------------------
-- Records of sys_user
-- ----------------------------
BEGIN;
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (1, 'admin', '81955f9029a52a13608579bafe24c20b', 'admin', 'admin', 'https://avataaars.io/?clotheColor=Gray01&accessoriesType=Prescription02&avatarStyle=Circle&clotheType=GraphicShirt&eyeType=WinkWacky&eyebrowType=SadConcerned&facialHairColor=Black&facialHairType=BeardLight&hairColor=BrownDark&hatColor=Blue03&mouthType=Default&skinColor=Light&topType=LongHairDreads', 0, 'admin@gmail.com', '', 0, 0, 0, '[1]', 1, 0, 'admin', '2022-08-11 06:19:45', '2024-06-28 16:41:00', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (2, 'demo', '81955f9029a52a13608579bafe24c20b', 'demo', 'nick', 'https://avataaars.io/?avatarStyle=Circle&topType=Hat&accessoriesType=Sunglasses&facialHairType=Blank&clotheType=Hoodie&clotheColor=Heather&eyeType=Hearts&eyebrowType=UpDown&mouthType=Tongue&skinColor=DarkBrown', 0, '', '', 2, 2, 1, '[2]', 1, 0, '', '2022-08-23 14:04:24', '2024-06-03 09:26:03', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (3, '111111', '81955f9029a52a13608579bafe24c20b', '222222', '222', 'https://avataaars.io/?clotheColor=Blue02&accessoriesType=Kurt&avatarStyle=Circle&clotheType=CollarSweater&eyeType=EyeRoll&eyebrowType=RaisedExcitedNatural&facialHairColor=Platinum&facialHairType=MoustacheMagnum&hairColor=Red&hatColor=PastelYellow&mouthType=Twinkle&skinColor=Tanned&topType=WinterHat2', 0, '', '18695601234', 1, 1, 1, '[2]', 1, 0, 'dsdd', '2024-05-24 17:40:40', '2024-06-03 09:26:03', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (4, '22222', '81955f9029a52a13608579bafe24c20b', '22222', 'rrrrr', 'https://avataaars.io/?clotheColor=Red&accessoriesType=Prescription02&avatarStyle=Circle&clotheType=ShirtVNeck&eyeType=Cry&eyebrowType=FlatNatural&facialHairColor=Brown&facialHairType=BeardMedium&hairColor=PastelPink&hatColor=PastelOrange&mouthType=Concerned&skinColor=Tanned&topType=LongHairBob', 0, '', '', 1, 1, 1, '[2]', 1, 1, '111', '2024-05-24 17:48:08', '2024-06-03 09:26:03', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (5, '114444', '81955f9029a52a13608579bafe24c20b', '12222', '123456', 'https://avataaars.io/?clotheColor=PastelOrange&accessoriesType=Kurt&avatarStyle=Circle&clotheType=ShirtCrewNeck&eyeType=Side&eyebrowType=Default&facialHairColor=Platinum&facialHairType=MoustacheFancy&hairColor=Blue&hatColor=Gray02&mouthType=Grimace&skinColor=DarkBrown&topType=LongHairCurly', 0, '', '18695605011', 1, 1, 1, '[]', 0, 0, '', '2024-05-24 17:50:53', '2024-06-13 15:42:01', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (10, '999999', '20b9bfd34b2ab8d960fde4a7cc73dbe2', '999999', '999999', 'https://avataaars.io/?clotheColor=PastelOrange&accessoriesType=Blank&avatarStyle=Circle&clotheType=GraphicShirt&eyeType=Wink&eyebrowType=Default&facialHairColor=Red&facialHairType=MoustacheMagnum&hairColor=Brown&hatColor=PastelRed&mouthType=Default&skinColor=Brown&topType=ShortHairShortWaved', 0, '', '18695601234', 1, 1, 1, '[2]', 1, 0, '', '2024-06-03 09:23:26', '2024-06-06 09:36:49', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (12, 'abcdef', 'c42683217567e41b39454ed326bef838', '20240613', '', 'https://avataaars.io/?clotheColor=Gray01&accessoriesType=Prescription01&avatarStyle=Circle&clotheType=Overall&eyeType=Happy&eyebrowType=Angry&facialHairColor=Red&facialHairType=Blank&hairColor=Red&hatColor=White&mouthType=Sad&skinColor=Brown&topType=Hijab', 0, '', '', 1, 1, 1, '[8]', 1, 0, '', '2024-06-13 10:03:01', '2024-06-13 15:56:07', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (13, '23123', 'c42683217567e41b39454ed326bef838', '23123', '', 'https://avataaars.io/?clotheColor=Blue02&accessoriesType=Prescription01&avatarStyle=Circle&clotheType=ShirtCrewNeck&eyeType=EyeRoll&eyebrowType=SadConcernedNatural&facialHairColor=BlondeGolden&facialHairType=MoustacheMagnum&hairColor=Blonde&hatColor=PastelYellow&mouthType=Concerned&skinColor=Tanned&topType=LongHairBigHair', 0, '', '', 1, 1, 1, '[3]', 1, 0, '', '2024-06-13 10:03:25', '2024-06-13 15:35:53', '2024-06-13 15:35:53');
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (14, '23123', 'c42683217567e41b39454ed326bef838', '23123', '', 'https://avataaars.io/?clotheColor=PastelGreen&accessoriesType=Round&avatarStyle=Circle&clotheType=ShirtCrewNeck&eyeType=Close&eyebrowType=UpDown&facialHairColor=Black&facialHairType=Blank&hairColor=Black&hatColor=Blue03&mouthType=Vomit&skinColor=Tanned&topType=WinterHat4', 0, '', '', 1, 1, 1, '[2]', 0, 0, '', '2024-06-13 15:40:18', '2024-06-13 15:41:25', '2024-06-13 15:41:25');
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (15, 'demo2', 'c42683217567e41b39454ed326bef838', '20240613', '', 'https://avataaars.io/?clotheColor=Blue03&accessoriesType=Blank&avatarStyle=Circle&clotheType=CollarSweater&eyeType=Wink&eyebrowType=Default&facialHairColor=Auburn&facialHairType=MoustacheFancy&hairColor=SilverGray&hatColor=PastelGreen&mouthType=Vomit&skinColor=Yellow&topType=LongHairMiaWallace', 0, '', '', 1, 1, 1, 'null', 1, 0, '', '2024-06-13 15:45:13', '2024-06-13 15:46:21', '2024-06-13 15:46:21');
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (21, '11111', '496dd3f0b47328173fbde704df3f1f23', '22222', '', 'https://avataaars.io/?clotheColor=Blue03&accessoriesType=Sunglasses&avatarStyle=Circle&clotheType=GraphicShirt&eyeType=Cry&eyebrowType=AngryNatural&facialHairColor=BrownDark&facialHairType=MoustacheFancy&hairColor=PastelPink&hatColor=Black&mouthType=Smile&skinColor=Pale&topType=WinterHat1', 0, '', '', 1, 1, 1, '[]', 1, 0, '', '2024-06-25 11:22:28', '2024-06-25 11:45:06', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (22, 'zhang', '81955f9029a52a13608579bafe24c20b', 'zhangcl', 'zcl', 'https://avataaars.io/?clotheColor=PastelRed&accessoriesType=Round&avatarStyle=Circle&clotheType=ShirtScoopNeck&eyeType=Surprised&eyebrowType=UpDown&facialHairColor=Brown&facialHairType=Blank&hairColor=Platinum&hatColor=White&mouthType=Default&skinColor=Yellow&topType=WinterHat1', 0, '', '15860288204', 1, 1, 1, '[9]', 1, 0, '', '2024-06-25 13:43:13', '2024-07-08 17:51:33', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (23, 'linxy', '81955f9029a52a13608579bafe24c20b', 'linxy', '', 'https://avataaars.io/?clotheColor=Blue02&accessoriesType=Wayfarers&avatarStyle=Circle&clotheType=CollarSweater&eyeType=Default&eyebrowType=SadConcernedNatural&facialHairColor=Red&facialHairType=BeardMedium&hairColor=Blonde&hatColor=PastelRed&mouthType=Concerned&skinColor=Pale&topType=LongHairStraightStrand', 0, '', '', 1, 1, 1, '[9]', 1, 0, '', '2024-06-26 10:56:23', '2024-07-01 11:24:06', NULL);
INSERT INTO `sys_user` (`id`, `account`, `password`, `username`, `nickname`, `avatar`, `gender`, `email`, `mobile`, `profession_id`, `job_id`, `dept_id`, `role_ids`, `status`, `order_num`, `remark`, `create_time`, `update_time`, `deleted_at`) VALUES (24, '52525', 'e2d9bc66577d84ef10d1d6b46df5df67', '525', '', 'https://avataaars.io/?clotheColor=PastelOrange&accessoriesType=Blank&avatarStyle=Circle&clotheType=CollarSweater&eyeType=Default&eyebrowType=UnibrowNatural&facialHairColor=BrownDark&facialHairType=BeardMajestic&hairColor=Black&hatColor=PastelBlue&mouthType=Twinkle&skinColor=DarkBrown&topType=LongHairBigHair', 0, '', '', 1, 1, 1, 'null', 1, 0, '', '2024-07-08 10:51:46', '2024-07-08 10:51:46', NULL);
COMMIT;

-- ----------------------------
-- Table structure for test_gorm
-- ----------------------------
DROP TABLE IF EXISTS `test_gorm`;
CREATE TABLE `test_gorm` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `text` varchar(255) DEFAULT NULL COMMENT '文本',
  PRIMARY KEY (`id`),
  KEY `idx_test_gorm_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of test_gorm
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
