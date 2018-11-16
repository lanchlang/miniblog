CREATE TABLE public.t_comment
(
    id bigint PRIMARY KEY NOT NULL,
    u_id bigint NOT NULL,
    u_name text NOT NULL,
    b_id bigint NOT NULL,
    content text NOT NULL,
    date timestamp DEFAULT now(),
    type smallint DEFAULT 0,
    reply_u_id bigint,
    reply_u_name text,
    deleted smallint DEFAULT 0
);