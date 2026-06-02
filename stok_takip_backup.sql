--
-- PostgreSQL database dump
--

\restrict TlqlZJ9U6xqRRBvSiBpfqbNOcslmJ7dYqOr4jRU3jGyfbCpf5d48wIHb5nE4N6A

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admin; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admin (
    id integer NOT NULL,
    "kullanıcıadı" character varying NOT NULL,
    "şifre" character varying NOT NULL
);


ALTER TABLE public.admin OWNER TO postgres;

--
-- Name: admin_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.admin_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.admin_id_seq OWNER TO postgres;

--
-- Name: admin_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.admin_id_seq OWNED BY public.admin.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    isim text NOT NULL
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer CONSTRAINT "products _id_not_null" NOT NULL,
    isim character varying(255) CONSTRAINT "products _isim_not_null" NOT NULL,
    fiyat numeric(10,0) CONSTRAINT "products _fiyat_not_null" NOT NULL,
    stok integer CONSTRAINT "products _stok_not_null" NOT NULL,
    kategori character varying(100) CONSTRAINT "products _kategori_not_null" NOT NULL,
    resim text,
    "tanım" text
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products _id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."products _id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."products _id_seq" OWNER TO postgres;

--
-- Name: products _id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."products _id_seq" OWNED BY public.products.id;


--
-- Name: admin id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin ALTER COLUMN id SET DEFAULT nextval('public.admin_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public."products _id_seq"'::regclass);


--
-- Data for Name: admin; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.admin (id, "kullanıcıadı", "şifre") FROM stdin;
1	orkunorkun	2005
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, isim) FROM stdin;
1	Telefon
2	Telefon
3	Bilgisayar
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, isim, fiyat, stok, kategori, resim, "tanım") FROM stdin;
1	rtx 5090 	55000	10	Bilgisayar	https://static.sinerji.gen.tr/Images/MD/L9Q-39f4b516a0a7e5344fe9a30b8cb9721a.jpg	oyun canavarı
2	rtx 2090	2500	4	Bilgisayar	https://encrypted-tbn3.gstatic.com/shopping?q=tbn:ANd9GcQV1ehdC2ZRpZ9D3I89rtgZdjBbME7VDS4SfDIlxw7frwaqKX8z2oPGOMjSVvJv6_G7NQxF-dfzoA6J6H60ddttBYR_v-Bojg	oyun canavarı
\.


--
-- Name: admin_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.admin_id_seq', 1, true);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 3, true);


--
-- Name: products _id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."products _id_seq"', 2, true);


--
-- Name: admin admin_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: products products _pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT "products _pkey" PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict TlqlZJ9U6xqRRBvSiBpfqbNOcslmJ7dYqOr4jRU3jGyfbCpf5d48wIHb5nE4N6A

