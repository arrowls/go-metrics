create table if not exists gauges (
    name varchar(50) not null primary key,
    value double precision not null
);

create table if not exists counters (
  name varchar(50) not null unique primary key,
  value bigint not null
);