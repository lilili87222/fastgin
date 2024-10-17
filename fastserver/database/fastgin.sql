
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_api`;
CREATE TABLE `sys_api`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'请求方式\'',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'访问路径\'',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'所属类别\'',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'说明\'',
  `creator` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'创建人\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_api_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_api
-- ----------------------------
INSERT INTO `sys_api` VALUES (1, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/user/logout', 'user', '用户登出', '系统');
INSERT INTO `sys_api` VALUES (2, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/user/refreshToken', 'user', '刷新JWT令牌', '系统');
INSERT INTO `sys_api` VALUES (3, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/user/info', 'user', '获取当前登录用户信息', '系统');
INSERT INTO `sys_api` VALUES (4, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/user/index', 'user', '获取用户列表', '系统');
INSERT INTO `sys_api` VALUES (5, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PUT', '/user/changePwd', 'user', '更新用户登录密码', '系统');
INSERT INTO `sys_api` VALUES (6, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/user/index', 'user', '创建用户', '系统');
INSERT INTO `sys_api` VALUES (7, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/user/index/:userId', 'user', '更新用户', '系统');
INSERT INTO `sys_api` VALUES (8, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'DELETE', '/user/index', 'user', '批量删除用户', '系统');
INSERT INTO `sys_api` VALUES (9, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/role/index', 'role', '获取角色列表', '系统');
INSERT INTO `sys_api` VALUES (10, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/role/index', 'role', '创建角色', '系统');
INSERT INTO `sys_api` VALUES (11, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/role/index/:roleId', 'role', '更新角色', '系统');
INSERT INTO `sys_api` VALUES (12, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/role/menus/:roleId', 'role', '获取角色的权限菜单', '系统');
INSERT INTO `sys_api` VALUES (13, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/role/menus/:roleId', 'role', '更新角色的权限菜单', '系统');
INSERT INTO `sys_api` VALUES (14, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/role/apis/:roleId', 'role', '获取角色的权限接口', '系统');
INSERT INTO `sys_api` VALUES (15, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/role/apis/:roleId', 'role', '更新角色的权限接口', '系统');
INSERT INTO `sys_api` VALUES (16, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'DELETE', '/role/index', 'role', '批量删除角色', '系统');
INSERT INTO `sys_api` VALUES (17, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/menu/index', 'menu', '获取菜单列表', '系统');
INSERT INTO `sys_api` VALUES (18, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/menu/tree', 'menu', '获取菜单树', '系统');
INSERT INTO `sys_api` VALUES (19, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/menu/index', 'menu', '创建菜单', '系统');
INSERT INTO `sys_api` VALUES (20, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/menu/index/:menuId', 'menu', '更新菜单', '系统');
INSERT INTO `sys_api` VALUES (21, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'DELETE', '/menu/index', 'menu', '批量删除菜单', '系统');
INSERT INTO `sys_api` VALUES (22, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/menu/user/:userId', 'menu', '获取用户的可访问菜单列表', '系统');
INSERT INTO `sys_api` VALUES (23, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/menu/user_tree/:userId', 'menu', '获取用户的可访问菜单树', '系统');
INSERT INTO `sys_api` VALUES (24, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/api/index', 'api', '获取接口列表', '系统');
INSERT INTO `sys_api` VALUES (25, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/api/tree', 'api', '获取接口树', '系统');
INSERT INTO `sys_api` VALUES (26, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'POST', '/api/index', 'api', '创建接口', '系统');
INSERT INTO `sys_api` VALUES (27, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'PATCH', '/api/index/:roleId', 'api', '更新接口', '系统');
INSERT INTO `sys_api` VALUES (28, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'DELETE', '/api/index', 'api', '批量删除接口', '系统');
INSERT INTO `sys_api` VALUES (29, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'GET', '/log/index', 'log', '获取操作日志列表', '系统');
INSERT INTO `sys_api` VALUES (30, '2024-10-14 09:10:10.378', '2024-10-14 09:10:10.378', NULL, 'DELETE', '/log/index', 'log', '批量删除操作日志', '系统');
INSERT INTO `sys_api` VALUES (31, '2024-10-14 12:31:43.848', '2024-10-14 12:31:43.848', NULL, 'GET', '/system/info', 'system', '获取操作系统和本软件的信息', 'admin');
INSERT INTO `sys_api` VALUES (32, '2024-10-14 12:32:15.689', '2024-10-14 12:32:15.689', NULL, 'GET', '/system/stop', 'system', '停止web服务', 'admin');
INSERT INTO `sys_api` VALUES (33, '2024-10-14 12:32:35.894', '2024-10-14 15:05:20.354', NULL, 'GET', '/system/restart', 'system', '重启web服务', 'admin');

-- ----------------------------
-- Table structure for sys_casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `sys_casbin_rule`;
CREATE TABLE `sys_casbin_rule`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_casbin_rule`(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 400 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_casbin_rule
-- ----------------------------
INSERT INTO `sys_casbin_rule` VALUES (99, 'p', 'admin', '/api/index', 'DELETE', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (95, 'p', 'admin', '/api/index', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (97, 'p', 'admin', '/api/index', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (98, 'p', 'admin', '/api/index/:roleId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (96, 'p', 'admin', '/api/tree', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (101, 'p', 'admin', '/log/index', 'DELETE', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (100, 'p', 'admin', '/log/index', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (92, 'p', 'admin', '/menu/index', 'DELETE', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (88, 'p', 'admin', '/menu/index', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (90, 'p', 'admin', '/menu/index', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (91, 'p', 'admin', '/menu/index/:menuId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (89, 'p', 'admin', '/menu/tree', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (93, 'p', 'admin', '/menu/user/:userId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (94, 'p', 'admin', '/menu/user_tree/:userId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (85, 'p', 'admin', '/role/apis/:roleId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (86, 'p', 'admin', '/role/apis/:roleId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (87, 'p', 'admin', '/role/index', 'DELETE', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (80, 'p', 'admin', '/role/index', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (81, 'p', 'admin', '/role/index', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (82, 'p', 'admin', '/role/index/:roleId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (83, 'p', 'admin', '/role/menus/:roleId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (84, 'p', 'admin', '/role/menus/:roleId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (102, 'p', 'admin', '/system/info', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (104, 'p', 'admin', '/system/restart', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (103, 'p', 'admin', '/system/stop', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (76, 'p', 'admin', '/user/changePwd', 'PUT', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (79, 'p', 'admin', '/user/index', 'DELETE', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (75, 'p', 'admin', '/user/index', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (77, 'p', 'admin', '/user/index', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (78, 'p', 'admin', '/user/index/:userId', 'PATCH', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (74, 'p', 'admin', '/user/info', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (72, 'p', 'admin', '/user/logout', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (73, 'p', 'admin', '/user/refreshToken', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (109, 'p', 'guest', '/menu/access/tree/:userId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (108, 'p', 'guest', '/user/changePwd', 'PUT', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (107, 'p', 'guest', '/user/info', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (105, 'p', 'guest', '/user/logout', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (106, 'p', 'guest', '/user/refreshToken', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (30, 'p', 'user', '/menu/access/tree/:userId', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (8, 'p', 'user', '/user/info', 'GET', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (2, 'p', 'user', '/user/logout', 'POST', '', '', '');
INSERT INTO `sys_casbin_rule` VALUES (5, 'p', 'user', '/user/refreshToken', 'POST', '', '', '');

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'菜单名称(英文名, 可用于国际化)\'',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'菜单标题(无法国际化时使用)\'',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'菜单图标\'',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'菜单访问路径\'',
  `redirect` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'重定向路径\'',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'前端组件路径\'',
  `sort` int(3) UNSIGNED NULL DEFAULT 999 COMMENT '\'菜单顺序(1-999)\'',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '\'菜单状态(正常/禁用, 默认正常)\'',
  `hidden` tinyint(1) NULL DEFAULT 2 COMMENT '\'菜单在侧边栏隐藏(1隐藏，2显示)\'',
  `no_cache` tinyint(1) NULL DEFAULT 2 COMMENT '\'菜单是否被 <keep-alive> 缓存(1不缓存，2缓存)\'',
  `always_show` tinyint(1) NULL DEFAULT 2 COMMENT '\'忽略之前定义的规则，一直显示根路由(1忽略，2不忽略)\'',
  `breadcrumb` tinyint(1) NULL DEFAULT 1 COMMENT '\'面包屑可见性(可见/隐藏, 默认可见)\'',
  `active_menu` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'在其它路由时，想在侧边栏高亮的路由\'',
  `parent_id` bigint(20) UNSIGNED NULL DEFAULT 0 COMMENT '\'父菜单编号(编号为0时表示根菜单)\'',
  `creator` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'创建人\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_menu_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'System', '系统管理', 'component', '/system', '/system/user', 'Layout', 10, 1, 2, 2, 2, 1, NULL, 0, '系统');
INSERT INTO `sys_menu` VALUES (2, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'User', '用户管理', 'user', 'user', NULL, '/system/user/index', 11, 1, 2, 2, 2, 1, NULL, 1, '系统');
INSERT INTO `sys_menu` VALUES (3, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'Role', '角色管理', 'peoples', 'role', NULL, '/system/role/index', 12, 1, 2, 2, 2, 1, NULL, 1, '系统');
INSERT INTO `sys_menu` VALUES (4, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'Menu', '菜单管理', 'tree-table', 'menu', NULL, '/system/menu/index', 13, 1, 2, 2, 2, 1, NULL, 1, '系统');
INSERT INTO `sys_menu` VALUES (5, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'Api', '接口管理', 'tree', 'api', NULL, '/system/api/index', 14, 1, 2, 2, 2, 1, NULL, 1, '系统');
INSERT INTO `sys_menu` VALUES (6, '2024-10-14 09:10:10.083', '2024-10-14 09:10:10.083', NULL, 'Log', '日志管理', 'example', '/log', '/log/operation-log', 'Layout', 20, 1, 2, 2, 2, 1, NULL, 0, '系统');
INSERT INTO `sys_menu` VALUES (7, '2024-10-14 09:10:10.083', '2024-10-15 10:05:44.069', NULL, 'OperationLog', '系统日志', 'documentation', 'operation-log', '', '/log/operation-log/index', 21, 1, 2, 2, 2, 1, '', 6, 'admin');

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `user_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'用户登录名\'',
  `ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'Ip地址\'',
  `ip_location` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'Ip所在地\'',
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'请求方式\'',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'访问路径\'',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'说明\'',
  `status` int(4) NULL DEFAULT NULL COMMENT '\'响应状态码\'',
  `start_time` datetime(3) NULL DEFAULT NULL COMMENT '\'发起时间\'',
  `time_cost` int(6) NULL DEFAULT NULL COMMENT '\'请求耗时(ms)\'',
  `user_agent` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'浏览器标识\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_operation_log_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `keyword` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '\'1正常, 2禁用\'',
  `sort` int(3) NULL DEFAULT 999 COMMENT '\'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)\'',
  `creator` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_sys_role_name`(`name`) USING BTREE,
  UNIQUE INDEX `uni_sys_role_keyword`(`keyword`) USING BTREE,
  INDEX `idx_sys_role_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '2024-10-14 09:10:10.076', '2024-10-15 09:58:08.056', NULL, '管理员', 'admin', '', 1, 1, '系统');
INSERT INTO `sys_role` VALUES (2, '2024-10-14 09:10:10.076', '2024-10-14 09:10:10.076', NULL, '普通用户', 'user', '', 1, 3, '系统');
INSERT INTO `sys_role` VALUES (3, '2024-10-14 09:10:10.076', '2024-10-15 10:03:59.862', NULL, '访客', 'guest', '', 1, 5, '系统');

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `menu_id` bigint(20) UNSIGNED NOT NULL,
  `role_id` bigint(20) UNSIGNED NOT NULL,
  PRIMARY KEY (`menu_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (1, 1);
INSERT INTO `sys_role_menu` VALUES (2, 1);
INSERT INTO `sys_role_menu` VALUES (3, 1);
INSERT INTO `sys_role_menu` VALUES (4, 1);
INSERT INTO `sys_role_menu` VALUES (5, 1);
INSERT INTO `sys_role_menu` VALUES (6, 1);
INSERT INTO `sys_role_menu` VALUES (6, 2);
INSERT INTO `sys_role_menu` VALUES (7, 1);
INSERT INTO `sys_role_menu` VALUES (7, 2);

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `user_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `nick_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `introduction` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '\'1正常, 2禁用\'',
  `creator` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_sys_user_mobile`(`mobile`) USING BTREE,
  UNIQUE INDEX `uni_sys_user_username`(`user_name`) USING BTREE,
  INDEX `idx_sys_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES (1, '2024-10-14 09:10:10.364', '2024-10-14 13:43:03.044', NULL, 'admin', '$2a$10$4Ugz/XI60xb6awSKJRDk5.kEUs5.eOpg0dVmXrhxoGVkMhHR7A3H.', '18888888888', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', '', '', 1, 'admin');
INSERT INTO `sys_user` VALUES (2, '2024-10-14 09:10:10.364', '2024-10-14 09:10:10.364', NULL, 'faker', '$2a$10$5GznAHkxjzUYGQsoTiNmpuy9Scs.8d1DBnKWLmexrYvVo50NnUlOW', '19999999999', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', '', '', 1, '系统');
INSERT INTO `sys_user` VALUES (3, '2024-10-14 09:10:10.364', '2024-10-14 09:10:10.364', NULL, 'nike', '$2a$10$kKk4kresgzIfukFx9vlE.ecesq8Sh11Qk1D0pskzq8wZfWOuGLQxa', '13333333333', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', '', '', 1, '系统');
INSERT INTO `sys_user` VALUES (4, '2024-10-14 09:10:10.364', '2024-10-14 09:10:10.364', NULL, 'bob', '$2a$10$DX1BFE154XFVu3rEkQ6Uy.8ENfU/vsaXreeUGPNusUw2o.0r1D/eq', '15555555555', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', '', '', 1, '系统');

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
  `role_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  PRIMARY KEY (`role_id`, `user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO `sys_user_role` VALUES (1, 1);
INSERT INTO `sys_user_role` VALUES (1, 2);
INSERT INTO `sys_user_role` VALUES (2, 2);
INSERT INTO `sys_user_role` VALUES (2, 3);
INSERT INTO `sys_user_role` VALUES (3, 4);

-- ----------------------------
-- Table structure for web_dictionary
-- ----------------------------
DROP TABLE IF EXISTS `web_dictionary`;
CREATE TABLE `web_dictionary`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `created_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  `key` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '配置的Key',
  `value` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '配置的值',
  `desc` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `key索引`(`key`) USING BTREE COMMENT 'key必须保持唯一性'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of web_dictionary
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
