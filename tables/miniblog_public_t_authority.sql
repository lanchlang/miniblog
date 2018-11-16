CREATE TABLE public.t_authority
(
    id integer DEFAULT nextval('t_authority_id_seq'::regclass) PRIMARY KEY NOT NULL,
    name text NOT NULL,
    department_id integer,
    CONSTRAINT t_authority_t_departments_id_fk FOREIGN KEY (department_id) REFERENCES public.t_departments (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX t_authority_name_uindex ON public.t_authority (name);
INSERT INTO public.t_authority (id, name, department_id) VALUES (1, 'ddd', null);