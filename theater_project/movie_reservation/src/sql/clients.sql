DROP TABLE IF EXISTS clients;
DROP PROCEDURE IF EXISTS while_ex;

SET @title_name := '';
SELECT title INTO @title_name FROM movies WHERE id = 1;

SET @query := CONCAT('CREATE TABLE clients (',
                     'id INT AUTO_INCREMENT, ',
                     '`', @title_name, '` SMALLINT, ',
                     'place_available VARCHAR(128) NOT NULL, ',
                     'PRIMARY KEY (id))');

PREPARE create_table_query FROM @query;
EXECUTE create_table_query;
DEALLOCATE PREPARE create_table_query;

DELIMITER $$
CREATE PROCEDURE while_ex()
BEGIN
 DECLARE i INT DEFAULT 1;
 WHILE i <= (SELECT maxticket FROM movies WHERE title = @title_name) DO
    SET @insert_query := CONCAT('INSERT INTO clients (',
                                '`', @title_name, '`, place_available) ',
                                'VALUES (',
                                i, ', "True")');
    PREPARE insert_query FROM @insert_query;
    EXECUTE insert_query;
    DEALLOCATE PREPARE insert_query;
    SET i = i + 1;
 END WHILE;
END $$
DELIMITER ;

CALL while_ex();



