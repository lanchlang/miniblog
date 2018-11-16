CREATE TABLE public.t_blog
(
    id bigint DEFAULT nextval('t_blog_id_seq'::regclass) NOT NULL,
    u_id bigint NOT NULL,
    u_name text NOT NULL,
    title text DEFAULT 'empty title'::text NOT NULL,
    date_create timestamp DEFAULT now(),
    comment_cnt integer DEFAULT 0,
    like_cnt integer DEFAULT 0,
    intro text DEFAULT '暂无简介'::text,
    c_name text NOT NULL,
    c_id integer NOT NULL,
    tags text[],
    referer text DEFAULT ''::text,
    cover jsonb,
    content jsonb,
    date_last_update timestamp DEFAULT now(),
    type smallint DEFAULT 1,
    access_limit integer DEFAULT 0,
    view_cnt integer DEFAULT 0,
    deleted smallint DEFAULT 0
);