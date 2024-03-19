alter session set container=freepdb1;
create user dvh_dmo no authentication;
create or replace package dvh_dmo.knaudit_api is
  procedure log(p_event_document in varchar2);
  end knaudit_api;
/

create or replace package body dvh_dmo.knaudit_api is
  procedure log(p_event_document in varchar2) is
    pragma autonomous_transaction;
    begin
      null;
    end log;
  end knaudit_api;
/

grant execute on dvh_dmo.knaudit_api to system;