/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-08-27 20:04:30
 * @LastEditTime: 2025-08-28 10:23:57
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /uiharuam/initdb.sql
 */
-- Enable foreign keys
PRAGMA FOREIGN_KEYS = ON;

BEGIN TRANSACTION;

-- Create Table DESC
CREATE TABLE `DESCRIPTION` (
    `META_DB_VERSION` INTEGER PRIMARY KEY
);

-- Create Table FILES
CREATE TABLE IF NOT EXISTS `FILES` (
    `FILE_ID` TEXT PRIMARY KEY,
    `PATH` TEXT NOT NULL,
    `CHECKSUM` TEXT NOT NULL,
    `FILE_IS_DIR` INTEGER NOT NULL, -- 0: is not dir, 1: is dir.
    `FILE_SIZE` INTEGER NOT NULL, -- as bytes.
    `ADDED_AT` INTEGER NOT NULL, -- unix time as ms.
    `TAPE_BARCODE` TEXT NOT NULL,
    `TAPE_FILE_NO` INTEGER NOT NULL
);

-- Create index
CREATE INDEX IF NOT EXISTS `IDX_BARCODE_FILE_NO` ON `FILES`(`TAPE_BARCODE`,`TAPE_FILE_NO`);

-- Insert description data
INSERT INTO `DESCRIPTION` VALUES (0);

COMMIT;