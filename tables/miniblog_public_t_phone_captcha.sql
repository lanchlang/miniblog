CREATE TABLE public.t_phone_captcha
(
    id bigint DEFAULT nextval('t_phone_captcha_id_seq'::regclass) PRIMARY KEY NOT NULL,
    phone text NOT NULL,
    verify_code text NOT NULL,
    expires timestamp
);