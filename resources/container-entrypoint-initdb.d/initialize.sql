CREATE OR REPLACE PROCEDURE dvh_dmo.knaudit_api.log(data IN BLOB) IS
BEGIN
        DBMS_OUTPUT.PUT_LINE('Mock log called with BLOB data of size: ' || DBMS_LOB.GETLENGTH(data));
        NULL;
END;
