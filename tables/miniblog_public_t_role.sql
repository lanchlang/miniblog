CREATE TABLE public.t_role
(
    id integer DEFAULT nextval('t_role_id_seq'::regclass) PRIMARY KEY NOT NULL,
    name text,
    auths integer[]
);
CREATE UNIQUE INDEX t_role_name_uindex ON public.t_role (name);
INSERT INTO public.t_role (id, name, auths) VALUES (1, 'aa', '{12,3}');