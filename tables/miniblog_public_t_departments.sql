CREATE TABLE public.t_departments
(
    id integer DEFAULT nextval('t_departments_id_seq'::regclass) PRIMARY KEY NOT NULL,
    name text NOT NULL
);
CREATE UNIQUE INDEX t_departments_name_uindex ON public.t_departments (name);
INSERT INTO public.t_departments (id, name) VALUES (2, 'aa');