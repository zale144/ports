create table if not exists port (
    id text unique not null,
    name text,
    city text,
    country text,
    alias json,
    regions json,
    coordinates json,
    province text,
    timezone text,
    unlocs json,
    code text
)
