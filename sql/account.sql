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
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号id',
  `account_name` varchar(100) NOT NULL DEFAULT '' COMMENT '账号名字',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账号状态',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `logout_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登出时间',
  `login_ip` varchar(20) NOT NULL DEFAULT '' COMMENT '登录ip',
  PRIMARY KEY (`account_id`),
  KEY `idx_tbl_account_account_name` (`account_name`)
) ENGINE=InnoDB AUTO_INCREMENT=234 DEFAULT CHARSET=utf8;


/*Table structure for table `tbl_player` */

DROP TABLE IF EXISTS `tbl_player`;

CREATE TABLE `tbl_player` (
  `player_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '玩家ID',
  `player_name` varchar(32) NOT NULL DEFAULT '' COMMENT '玩家名字',
  `account_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '账号ID',
  `delete_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标志',
  PRIMARY KEY (`player_id`),
  KEY `idx_tbl_player_account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=237 DEFAULT CHARSET=utf8;

/* Procedure structure for procedure `usp_activeaccount` */

/*!50003 DROP PROCEDURE IF EXISTS  `usp_activeaccount` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `usp_activeaccount`(IN _userid VARCHAR(50), IN _password VARCHAR(32), IN _uid bigint)
BEGIN
     SET @accountid = -1;
	 SET @result = '0000';
     
    IF @result = '0000' AND EXISTS(SELECT 1 FROM tbl_account A  WHERE A.account_name = _userid) THEN
		SET @result = '0002';
	END IF;
    
	IF @result = '0000' THEN
		-- 开始事务
		-- BEGIN
		SELECT @accountid = A.account_id FROM tbl_account A WHERE A.account_name = _userid;
		IF FOUND_ROWS() = 0 THEN
			-- 开始插入帐号信息
			INSERT INTO tbl_account(account_name,password, account_id)
				SELECT _userid,MD5(_password), _uid;
			IF ROW_COUNT() = 1 THEN
			  SET @accountId = _uid;
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
    SELECT @result, @accountId;
END */$$
DELIMITER ;

/* Procedure structure for procedure `usp_createplayer` */

/*!50003 DROP PROCEDURE IF EXISTS  `usp_createplayer` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`localhost` PROCEDURE `usp_createplayer`(in _accountId bigint,
in _playerName varchar(32), in _uid bigint)
begin
	set @err = -1;
    set @playerId = 0;
        
	if @err = -1 then
		insert into tbl_player(account_id, player_name, player_id)
			value(_accountId, _playerName, _uid);
		if row_count() <> 0 then
		  set @playerId = _uid;
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
    select @accountId:= A.account_id, @pwd:= A.password from tbl_account A where account_name = _userid;
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
