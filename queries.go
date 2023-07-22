package main

var CreateTableQuery = `CREATE TABLE apply_leave1 (
	leaveid serial PRIMARY KEY NOT NULL ,
	mdays double precision NOT NULL DEFAULT 0 ,
	leavetype varchar(20) NOT NULL DEFAULT '' ,
	daytype text NOT NULL '',
	leavefrom timestamp with time zone NOT NULL,
	leaveto timestamp with time zone NOT NULL,
	applieddate timestamp with time zone NOT NULL,
	leavestatus varchar(15) NOT NULL DEFAULT ''  ,
	resultdate timestamp with time zone,
	certificatestatus bool NOT NULL DEFAULT FALSE 
	certificate json NULL)`
