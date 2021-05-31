/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50649
Source Host           : localhost:3306
Source Database       : zhuque

Target Server Type    : MYSQL
Target Server Version : 50649
File Encoding         : 65001

Date: 2021-01-27 19:43:41
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `permission`
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `order_number` int(11) DEFAULT NULL,
  `href` varchar(255) DEFAULT NULL,
  `icon` varchar(255) DEFAULT NULL,
  `authority` varchar(255) DEFAULT NULL,
  `checked` int(11) DEFAULT NULL,
  `menu_type` int(11) DEFAULT NULL,
  `parent_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of permission
-- ----------------------------
INSERT INTO `permission` VALUES ('1', '系统管理', '1', '', '', '', '0', '0', '-1');
INSERT INTO `permission` VALUES ('2', '命令部署', '14', 'page/task.html', 'fa fa-adjust', 'exec.Send', '0', '0', '1');
INSERT INTO `permission` VALUES ('3', '项目管理', '15', 'page/projects.html', 'fa fa-navicon', 'project.List', '0', '0', '1');
INSERT INTO `permission` VALUES ('4', '查看项目', '0', '', '', 'project.GetProjectInfo', '0', '1', '3');
INSERT INTO `permission` VALUES ('5', '添加项目', '0', '', '', 'project.CreateProject', '0', '1', '3');
INSERT INTO `permission` VALUES ('6', '修改项目', '0', '', '', 'project.UpdateProject', '0', '1', '3');
INSERT INTO `permission` VALUES ('7', '删除项目', '0', '', '', 'project.DeleteProject', '0', '1', '3');
INSERT INTO `permission` VALUES ('8', '用户管理', '16', 'page/users.html', 'fa fa-users', 'user.List', '0', '0', '1');
INSERT INTO `permission` VALUES ('9', '查询用户', '0', '', '', 'user.GetUserInfo', '0', '1', '8');
INSERT INTO `permission` VALUES ('10', '添加用户', '0', '', '', 'user.CreateUser', '0', '1', '8');
INSERT INTO `permission` VALUES ('11', '修改用户', '0', '', '', 'user.UpdateUser', '0', '1', '8');
INSERT INTO `permission` VALUES ('12', '删除用户', '0', '', '', 'user:DeleteUser', '0', '1', '8');
INSERT INTO `permission` VALUES ('13', '角色管理', '17', 'page/role.html', 'fa fa-user-circle-o', 'role.List', '0', '0', '1');
INSERT INTO `permission` VALUES ('14', '查询角色', '0', '', '', 'role.List', '0', '1', '13');
INSERT INTO `permission` VALUES ('15', '添加角色', '0', '', '', 'role.CreateRole', '0', '1', '13');
INSERT INTO `permission` VALUES ('16', '修改角色', '0', '', '', 'role.UpdateRole', '0', '1', '13');
INSERT INTO `permission` VALUES ('17', '删除角色', '0', '', '', 'role.DeleteRole', '0', '1', '13');
INSERT INTO `permission` VALUES ('18', '权限管理', '18', 'page/menu.html', 'fa fa-list-alt', 'permission.List', '0', '0', '1');
INSERT INTO `permission` VALUES ('19', '查询权限', '0', '', '', 'authorities:view', '0', '1', '18');
INSERT INTO `permission` VALUES ('20', '添加权限', '0', '', '', 'authorities:add', '0', '1', '18');
INSERT INTO `permission` VALUES ('21', '修改权限', '0', '', '', 'authorities:edit', '0', '1', '18');
INSERT INTO `permission` VALUES ('22', '删除权限', '0', '', ' ', 'authorities:delete', '0', '1', '18');
INSERT INTO `permission` VALUES ('23', '分配项目', '0', '', '', 'user.CreateUserProject', '0', '1', '8');
INSERT INTO `permission` VALUES ('24', '分配权限', '0', '', '', 'role.SetPermissio', '0', '1', '13');
INSERT INTO `permission` VALUES ('25', '快捷发布', '11', 'page/quick-release.html', 'fa fa-bolt', null, '0', '0', '1');
INSERT INTO `permission` VALUES ('26', '快捷发布-多项目', '12', 'page/quick-release-v2.html', 'fa fa-bolt', null, '0', '0', '1');
INSERT INTO `permission` VALUES ('27', '发布记录', '13', 'page/deploy.html', 'fa fa-tasks', 'deploy.List', '0', '0', '1');

-- ----------------------------
-- Table structure for `project`
-- ----------------------------
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT '1',
  `env` varchar(255) DEFAULT NULL,
  `namespace` varchar(255) DEFAULT NULL,
  `config` text,
  `deploy_type` varchar(5) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  UNIQUE(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of project
-- ----------------------------
INSERT INTO `project` VALUES ('1', 'test-project', 1, 'fat', '测试项目', '', 'scp', '2020-09-11 12:51:42', '2020-09-11 12:51:42');

-- ----------------------------
-- Table structure for `role`
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_num` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `permissions` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES ('1', '1', '超管', '2,3,4,5,6,7,8,9,10,11,12,23,13,14,15,16,17,24,18,19,20,21,22,25,26,27');
INSERT INTO `role` VALUES ('2', '2', '管理员', '3,4,5,6,7,8,9,10,11,12,23,25,26,27');
INSERT INTO `role` VALUES ('3', '3', '开发', '2,25,26,27');

-- ----------------------------
-- Table structure for `task`
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_name` varchar(255) DEFAULT NULL,
  `project` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `release_state` int(11) DEFAULT NULL,
  `now_release` tinyint(1) DEFAULT '0',
  `approve_msg` varchar(255) DEFAULT NULL,
  `from` varchar(255) DEFAULT NULL,
  `deploy_type` varchar(5) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of task
-- ----------------------------

-- ----------------------------
-- Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `role_num` int(11) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', 'test', 'test@163.com', '098f6bcd4621d373cade4e832627b4f6', '1', '1', '2020-08-21 15:51:24', '2021-01-23 17:35:09');
-- ----------------------------
-- Table structure for `user_project`
-- ----------------------------
DROP TABLE IF EXISTS `user_project`;
CREATE TABLE `user_project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `project_id` int(11) DEFAULT NULL,
  `create_user` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of user_project
-- ----------------------------
INSERT INTO `user_project` VALUES ('1', '1', '1', '1');

