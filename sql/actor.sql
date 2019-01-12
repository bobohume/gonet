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
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `tbl_mail_deleted` */

/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号ID',
  `player_id` bigint(20) NOT NULL COMMENT '玩家ID',
  `player_name` varchar(32) DEFAULT '' COMMENT '玩家名字',
  `sex` int(11) NOT NULL DEFAULT '0' COMMENT '性别',
  `level` int(11) NOT NULL DEFAULT '0' COMMENT '等级',
  `gold` int(11) NOT NULL DEFAULT '0' COMMENT '元宝',
  `draw_gold` int(11) NOT NULL DEFAULT '0' COMMENT '充值元宝',
  `vip` int(11) NOT NULL DEFAULT '0' COMMENT 'Vip等级',
  `last_login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `last_logout_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `last_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  PRIMARY KEY (`player_id`),
  KEY `idx_tbl_player_accountId` (`account_id`),
  KEY `idx_tbl_player_playerName` (`player_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `tbl_social` */

DROP TABLE IF EXISTS `tbl_social`;

CREATE TABLE `tbl_social` (
  `player_id` bigint(20) NOT NULL COMMENT '玩家id',
  `target_id` bigint(20) NOT NULL COMMENT '目标玩家id',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '关系类型',
  `friend_value` int(11) NOT NULL DEFAULT '0' COMMENT '好友度',
  PRIMARY KEY (`player_id`, `target_id`, `type`)
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
  `last_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_toprank_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* Procedure structure for procedure `sp_checkcreateplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_checkcreateplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_checkcreateplayer`(in _accountId bigint)
begin
	set @err = 0;
	select @err := case when count(player_id) >= 1 then -1 else  0 end from tbl_player where account_id = _accountId;
	select @err;
end */$$
DELIMITER ;

/* Procedure structure for procedure `sp_createplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `sp_createplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_createplayer`(in _accountId bigint,
in _playerName varchar(32),
in _sex int,
in _playerId bigint)
begin
	set @err = -1;
    set @playerId = _playerId;
    
    select 1 from tbl_player where player_id = _playerId;
    if found_rows() <> 0 then
		set @err = 1;
    end if;
        
	if @err = -1 then
		select @err := case when count(player_id) >= 1 then -3 else  -1 end from tbl_player where account_id = _accountId;
		if @err = -1 then
				insert into tbl_player(account_id, player_id, player_name, sex, level, gold, draw_gold)
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

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_updatemail`(in _mailid bigint, in _sender bigint, in _sendername varchar(32),
in _money int, in _itemid int, in _itemcount int, in _recver bigint, in _recvername varchar(32),
in _issystem tinyint, in _title varchar(128),  in _message varchar(2048))
BEGIN
	set @count = 0, @err = 0, @mailid = _mailid, @recver = _recver, @recvername = _recvername;
    
    -- 检查收件人ID是否存在alter
    if _recver <> 0 then
		select @recvername := player_name from tbl_player where player_id = _recver;
        if found_rows() = 0 then
			set @err = 1;	    -- 收件人ID不存在
        end if;
    else
		select _recver = player_id from tbl_player where player_name = _recvername;
        if found_rows() = 0 then
			set @err = 2;	    -- 收件人名称不存在
            set @recver = _recver;
        end if;
    end if;
    
    if @err = 0 then
		if _issystem = 0 then -- 非系统邮件
			if _money <> 0 or _itemid <> 0 then
				select @count = count(recver) from tbl_mail where recver=_recver AND is_system=0 AND (money<>0 OR item_id<>0);
                if @count >= 30 then
					set @err = 3;		-- 带物品邮件数量超限
				end if;
			else 
				select @count =count(recver)  from tbl_mail where recver=_recver AND is_system=0 AND (money=0 OR item_id=0);
                if @count >= 90 then
					set @err = 3;		-- 文本邮件数量超限
                end if;
            end if;
        end if;
	end if;
    
    if @err = 0 then
		if exists(select 1 from tbl_mail where id = _mailid) then
			update tbl_mail set money = _money,
								item_count = _itemcount,
                                title = _title,
                                item_id = _itemid,
                                content = _mssage
                                where id = _mailid;
		else
			insert into tbl_mail(id, sender, sender_name, money, item_id, item_count,
						send_time,recver,recver_name, is_system, title, content)
                        values(_mailid, _sender, _sendername,_money, _itemid, _itemcount,
                        current_timestamp, _recver, @recvername, _issystem, _title, _message);
        end if;
	end if;
    
    select @err, @mailid, @recver;
END */$$
DELIMITER ;

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
