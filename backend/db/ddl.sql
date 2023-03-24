create schema partner;

create table partner.p_partner (
    id uuid not null primary key,
    name varchar(150) not null,
    description varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null
);

create table partner.p_form (
    id uuid not null primary key,
    partner_id uuid not null,
    name varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null,
    foreign key (partner_id) references partner.p_partner (id)
    on update cascade on delete restrict
);

create table partner.p_field_type (
    id uuid not null primary key,
    name varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null
);

create table partner.p_form_field (
    id uuid not null primary key,
    p_form_id uuid not null,
    p_field_type_id uuid not null,
    name varchar(150) not null,
    element varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null,
    foreign key (p_form_id) references partner.p_form (id)
    on update cascade on delete restrict,
    foreign key (p_field_type_id) references partner.p_field_type (id)
    on update cascade on delete restrict
);


create table partner.p_filled_form_field (
    id uuid not null primary key,
    p_form_field_id uuid null,
    value varchar(150) not null,
    created_at timestamp not null,
    updated_at timestamp null,
    foreign key (p_form_field_id) references partner.p_form_field (id)
    on update cascade on delete restrict
);