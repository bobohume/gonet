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

/*Table structure for table `tbl_account` */

DROP TABLE IF EXISTS `tbl_account`;

CREATE TABLE `tbl_account` (
  `account_name` varchar(100) NOT NULL DEFAULT '' COMMENT '账号名字',
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号id',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账号状态',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `logout_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `login_ip` varchar(20) NOT NULL DEFAULT '' COMMENT '登录ip',
  PRIMARY KEY (`account_name`),
  KEY `idx_tbl_account_account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=234 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `player_id` bigint(20) NOT NULL COMMENT '玩家ID',
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号ID',
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
  KEY `idx_tbl_player_playerName` (`player_name`),
  KEY `idx_tbl_player_accountId` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  PARTITION BY HASH(player_id) PARTITIONS 3;

/*Table structure for table `tbl_player_kv` */

DROP TABLE IF EXISTS `tbl_player_kv`;

CREATE TABLE `tbl_player_kv` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `data_map` mediumblob COMMENT 'KV数据',
  PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 PARTITION BY HASH(player_id) PARTITIONS 3;

/*Data for the table `tbl_player_kv` */

/*Table structure for table `tbl_item` */

DROP TABLE IF EXISTS `tbl_item`;

CREATE TABLE `tbl_item` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `data_map` mediumblob COMMENT '物品',
  PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 PARTITION BY HASH(player_id) PARTITIONS 3;

/*Data for the table `tbl_item` */

/*Table structure for table `tbl_equip` */

DROP TABLE IF EXISTS `tbl_equip`;

CREATE TABLE `tbl_equip` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `data_map` mediumblob COMMENT '装备',
  PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 PARTITION BY HASH(player_id) PARTITIONS 3;

/*Data for the table `tbl_equip` */

/*Table structure for table `tbl_mail` */

DROP TABLE IF EXISTS `tbl_mail`;

CREATE TABLE `tbl_mail` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `data_map` mediumblob COMMENT '邮件',
  PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 PARTITION BY HASH(player_id) PARTITIONS 3;

/*Data for the table `tbl_mail` */

/*Table structure for table `tbl_social` */

DROP TABLE IF EXISTS `tbl_social`;

CREATE TABLE `tbl_social` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '角色id',
  `data_map` mediumblob COMMENT '好友',
  PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 PARTITION BY HASH(player_id) PARTITIONS 3;

/*Data for the table `tbl_social` */


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

/* Procedure structure for procedure `sp_updateplayerGold` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_updateplayerGold` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_updateplayerGold`(in playerId bigint,
in _gold int)
begin
	set @curGold = 0;
	set @err = 0;
    
    select @curGold = gold FROM tbl_player where player_Id = playerId;
    if found_rows() <> 0 then
		set @curGold = @curGold + _gold;
        update tbl_player set gold = @curGold where player_Id = playerId;
    else
		set @err = 1;
    end if;
    
    if @err <> 0 then
		rollback;
	else
		commit;
	end if;
	
   
    select @err;
end */$$
DELIMITER ;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
