create schema partner;

create table partner.p_partner (
    id uuid not null primary key,
    name varchar(150) not null,
    description varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null
);