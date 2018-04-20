/*
SQLyog Ultimate v11.52 (64 bit)
MySQL - 5.7.17-log : Database - md_account
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`md_account` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `md_account`;

/*Table structure for table `tbl_account` */

DROP TABLE IF EXISTS `tbl_account`;

CREATE TABLE `tbl_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增量',
  `accountId` int(11) NOT NULL DEFAULT '0' COMMENT '账号id',
  `accountName` varchar(100) NOT NULL DEFAULT '' COMMENT '账号名字',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账号状态',
  `loginTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `logoutTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `loginIp` varchar(20) NOT NULL DEFAULT '' COMMENT '登录ip',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_account_accountName` (`accountName`),
  KEY `idx_tbl_account_accountId` (`accountId`)
) ENGINE=InnoDB AUTO_INCREMENT=234 DEFAULT CHARSET=utf8;

/*Data for the table `tbl_account` */

insert  into `tbl_account`(`id`,`accountId`,`accountName`,`password`,`status`,`loginTime`,`logoutTime`,`loginIp`) values (232,50000232,'test66666','e10adc3949ba59abbe56e057f20f883e',0,'2018-01-18 20:53:13','2018-01-18 20:53:13',''),(233,10000233,'test166666','e10adc3949ba59abbe56e057f20f883e',0,'2018-01-19 12:08:55','2018-01-19 12:08:55','');

/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增量',
  `playerId` int(11) NOT NULL DEFAULT '0' COMMENT '玩家ID',
  `playerName` varchar(32) NOT NULL DEFAULT '' COMMENT '玩家名字',
  `accountId` int(11) NOT NULL DEFAULT '0' COMMENT '账号ID',
  `deleteFlag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标志',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_player_accountId` (`accountId`),
  KEY `idx_tbl_player_playerId` (`playerId`)
) ENGINE=InnoDB AUTO_INCREMENT=237 DEFAULT CHARSET=utf8;

/*Data for the table `tbl_player` */

insert  into `tbl_player`(`id`,`playerId`,`playerName`,`accountId`,`deleteFlag`) values (229,50000229,'我是大坏蛋11',228,1),(230,50000230,'我是大坏蛋11',229,0),(231,50000231,'我是大坏蛋11',230,0),(232,50000232,'我是大坏蛋11',230,0),(233,50000233,'我是大坏蛋11',230,0),(234,50000234,'我是大坏蛋11',230,0),(235,50000235,'我是大坏蛋11',231,1),(236,50000236,'我是大坏蛋11((',50000232,0);

/* Procedure structure for procedure `usp_activeaccount` */

/*!50003 DROP PROCEDURE IF EXISTS  `usp_activeaccount` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `usp_activeaccount`(IN _userid VARCHAR(50), IN _password VARCHAR(32))
BEGIN
     SET @accountid = -1;
	 SET @result = '0000';
     
    IF @result = '0000' AND EXISTS(SELECT 1 FROM tbl_account A  WHERE A.accountName = _userid) THEN
		SET @result = '0002';
	END IF;
    
	IF @result = '0000' THEN
		-- 开始事务
		-- BEGIN
		SELECT @accountid = A.accountid FROM tbl_account A WHERE A.accountName = _userid;
		IF FOUND_ROWS() = 0 THEN
			-- 开始插入帐号信息
			INSERT INTO tbl_account(accountName,password)		
				SELECT _userid,MD5(_password);
			IF ROW_COUNT() = 1 THEN
				SET @accountid = 10000000 + @@IDENTITY;
                update tbl_account set accountid = @accountid where id=@@IDENTITY;
			ELSE
				SET @result = '0003';
			END IF;
		END IF;
		IF @result = '0000' THEN
			COMMIT;
		ELSE
			ROLLBACK;
		END IF;
	
	END IF;
    SELECT @result, @accountid;
END */$$
DELIMITER ;

/* Procedure structure for procedure `usp_createplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `usp_createplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `usp_createplayer`(in _accountId int,
in _playerName varchar(32))
begin
	set @err = -1;
    set @playerId = 0;
        
	if @err = -1 then
		insert into tbl_player(accountId, playerName)
			value(_accountId, _playerName);
		if row_count() <> 0 then
			set @playerid = 10000000 + @@IDENTITY;
			update tbl_player set playerId = @playerId where id=@@IDENTITY;
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

/* Procedure structure for procedure `usp_login` */

/*!50003 DROP PROCEDURE IF EXISTS  `usp_login` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `usp_login`(in _userid varchar(50), in _password varchar(32))
BEGIN
    set @result = '0000';
    set @accountId = 0;
    set @pwd = '';
    select @accountId:= A.accountId, @pwd:= A.password from tbl_account A where accountName = _userid;
    if found_rows() = 0 then
		set @result = '0001';
	else
		if @result = '0000' then
			if md5(_password) <> @pwd then
				set @result = '0002';
			end if;
        end if;
    end if;
    
    if @result = '0000' then
		commit;
	else
		rollback;
    end if;
    
    select @result, @accountId;
END */$$
DELIMITER ;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
