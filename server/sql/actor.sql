/*
SQLyog Ultimate v11.52 (64 bit)
MySQL - 5.7.17-log : Database - md_actor
*********************************************************************
*/


/*!40101 SET NAMES utf8mb4 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`md_actor` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `md_actor`;

/*Table structure for table `tbl_mail` */

DROP TABLE IF EXISTS `tbl_mail`;

CREATE TABLE `tbl_mail` (
  `id` bigint(20) NOT NULL  DEFAULT '0' COMMENT '邮件ID',
  `sender` bigint(20) NOT NULL DEFAULT '0' COMMENT '发送者',
  `sender_name` varchar(32) NOT NULL DEFAULT '' COMMENT '发送者名字',
  `recver` bigint(20) NOT NULL DEFAULT '0' COMMENT '接收者',
  `recver_name` varchar(32) NOT NULL DEFAULT '' COMMENT '接收者名字',
  `money` int(11) NOT NULL DEFAULT '0' COMMENT '金钱',
  `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '物品ID',
  `item_count` int(11) NOT NULL DEFAULT '0' COMMENT '物品个数',
  `is_read` tinyint(4) NOT NULL DEFAULT '0' COMMENT '读取标志',
  `is_system` tinyint(4) NOT NULL DEFAULT '0' COMMENT '系统标志',
  `recv_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '领取标志',
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '邮件标题',
  `content` text NOT NULL COMMENT '邮件内容',
  `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_mail_recver` (`recver`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `tbl_mail_deleted` */

DROP TABLE IF EXISTS `tbl_mail_deleted`;

CREATE TABLE `tbl_mail_deleted` (
  `id` bigint(20) NOT NULL COMMENT '邮件ID',
  `sender` bigint(20) NOT NULL DEFAULT '0' COMMENT '发送者',
  `sender_name` varchar(32) NOT NULL DEFAULT '' COMMENT '发送者名字',
  `recver` bigint(20) NOT NULL DEFAULT '0' COMMENT '接收者',
  `recver_name` varchar(32) NOT NULL DEFAULT '' COMMENT '接收者名字',
  `money` int(11) NOT NULL DEFAULT '0' COMMENT '金钱',
  `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '物品ID',
  `item_count` int(11) NOT NULL DEFAULT '0' COMMENT '物品个数',
  `is_read` tinyint(4) NOT NULL DEFAULT '0' COMMENT '读取标志',
  `is_system` tinyint(4) NOT NULL DEFAULT '0' COMMENT '系统标志',
  `recv_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '领取标志',
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '邮件标题',
  `content` text NOT NULL COMMENT '邮件内容',
  `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_mail_deleted_recver` (`recver`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Data for the table `tbl_mail_deleted` */

/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号ID',
  `player_id` bigint(20) NOT NULL COMMENT '玩家ID',
  `player_name` varchar(32) DEFAULT '' COMMENT '玩家名字',
  `sex` int(11) NOT NULL DEFAULT '0' COMMENT '性别',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '等级',
  `gold` int(11) NOT NULL DEFAULT '0' COMMENT '元宝',
  `draw_gold` int(11) NOT NULL DEFAULT '0' COMMENT '充值元宝',
  `vip` int(11) NOT NULL DEFAULT '0' COMMENT 'Vip等级',
  `last_login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `last_logout_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `update_time` timestamp ON UPDATE CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`player_id`),
  KEY `idx_tbl_player_accountId` (`account_id`),
  KEY `idx_tbl_player_playerName` (`player_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Table structure for table `tbl_social` */

DROP TABLE IF EXISTS `tbl_social`;

CREATE TABLE `tbl_social` (
  `player_id` bigint(20) NOT NULL COMMENT '玩家id',
  `target_id` bigint(20) NOT NULL COMMENT '目标玩家id',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '关系类型',
  `friend_value` int(11) NOT NULL DEFAULT '0' COMMENT '好友度',
  PRIMARY KEY (`player_id`, `target_id`, `type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Data for the table `tbl_social` */

/*Table structure for table `tbl_player_kv` */

DROP TABLE IF EXISTS `tbl_player_kv`;

CREATE TABLE `tbl_player_kv` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `key` int(11) NOT NULL DEFAULT '0' COMMENT '记录类型',
  `value` bigint(20) NOT NULL DEFAULT '0' COMMENT '记录值',
  PRIMARY KEY (`player_id`,`key`),
  KEY `idx_tbl_player_kv_player_id` (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Data for the table `tbl_player_kv` */

/*Table structure for table `tbl_toprank` */

DROP TABLE IF EXISTS `tbl_toprank`;

CREATE TABLE `tbl_toprank` (
  `id` bigint(20) unsigned NOT NULL COMMENT 'ID',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '排行榜类型',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '排行榜名字',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '分数',
  `value0` int(11) NOT NULL DEFAULT '0' COMMENT '附加信息',
  `value1` int(11) NOT NULL DEFAULT '0' COMMENT '附加信息',
  `last_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_toprank_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Data for the table `tbl_toprank` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
