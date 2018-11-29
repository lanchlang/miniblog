CREATE TABLE public.t_tumblr_blog
(
    id bigint PRIMARY KEY NOT NULL,
    url text DEFAULT ''::text,
    date timestamp DEFAULT now(),
    type text,
    width integer,
    height integer,
    photo_url_1280 text,
    photo_url_500 text,
    photo_url_400 text,
    photo_url_250 text,
    photo_url_100 text,
    reblogged_root_url text,
    photos jsonb,
    video_caption text,
    video_source text,
    video_player text,
    video_player_500 text,
    video_player_250 text,
    blog text
);
CREATE UNIQUE INDEX t_tumblr_blog_id_uindex ON public.t_tumblr_blog (id);