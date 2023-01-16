create table app_feedback
(
    id           varchar(36)                            not null
        primary key,
    feedback_id  varchar(36)                            not null,
    store_id     int                                    not null,
    app_id       varchar(64)  default ''                not null,
    version      int          default 1                 not null,
    `key`        varchar(50)                            not null,
    value        varchar(500) default ''                not null,
    key_type     int          default 0                 not null,
    trigger_type int          default 1                 not null,
    from_type    int          default 1                 not null,
    created_at   datetime     default CURRENT_TIMESTAMP not null,
    updated_at   datetime     default CURRENT_TIMESTAMP not null
);

create index index_app_feedback_on_app_id_idx
    on app_feedback (app_id);

create index index_app_feedback_on_store_id_idx
    on app_feedback (store_id);

create table app_use_history
(
    id            varchar(36)                           not null
        primary key,
    store_id      int                                   not null,
    app_id        varchar(64) default ''                not null,
    version       int         default 1                 not null,
    count         int         default 0                 not null,
    last_pop_time datetime    default CURRENT_TIMESTAMP not null,
    created_at    datetime    default CURRENT_TIMESTAMP not null,
    updated_at    datetime    default CURRENT_TIMESTAMP not null
);

create index index_app_feedback_on_app_id_idx
    on app_use_history (app_id);

create index index_app_feedback_on_store_id_idx
    on app_use_history (store_id);

create table app_version
(
    id           varchar(36)                           not null
        primary key,
    app_id       varchar(64)                           not null,
    app_name     varchar(50)                           not null,
    version      int         default 1                 not null,
    version_name varchar(50) default ''                not null,
    type         int         default 0                 not null,
    created_at   datetime    default CURRENT_TIMESTAMP not null,
    updated_at   datetime    default CURRENT_TIMESTAMP not null,
    constraint uni_index_app_version_on_app_id_version_idx
        unique (app_id, version)
);

create table application_charges
(
    id               bigint auto_increment comment '账单id'
        primary key,
    store_id         varchar(36)  null,
    application_id   varchar(36)  null,
    status           varchar(50)  null,
    name             varchar(255) null,
    return_url       varchar(255) null,
    test             tinyint(1)   null,
    created_at       datetime     not null,
    updated_at       datetime     null,
    price            bigint       null,
    confirmation_url varchar(255) null
);

create table application_locales
(
    id             varbinary(16) not null
        primary key,
    name           varchar(255)  not null,
    subtitle       text          null,
    application_id bigint        not null,
    locale         varchar(255)  not null,
    is_primary         tinyint(1)  null,
    icon         varchar(255)  null,
    created_at       datetime     not null,
    updated_at       datetime     null,
    `desc`         mediumtext    null,
    constraint index_application_locales_on_application_id_and_locale
        unique (application_id, locale)
);

create index index_application_locales_on_application_id
    on application_locales (application_id);

create table application_plans
(
    id                       bigint auto_increment comment '应用套餐id'
        primary key,
    application_id           bigint       null comment '应用插件id',
    name                     varchar(255) null comment '套餐名',
    price                    bigint       null comment '套餐价格',
    subscription_type        varchar(50)  null comment 'monthly、yearly',
    trial_days               int(10)      null comment '试用天数',
    usage_charge_description varchar(255) null comment 'additional usage charge的描述',
    plan_type                varchar(50)  null comment 'free、recurring、one-time',
    created_at               datetime     not null,
    updated_at               datetime     null
);

create table application_platforms
(
    id             varbinary(16) not null
        primary key,
    application_id bigint        not null,
    platform       varchar(255)  not null,
    listed        tinyint(1)   DEFAULT '0' COMMENT '指明是否在应用市场显示',
    unlisted_reason varchar(512) DEFAULT '' COMMENT '说明unlisted原因',
    sended_email tinyint(1) DEFAULT '0' COMMENT '指明是否需要发送mail给开发者',
    constraint index_application_platforms_on_application_id_and_platform
        unique (application_id, platform)
);

create index index_application_platforms_on_application_id
    on application_platforms (application_id);

create table application_subscriptions
(
    id             int auto_increment
        primary key,
    store_id       bigint      null,
    application_id bigint      null,
    plan_id        int         null,
    expired_at     datetime    null,
    created_at     datetime    not null,
    updated_at     datetime    not null,
    status         varchar(20) null
);

create table ar_internal_metadata
(
    `key`      varchar(255) charset utf8 not null
        primary key,
    value      varchar(255)              null,
    created_at datetime                  not null,
    updated_at datetime                  not null
);

create table oauth_applications
(
    id                  bigint auto_increment
        primary key,
    name                varchar(255)             not null,
    uid                 varchar(255)             not null,
    secret              varchar(255)             not null,
    redirect_uri        text                     not null,
    scopes              varchar(1024) default '' not null,
    confidential        tinyint(1)    default 1  not null,
    created_at          datetime                 not null,
    updated_at          datetime                 not null,
    icon                text                     null comment '插件图标',
    embbed              tinyint(1)    default 0  null,
    private_app         tinyint(1)    default 0  null,
    store_id            varbinary(16)            null,
    email               varchar(255)             null,
    category            varchar(255)             null,
    oauth_dancable      tinyint(1)    default 1  null,
    subscribable        tinyint(1)    default 0  null,
    store_charge        tinyint(1)    default 0  null,
    notify_uri          text                     null,
    link                text                     null,
    app_uri             text                     null,
    charge_min_plan     varchar(255)             null,
    free_stores         text                     null,
    position            int                      null,
    webhook_api_version varchar(255)             null,
    tags                varchar(1024)            null comment 'add tags for recommended',
    constraint index_oauth_applications_on_uid
        unique (uid)
);

create table oauth_access_grants
(
    id                bigint auto_increment
        primary key,
    resource_owner_id varbinary(16)            not null,
    application_id    bigint                   not null,
    token             varchar(255)             not null,
    expires_in        int                      not null,
    redirect_uri      text                     not null,
    created_at        datetime                 not null,
    revoked_at        datetime                 null,
    scopes            varchar(1024) default '' not null,
    constraint index_oauth_access_grants_on_token
        unique (token)
);

create index index_oauth_access_grants_on_application_id
    on oauth_access_grants (application_id);

create index index_oauth_access_grants_on_resource_owner_id
    on oauth_access_grants (resource_owner_id);

create table oauth_access_tokens
(
    id                     bigint auto_increment
        primary key,
    resource_owner_id      varbinary(16)            null,
    application_id         bigint                   null,
    token                  varchar(255)             not null,
    refresh_token          varchar(255)             null,
    expires_in             int                      null,
    revoked_at             datetime                 null,
    created_at             datetime                 not null,
    scopes                 varchar(1024) default '' not null,
    previous_refresh_token varchar(255)  default '' not null,
    constraint index_oauth_access_tokens_on_refresh_token
        unique (refresh_token),
    constraint index_oauth_access_tokens_on_token
        unique (token)
);

create index idx_oauth_access_tokens_res_app_rvk
    on oauth_access_tokens (resource_owner_id, application_id, revoked_at);

create index index_oauth_access_tokens_on_application_id
    on oauth_access_tokens (application_id);

create table plans
(
    id                varbinary(16)               not null
        primary key,
    application_id    bigint                      null,
    plan_type         varchar(255)                null,
    subscription_type varchar(255)                null,
    per_price         decimal(15, 2) default 0.00 null,
    price             decimal(15, 2) default 0.00 null,
    ori_price         decimal(15, 2) default 0.00 null,
    ori_per_price     decimal(15, 2) default 0.00 null,
    month             int            default 0    null,
    day               int            default 0    null,
    enabled           tinyint(1)     default 1    null,
    created_at        datetime                    not null,
    updated_at        datetime                    not null
);

create index index_plans_on_application_id
    on plans (application_id);

create table recurring_application_charges
(
    id               bigint auto_increment
        primary key,
    store_id         varchar(36)  null,
    application_id   bigint       null,
    status           varchar(20)  null,
    activated_on     datetime     null,
    billing_on       datetime     null,
    trial_end_on     datetime     null,
    name             varchar(255) null,
    capped_amount    int(255)     null,
    terms            varchar(255) null,
    cancelled_on     datetime     null,
    confirmation_url varchar(255) null,
    return_url       varchar(255) null,
    created_at       datetime     not null,
    updated_at       datetime     null,
    price            int(255)     null,
    test             tinyint(1)   null,
    trial_days       tinyint(10)  null
);

create table stores
(
    id         varbinary(16) not null
        primary key,
    origin_id  varchar(255)  null,
    name       varchar(255)  null comment '店铺名称',
    created_at datetime      not null,
    updated_at datetime      not null,
    platform   varchar(255)  null,
    constraint index_stores_on_origin_id
        unique (origin_id)
);

create table exclusive_stores
(
    id             varbinary(16) not null
        primary key,
    application_id bigint        not null,
    store_id       varbinary(16) null,
    created_at     datetime(6)   not null,
    updated_at     datetime(6)   not null,
    constraint index_exclusive_stores_on_store_id_and_application_id
        unique (store_id, application_id)
);

create index index_exclusive_stores_on_store_id
    on exclusive_stores (store_id);

create table install_tracks
(
    id             varbinary(16) not null
        primary key,
    application_id bigint        null,
    store_id       varbinary(16) null,
    installed_at   datetime      null,
    created_at     datetime      not null,
    updated_at     datetime      not null,
    constraint index_install_tracks_on_store_id_and_application_id
        unique (store_id, application_id)
);

create index index_install_tracks_on_store_id
    on install_tracks (store_id);

create table subscriptions
(
    id               varbinary(16)        not null
        primary key,
    application_id   bigint               null,
    store_id         varbinary(16)        null,
    plan_id          varbinary(16)        null,
    expired          tinyint(1) default 0 null,
    expired_at       datetime             null,
    created_at       datetime             not null,
    updated_at       datetime             not null,
    five_days_notify tinyint(1) default 0 null,
    two_days_notify  tinyint(1) default 0 null
);

create index index_subscriptions_on_application_id
    on subscriptions (application_id);

create index index_subscriptions_on_five_days_notify
    on subscriptions (five_days_notify);

create index index_subscriptions_on_plan_id
    on subscriptions (plan_id);

create index index_subscriptions_on_store_id
    on subscriptions (store_id);

create index index_subscriptions_on_two_days_notify
    on subscriptions (two_days_notify);

create table transactions
(
    id                varbinary(16)               not null
        primary key,
    application_id    bigint                      null,
    store_id          varbinary(16)               null,
    plan_id           varbinary(16)               null,
    order_number      varchar(255)                null,
    price             decimal(15, 2) default 0.00 null,
    before_expired_at datetime                    null,
    after_expired_at  datetime                    null,
    payment_type      varchar(255)                null,
    pre_plan_id       varbinary(16)               null,
    created_at        datetime                    not null,
    updated_at        datetime                    not null
);

create index index_transactions_on_application_id
    on transactions (application_id);

create index index_transactions_on_plan_id
    on transactions (plan_id);

create index index_transactions_on_store_id
    on transactions (store_id);

create table app_tags
(
    id                bigint               not null
        primary key,
    application_id    bigint                      null,
    tag_id    bigint                      null,
    created_at        datetime                    not null,
    updated_at        datetime                    not null
);

create table tags
(
    id                bigint               not null
        primary key,
    tag_level    varchar(255)                      not null,
    parent_id    bigint                      not null,
     name_zh    varchar(255)                       null,
      name_en    varchar(255)                      null,
    created_at        datetime                    not null,
    updated_at        datetime                    not null
);

create table dev_oauth_applications
(
    id                bigint       auto_increment
        primary key,
    name    varchar(255)                      not null,
    redirect_uri        text                     not null,
    icon                text                     null comment '插件图标',
    email               varchar(255)             null,
    app_uri             text                     null,
    webhook_api_version varchar(255)             null,
    additional_info varchar(255)             null,
    created_at        datetime                    not null,
    updated_at        datetime                    not null
);

create table dev_application_locales
(
    id             bigint auto_increment
        primary key,
    name           varchar(255)  not null,
    subtitle       text          null,
    application_id bigint        not null,
    locale         varchar(255)  not null,
    is_primary         tinyint(1)  null,
    icon         varchar(255)  null,
    created_at       datetime     not null,
    updated_at       datetime     null,
    `desc`         mediumtext    null,
    constraint index_application_locales_on_application_id_and_locale
        unique (application_id, locale)
);

create index index_application_locales_on_application_id
    on dev_application_locales (application_id);
