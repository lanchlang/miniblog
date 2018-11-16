CREATE TABLE public.t_users
(
    id bigint DEFAULT nextval('t_users_id_seq'::regclass) PRIMARY KEY NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    phone text,
    email text,
    roles integer[],
    date_create timestamp DEFAULT now(),
    last_login timestamp DEFAULT now(),
    likes bigint[] DEFAULT '{}'::bigint[],
    access_level integer DEFAULT 0
);
CREATE UNIQUE INDEX t_users_username_uindex ON public.t_users (username);
CREATE UNIQUE INDEX t_users_phone_uindex ON public.t_users (phone);
CREATE UNIQUE INDEX t_users_email_uindex ON public.t_users (email);
INSERT INTO public.t_users (id, username, password, phone, email, roles, date_create, last_login, likes, access_level) VALUES (1, 'admin', 'admin', null, null, null, '2018-11-13 03:49:38.076465', '2018-11-13 03:49:38.085451', null, 0);