create table if not exists groups
(
    id          bigserial,
    name        text                     not null,
    permissions bigint                   not null,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone,
    constraint groups_pkey
        primary key (id)
);

create table if not exists cities
(
    id         bigserial,
    name       text                     not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    constraint cities_pk
        primary key (id)
);

create table if not exists users
(
    id            bigserial,
    surname       text                     not null,
    name          text                     not null,
    patronymic    text                     not null,
    date_of_birth date                     not null,
    phone_number  text                     not null,
    email         text                     not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone,
    city_id       bigint                   not null,
    constraint users_pkey
        primary key (id),
    constraint users_city_id_fkey
        foreign key (city_id) references cities
);

create table if not exists sessions
(
    refresh_token text                     not null,
    expires_at    timestamp with time zone not null,
    user_id       bigint,
    constraint sessions_pkey
        primary key (refresh_token),
    constraint sessions_user_id_fkey
        foreign key (user_id) references users
);

create table if not exists donor_companies
(
    id              bigserial,
    name            text                     not null,
    contract_date   date                     not null,
    contract_number integer                  not null,
    created_at      timestamp with time zone not null,
    updated_at      timestamp with time zone,
    city_id         bigint                   not null,
    constraint donor_companies_pkey
        primary key (id),
    constraint donor_companies_city_id_fkey
        foreign key (city_id) references cities
);

create table if not exists acts
(
    id               bigserial,
    user_id          bigint                   not null,
    donor_company_id bigint                   not null,
    created_at       timestamp with time zone not null,
    updated_at       timestamp with time zone,
    constraint acts_pkey
        primary key (id),
    constraint acts_user_id_fkey
        foreign key (user_id) references users,
    constraint acts_donor_company_id_fkey
        foreign key (donor_company_id) references donor_companies
);

create table if not exists act_contents
(
    id              bigserial,
    act_id          bigint                   not null,
    number          integer                  not null,
    name            text                     not null,
    count           integer                  not null,
    price           integer                  not null,
    expiration_date date                     not null,
    comment         text                     not null,
    created_at      timestamp with time zone not null,
    updated_at      timestamp with time zone,
    constraint act_contents_pkey
        primary key (id),
    constraint act_contents_act_id_fkey
        foreign key (act_id) references acts
);

create table if not exists users_to_groups
(
    user_id  bigint not null,
    group_id bigint not null,
    constraint users_to_groups_pkey
        primary key (user_id, group_id),
    constraint users_to_groups_user_id_fkey
        foreign key (user_id) references users,
    constraint users_to_groups_group_id_fkey
        foreign key (group_id) references groups
);

create unique index if not exists cities_id_uindex
    on cities (id);

create table if not exists files
(
    id           bigserial,
    user_id      bigint                     not null,
    type         text default 'other'::text not null,
    content_type text                       not null,
    name         text                       not null,
    size         bigint                     not null,
    status       integer                    not null,
    url          text,
    created_at   timestamp with time zone   not null,
    updated_at   timestamp with time zone,
    constraint files_pk
        primary key (id),
    constraint files_users_id_fk
        foreign key (user_id) references users
            on update set null on delete set null
);

create table if not exists files_to_acts
(
    file_id bigint not null,
    act_id  bigint not null,
    constraint files_to_acts_pk
        unique (file_id, act_id),
    constraint files_to_acts_files_id_fk
        foreign key (file_id) references files,
    constraint files_to_acts_acts_id_fk
        foreign key (act_id) references acts
);


