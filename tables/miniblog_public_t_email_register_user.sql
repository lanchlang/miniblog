CREATE TABLE public.t_email_register_user
(
    id bigint DEFAULT nextval('t_email_register_user_id_seq'::regclass) PRIMARY KEY NOT NULL,
    username text,
    email text,
    password text,
    expires timestamp,
    verify_code text
);
CREATE UNIQUE INDEX t_email_register_user_verifycode_uindex ON public.t_email_register_user (verify_code);