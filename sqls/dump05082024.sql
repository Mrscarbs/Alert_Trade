-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: alert_trade_db
-- ------------------------------------------------------
-- Server version	8.0.36

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `tbl_access_token_master`
--

DROP TABLE IF EXISTS `tbl_access_token_master`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tbl_access_token_master` (
  `s_acess_token` varchar(500) NOT NULL,
  `n_expiry` bigint DEFAULT NULL,
  `n_api_id` int DEFAULT NULL,
  PRIMARY KEY (`s_acess_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_access_token_master`
--

LOCK TABLES `tbl_access_token_master` WRITE;
/*!40000 ALTER TABLE `tbl_access_token_master` DISABLE KEYS */;
INSERT INTO `tbl_access_token_master` VALUES ('ibb3F_Ty8tylG1Phj67DAWCIXUKlg54o8Kr0W9stVIPokFZ_7cajbercgYUM6cdoqbH5rOchU7mL617divYy2B8T1_yMLhfRqEic_pEY9eSzXHMdPsfc_qvQoG0Nt_KgXtIG70RI9jnJNIAX2HT6NQwju6RmjWYyPTuQRUpj3ztPg6TkqgdEGDjNNDUv1HpJGn2lJKVQddExMwlneEz9CwZ9hHqh95XQi37ft7vgHtOD93i9UoK2pOic__IbzAzSrnP9bQB7BVd3vZpc7fq1Sw',1722896999,1);
/*!40000 ALTER TABLE `tbl_access_token_master` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_api_sys_config`
--

DROP TABLE IF EXISTS `tbl_api_sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tbl_api_sys_config` (
  `n_api_id` int NOT NULL,
  `s_provider` varchar(255) DEFAULT NULL,
  `s_api_key` varchar(255) DEFAULT NULL,
  `s_secret_key` varchar(255) DEFAULT NULL,
  `n_start_time` bigint DEFAULT NULL,
  `n_last_update_time` bigint DEFAULT NULL,
  PRIMARY KEY (`n_api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_api_sys_config`
--

LOCK TABLES `tbl_api_sys_config` WRITE;
/*!40000 ALTER TABLE `tbl_api_sys_config` DISABLE KEYS */;
INSERT INTO `tbl_api_sys_config` VALUES (1,'true_data','Trial197','saum197',1722428930,1722428930);
/*!40000 ALTER TABLE `tbl_api_sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_entity_details`
--

DROP TABLE IF EXISTS `tbl_entity_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tbl_entity_details` (
  `n_entity_id` bigint NOT NULL,
  `s_name` varchar(45) DEFAULT NULL,
  `s_apikey` varchar(45) DEFAULT NULL,
  `s_email` varchar(45) DEFAULT NULL,
  `s_password` varchar(45) DEFAULT NULL,
  `n_start_timestamp` bigint DEFAULT NULL,
  `n_last_update_time` bigint DEFAULT NULL,
  PRIMARY KEY (`n_entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_entity_details`
--

LOCK TABLES `tbl_entity_details` WRITE;
/*!40000 ALTER TABLE `tbl_entity_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `tbl_entity_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_user_position_details`
--

DROP TABLE IF EXISTS `tbl_user_position_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tbl_user_position_details` (
  `n_entity_id` bigint DEFAULT NULL,
  `s_ticker` varchar(45) DEFAULT NULL,
  `n_quantity` double DEFAULT NULL,
  `n_trade_price` double DEFAULT NULL,
  `n_trade_type` int DEFAULT NULL,
  `n_trade_id` bigint NOT NULL AUTO_INCREMENT,
  `n_timestamp` bigint DEFAULT NULL,
  PRIMARY KEY (`n_trade_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_user_position_details`
--

LOCK TABLES `tbl_user_position_details` WRITE;
/*!40000 ALTER TABLE `tbl_user_position_details` DISABLE KEYS */;
INSERT INTO `tbl_user_position_details` VALUES (1,'RELIANCE',50,150.25,1,1,1722870000),(1,'RELIANCE',50,150.25,1,2,1627289197),(1,'RELIANCE',50,150.25,1,3,1627289196),(1,'RELIANCE',50,150.25,1,4,1627289193),(1,'RELIANCE',50,150.25,1,5,1627289184),(1,'RELIANCE',50,150.25,1,6,1627289164),(1,'RELIANCE',50,150.25,1,7,1627289192),(1,'TCS',50,13939,1,8,1627289192);
/*!40000 ALTER TABLE `tbl_user_position_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'alert_trade_db'
--

--
-- Dumping routines for database 'alert_trade_db'
--
/*!50003 DROP PROCEDURE IF EXISTS `stp_get_access_token_api_id` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_get_access_token_api_id`(IN api_id INT)
BEGIN
    SELECT s_acess_token, n_expiry, n_api_id
    FROM tbl_access_token_master
    WHERE n_api_id = api_id;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `stp_get_api_config` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_get_api_config`(
    IN p_api_id INT
)
BEGIN
    SELECT 
        n_api_id,
        s_provider,
        s_api_key,
        s_secret_key,
        n_start_time,
        n_last_update_time
    FROM 
        tbl_api_sys_config
    WHERE 
        n_api_id = p_api_id;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `stp_get_user_position_details` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_get_user_position_details`(
    IN p_entity_id BIGINT,
    IN p_trade_type INT,
    IN p_ticker VARCHAR(45)
)
BEGIN
    SELECT 
        n_entity_id,
        s_ticker,
        n_quantity,
        n_trade_price,
        n_trade_type,
        n_trade_id,
        n_timestamp
    FROM 
        tbl_user_position_details
    WHERE 
        n_entity_id = p_entity_id 
        AND n_trade_type = p_trade_type
        AND s_ticker = p_ticker;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `stp_insert_user_position_details` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_insert_user_position_details`(
    IN p_entity_id BIGINT,
    IN p_ticker VARCHAR(45),
    IN p_quantity DOUBLE,
    IN p_trade_price DOUBLE,
    IN p_trade_type INT,
    IN p_timestamp BIGINT
)
BEGIN
    INSERT INTO tbl_user_position_details (
        n_entity_id,
        s_ticker,
        n_quantity,
        n_trade_price,
        n_trade_type,
        n_timestamp
    ) VALUES (
        p_entity_id,
        p_ticker,
        p_quantity,
        p_trade_price,
        p_trade_type,
        p_timestamp
    );
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `stp_update_access_token` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_update_access_token`(
    IN p_api_id INT,
    IN p_access_token VARCHAR(500),
    IN p_expiry BIGINT
)
BEGIN
    UPDATE tbl_access_token_master
    SET 
        s_acess_token = p_access_token,
        n_expiry = p_expiry
    WHERE 
        n_api_id = p_api_id;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-05 15:32:10
