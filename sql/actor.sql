/*
SQLyog Ultimate v11.52 (64 bit)
MySQL - 5.7.17-log : Database - md_actor
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`md_actor` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `md_actor`;

/*Table structure for table `tbl_mail` */

DROP TABLE IF EXISTS `tbl_mail`;

CREATE TABLE `tbl_mail` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '邮件ID',
  `sender` int(11) NOT NULL DEFAULT '0' COMMENT '发送者',
  `sendername` varchar(32) NOT NULL DEFAULT '' COMMENT '发送者名字',
  `recver` int(11) NOT NULL DEFAULT '0' COMMENT '接收者',
  `recvername` varchar(32) NOT NULL DEFAULT '' COMMENT '接收者名字',
  `money` int(11) NOT NULL DEFAULT '0' COMMENT '金钱',
  `itemid` int(11) NOT NULL DEFAULT '0' COMMENT '物品ID',
  `itemcount` int(11) NOT NULL DEFAULT '0' COMMENT '物品个数',
  `isread` tinyint(4) NOT NULL DEFAULT '0' COMMENT '读取标志',
  `issystem` tinyint(4) NOT NULL DEFAULT '0' COMMENT '系统标志',
  `recvflag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '领取标志',
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '邮件标题',
  `content` blob NOT NULL COMMENT '邮件内容',
  `sendtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_mail_recver` (`recver`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;

/*Data for the table `tbl_mail` */

insert  into `tbl_mail`(`id`,`sender`,`sendername`,`recver`,`recvername`,`money`,`itemid`,`itemcount`,`isread`,`issystem`,`recvflag`,`title`,`content`,`sendtime`) values (14,50000055,'',50000055,'test',1000,60010,10,0,1,0,'test','test1111','2017-12-25 16:40:19'),(15,50000055,'',50000055,'test',1000,60010,10,0,1,0,'test','test1111','2017-12-25 16:41:04'),(16,50000055,'',50000055,'test',1000,60010,10,0,1,0,'test','test1111','2017-12-25 16:41:18'),(17,50000055,'',50000055,'test',1000,60010,10,0,1,0,'test','test1111','2017-12-25 16:41:33');

/*Table structure for table `tbl_mail_deleted` */

DROP TABLE IF EXISTS `tbl_mail_deleted`;

CREATE TABLE `tbl_mail_deleted` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '邮件ID',
  `sender` int(11) NOT NULL DEFAULT '0' COMMENT '发送者',
  `sendername` varchar(32) NOT NULL DEFAULT '' COMMENT '发送者名字',
  `recver` int(11) NOT NULL DEFAULT '0' COMMENT '接收者',
  `recvername` varchar(32) NOT NULL DEFAULT '' COMMENT '接收者名字',
  `money` int(11) NOT NULL DEFAULT '0' COMMENT '金钱',
  `itemid` int(11) NOT NULL DEFAULT '0' COMMENT '物品ID',
  `itemcount` int(11) NOT NULL DEFAULT '0' COMMENT '物品个数',
  `isread` tinyint(4) NOT NULL DEFAULT '0' COMMENT '读取标志',
  `issystem` tinyint(4) NOT NULL DEFAULT '0' COMMENT '系统标志',
  `recvflag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '领取标志',
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '邮件标题',
  `content` blob NOT NULL COMMENT '邮件内容',
  `sendtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `deletetime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_mail_deleted_recver` (`recver`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `tbl_mail_deleted` */

/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `accountId` int(11) NOT NULL DEFAULT '0' COMMENT '账号ID',
  `playerId` int(11) NOT NULL COMMENT '玩家ID',
  `playerName` varchar(32) DEFAULT '' COMMENT '玩家名字',
  `sex` int(11) NOT NULL DEFAULT '0' COMMENT '性别',
  `level` int(11) NOT NULL DEFAULT '0' COMMENT '等级',
  `gold` int(11) NOT NULL DEFAULT '0' COMMENT '元宝',
  `drawGold` int(11) NOT NULL DEFAULT '0' COMMENT '充值元宝',
  `vip` int(11) NOT NULL DEFAULT '0' COMMENT 'Vip等级',
  `lastLoginTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `lastLogoutTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `lastUpdateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleteTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  PRIMARY KEY (`playerId`),
  KEY `idx_tbl_player_accountId` (`accountId`),
  KEY `idx_tbl_player_playerName` (`playerName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `tbl_player` */

insert  into `tbl_player`(`accountId`,`playerId`,`playerName`,`sex`,`level`,`gold`,`drawGold`,`vip`,`lastLoginTime`,`lastLogoutTime`,`lastUpdateTime`,`deleteTime`) values (10000233,10000237,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 21:01:59','2018-01-18 21:01:59','2018-01-18 21:01:59','2018-01-18 21:01:59'),(203,50000047,'test',0,0,0,0,0,'2017-12-19 13:18:11','2017-12-19 13:18:11','2017-12-19 13:18:11','2017-12-19 13:18:11'),(204,50000048,'test',0,0,0,0,0,'2017-12-19 14:12:43','2017-12-19 14:12:43','2017-12-19 14:12:43','2017-12-19 14:12:43'),(205,50000049,'test',0,0,0,0,0,'2017-12-19 14:23:43','2017-12-19 14:23:43','2017-12-19 14:23:43','2017-12-19 14:23:43'),(207,50000050,'test',0,0,0,0,0,'2017-12-19 14:39:21','2017-12-19 14:39:21','2017-12-19 14:39:21','2017-12-19 14:39:21'),(208,50000051,'test',0,0,0,0,0,'2017-12-19 15:02:33','2017-12-19 15:02:33','2017-12-19 15:02:33','2017-12-19 15:02:33'),(211,50000052,'test',0,0,0,0,0,'2017-12-19 15:07:50','2017-12-19 15:07:50','2017-12-19 15:07:50','2017-12-19 15:07:50'),(212,50000053,'test',0,0,0,0,0,'2017-12-19 16:46:34','2017-12-19 16:46:34','2017-12-19 16:46:34','2017-12-19 16:46:34'),(213,50000054,'test',0,0,0,0,0,'2017-12-19 17:05:35','2017-12-19 17:05:35','2017-12-19 17:05:35','2017-12-19 17:05:35'),(214,50000055,'test',0,0,0,0,0,'2017-12-19 17:27:56','2017-12-19 17:27:56','2017-12-19 17:27:56','2017-12-19 17:27:56'),(215,50000056,'test',0,0,0,0,0,'2017-12-19 17:32:18','2017-12-19 17:32:18','2017-12-19 17:32:18','2017-12-19 17:32:18'),(216,50000057,'我是大坏蛋',0,0,0,0,0,'2018-01-05 10:14:02','2018-01-05 10:14:02','2018-01-05 10:14:02','2018-01-05 10:14:02'),(217,50000058,'我是大坏蛋',0,0,0,0,0,'2018-01-05 10:16:15','2018-01-05 10:16:15','2018-01-05 10:16:15','2018-01-05 10:16:15'),(218,50000059,'我是大坏蛋',0,0,0,0,0,'2018-01-05 10:21:00','2018-01-05 10:21:00','2018-01-05 10:21:00','2018-01-05 10:21:00'),(219,50000060,'我是大坏蛋11',0,0,0,0,0,'2018-01-05 10:26:03','2018-01-05 10:26:03','2018-01-05 10:26:03','2018-01-05 10:26:03'),(223,50000061,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 14:45:05','2018-01-18 14:45:05','2018-01-18 14:45:05','2018-01-18 14:45:05'),(221,50000223,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 14:15:11','2018-01-18 14:15:11','2018-01-18 14:15:11','2018-01-18 14:15:11'),(222,50000224,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 14:36:22','2018-01-18 14:36:22','2018-01-18 14:36:22','2018-01-18 14:36:22'),(224,50000225,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 14:58:55','2018-01-18 14:58:55','2018-01-18 14:58:55','2018-01-18 14:58:55'),(225,50000226,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 15:08:55','2018-01-18 15:08:55','2018-01-18 15:08:55','2018-01-18 15:08:55'),(226,50000227,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 15:17:32','2018-01-18 15:17:32','2018-01-18 15:17:32','2018-01-18 15:17:32'),(227,50000228,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 15:20:23','2018-01-18 15:20:23','2018-01-18 15:20:23','2018-01-18 15:20:23'),(228,50000229,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 15:25:00','2018-01-18 15:25:00','2018-01-18 15:25:00','2018-01-18 15:25:00'),(229,50000230,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 15:27:16','2018-01-18 15:27:16','2018-01-18 15:27:16','2018-01-18 15:27:16'),(230,50000234,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 16:00:13','2018-01-18 16:00:13','2018-01-18 16:00:13','2018-01-18 16:00:13'),(231,50000235,'我是大坏蛋11',0,0,0,0,0,'2018-01-18 16:01:14','2018-01-18 16:01:14','2018-01-18 16:01:14','2018-01-18 16:01:14'),(50000232,50000236,'我是大坏蛋11((',0,0,0,0,0,'2018-01-18 20:53:13','2018-01-18 20:53:13','2018-01-18 20:53:13','2018-01-18 20:53:13');

/*Table structure for table `tbl_social` */

DROP TABLE IF EXISTS `tbl_social`;

CREATE TABLE `tbl_social` (
  `PlayerId` int(11) NOT NULL COMMENT '玩家id',
  `Type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '关系类型',
  `FriendValue` int(11) NOT NULL DEFAULT '0' COMMENT '好友度',
  PRIMARY KEY (`PlayerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
  `lastTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_toprank_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `tbl_toprank` */

insert  into `tbl_toprank`(`id`,`type`,`name`,`score`,`value0`,`value1`,`lastTime`) values (2,0,'test',1000,0,0,'1970-01-01 08:00:00'),(3,0,'test',200,0,0,'1970-01-01 08:00:00'),(4,0,'test',600,0,0,'1970-01-01 08:00:00'),(5,0,'test',1000,0,0,'1970-01-01 08:00:00'),(7,0,'test',1000,0,0,'1970-01-01 08:00:00'),(100,0,'test',5000,0,0,'1970-01-01 08:00:00'),(200,0,'test',4000,0,0,'1970-01-01 08:00:00'),(300,0,'test',6000,0,0,'1970-01-01 08:00:00'),(500,0,'test',8000,0,0,'1970-01-01 08:00:00'),(600,0,'test',10000,0,0,'1970-01-01 08:00:00'),(700,0,'test',9000,0,0,'1970-01-01 08:00:00');

/* Procedure structure for procedure `sp_checkcreateplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_checkcreateplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_checkcreateplayer`(in _accountId int)
begin
	set @err = 0;
	select @err := case when count(playerId) >= 1 then -1 else  0 end from tbl_player where accountId = _accountId;
	select @err;
end */$$
DELIMITER ;

/* Procedure structure for procedure `sp_createplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_createplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_createplayer`(in _accountId int,
in _playerName varchar(32),
in _sex int,
in _playerId int)
begin
	set @err = -1;
    set @playerId = 0;
    
    select 1 from tbl_player where playerId = _playerId;
    if found_rows() <> 0 then
		set @err = 1;
    end if;
        
	if @err = -1 then
		select @err := case when count(playerId) >= 1 then -3 else  -1 end from tbl_player where accountId = _accountId;
		if @err = -1 then
				set @playerid = _playerid;
				insert into tbl_player(accountId, playerId, playerName, sex, level, gold, drawGold)
							values(_accountId, @playerId, _playerName, _sex, 0,		0,		0);
				set @err = 0;
		end if;
	end if;
	
    if @err <> 0 then
		rollback;
	else
		commit;
	end if;
        
	select @err, @playerId;
end */$$
DELIMITER ;

/* Procedure structure for procedure `sp_updatemail` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_updatemail` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_updatemail`(in _mailid int, in _sender int, in _sendername varchar(32),
in _money int, in _itemid int, in _itemcount int, in _recver int, in _recvername varchar(32),
in _issystem tinyint, in _title varchar(128),  in _message varchar(2048))
BEGIN
	set @count = 0, @err = 0, @mailid = _mailid, @recver = _recver, @recvername = _recvername;
    
    -- 检查收件人ID是否存在alter
    if _recver <> 0 then
		select @recvername := playername from tbl_player where playerid = _recver;
        if found_rows() = 0 then
			set @err = 1;	    -- 收件人ID不存在
        end if;
    else
		select _recver = playerid from tbl_player where playername = _recvername;
        if found_rows() = 0 then
			set @err = 2;	    -- 收件人名称不存在
            set @recver = _recver;
        end if;
    end if;
    
    if @err = 0 then
		if _issystem = 0 then -- 非系统邮件
			if _money <> 0 or _itemid <> 0 then
				select @count = count(recver) from tbl_mail where recver=_recver AND isSystem=0 AND (money<>0 OR itemid<>0);
                if @count >= 30 then
					set @err = 3;		-- 带物品邮件数量超限
				end if;
			else 
				select @count =count(recver)  from tbl_mail where recver=_recver AND isSystem=0 AND (money=0 OR itemid=0);
                if @count >= 90 then
					set @err = 3;		-- 文本邮件数量超限
                end if;
            end if;
        end if;
	end if;
    
    if @err = 0 then
		if exists(select 1 from tbl_mail where id = _mailid) then
			update tbl_mail set money = _money,
								itemcount = _itemcount,
                                title = _title,
                                itemid = _itemid,
                                content = _mssage
                                where id = _mailid;
		else
			insert into tbl_mail(sender, sendername, money, itemid, itemcount, 
						sendtime,recver,recvername, issystem, title, content)
                        values(_sender, _sendername,_money, _itemid, _itemcount,
                        current_timestamp, _recver, @recvername, _issystem, _title, _message);
                        set _mailid = @@IDENTITY;
                        set @mailid = _mailid;
        end if;
	end if;
    
    select @err, @mailid, @recver;
END */$$
DELIMITER ;

/* Procedure structure for procedure `sp_updateplayerGold` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_updateplayerGold` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_updateplayerGold`(in playerId int,
in _gold int)
begin
	set @curGold = 0;
	set @err = 0;
    
    select @curGold = gold FROM tbl_player where playerId = playerId;
    if found_rows() <> 0 then
		set @curGold = @curGold + _gold;
        update tbl_player set gold = @curGold where playerId = playerId;
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
