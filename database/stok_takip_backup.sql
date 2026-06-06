--
-- PostgreSQL database dump
--

\restrict 5wmFbOkjelWyB6El8LcU1ZlwRhvMf8JZVE5SOQNBV3Qzup6RLfUo3AhElqc0YY1

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
-- Name: admin; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin (
    id integer NOT NULL,
    "kullanıcıadı" character varying NOT NULL,
    "şifre" character varying NOT NULL
);


--
-- Name: admin_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.admin_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: admin_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.admin_id_seq OWNED BY public.admin.id;


--
-- Name: adresler; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.adresler (
    id integer NOT NULL,
    musteri_id integer NOT NULL,
    baslik text DEFAULT 'Ev'::text NOT NULL,
    il text NOT NULL,
    ilce text NOT NULL,
    mahalle text NOT NULL,
    adres_satir text NOT NULL,
    posta_kodu text,
    varsayilan boolean DEFAULT false
);


--
-- Name: adresler_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.adresler_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: adresler_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.adresler_id_seq OWNED BY public.adresler.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    isim text NOT NULL
);


--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: musteriler; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.musteriler (
    id integer NOT NULL,
    ad text NOT NULL,
    soyad text NOT NULL,
    email text NOT NULL,
    sifre text NOT NULL,
    tel text,
    olusturulma timestamp with time zone DEFAULT now()
);


--
-- Name: musteriler_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.musteriler_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: musteriler_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.musteriler_id_seq OWNED BY public.musteriler.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: products _id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."products _id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: products _id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."products _id_seq" OWNED BY public.products.id;


--
-- Name: siparis_urunleri; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.siparis_urunleri (
    id integer NOT NULL,
    siparis_id integer NOT NULL,
    urun_id integer,
    urun_isim text NOT NULL,
    adet integer NOT NULL,
    birim_fiyat numeric(10,2) NOT NULL
);


--
-- Name: siparis_urunleri_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.siparis_urunleri_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: siparis_urunleri_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.siparis_urunleri_id_seq OWNED BY public.siparis_urunleri.id;


--
-- Name: siparisler; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.siparisler (
    id integer NOT NULL,
    musteri_id integer,
    adres_id integer,
    toplam_tutar numeric(12,2) NOT NULL,
    durum text DEFAULT 'tamamlandi'::text NOT NULL,
    olusturulma timestamp with time zone DEFAULT now()
);


--
-- Name: siparisler_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.siparisler_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: siparisler_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.siparisler_id_seq OWNED BY public.siparisler.id;


--
-- Name: admin id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin ALTER COLUMN id SET DEFAULT nextval('public.admin_id_seq'::regclass);


--
-- Name: adresler id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.adresler ALTER COLUMN id SET DEFAULT nextval('public.adresler_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: musteriler id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.musteriler ALTER COLUMN id SET DEFAULT nextval('public.musteriler_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public."products _id_seq"'::regclass);


--
-- Name: siparis_urunleri id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparis_urunleri ALTER COLUMN id SET DEFAULT nextval('public.siparis_urunleri_id_seq'::regclass);


--
-- Name: siparisler id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparisler ALTER COLUMN id SET DEFAULT nextval('public.siparisler_id_seq'::regclass);


--
-- Data for Name: admin; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.admin (id, "kullanıcıadı", "şifre") FROM stdin;
1	orkunorkun	2005
\.


--
-- Data for Name: adresler; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.adresler (id, musteri_id, baslik, il, ilce, mahalle, adres_satir, posta_kodu, varsayilan) FROM stdin;
1	1	kayabası mah 75.yıl caddesi	istanbul	başakşehir	kayabası	b 27 d2	34900	t
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.categories (id, isim) FROM stdin;
2	Telefon
3	Bilgisayar
\.


--
-- Data for Name: musteriler; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.musteriler (id, ad, soyad, email, sifre, tel, olusturulma) FROM stdin;
1	orkun	teke	orkunhulya34@gmail.com	d614bf06cbefa888cc9f4743c478abb5798cfe60a562a21615bb70ba7258ddd7	5383814957	2026-06-06 13:52:59.360452+03
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.products (id, isim, fiyat, stok, kategori, resim, "tanım") FROM stdin;
2	rtx 2090	2500	2	Bilgisayar	https://encrypted-tbn3.gstatic.com/shopping?q=tbn:ANd9GcQV1ehdC2ZRpZ9D3I89rtgZdjBbME7VDS4SfDIlxw7frwaqKX8z2oPGOMjSVvJv6_G7NQxF-dfzoA6J6H60ddttBYR_v-Bojg	oyun canavarı
1	rtx 5090	55000	10	Telefon	https://static.sinerji.gen.tr/Images/MD/L9Q-39f4b516a0a7e5344fe9a30b8cb9721a.jpg	oyun canavarı
\.


--
-- Data for Name: siparis_urunleri; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.siparis_urunleri (id, siparis_id, urun_id, urun_isim, adet, birim_fiyat) FROM stdin;
1	1	2	rtx 2090	1	2500.00
2	2	2	rtx 2090	1	2500.00
\.


--
-- Data for Name: siparisler; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.siparisler (id, musteri_id, adres_id, toplam_tutar, durum, olusturulma) FROM stdin;
1	1	1	2500.00	tamamlandi	2026-06-06 13:53:40.259458+03
2	1	1	2500.00	tamamlandi	2026-06-06 14:10:23.483924+03
\.


--
-- Name: admin_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.admin_id_seq', 1, true);


--
-- Name: adresler_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.adresler_id_seq', 1, true);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.categories_id_seq', 3, true);


--
-- Name: musteriler_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.musteriler_id_seq', 1, true);


--
-- Name: products _id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."products _id_seq"', 2, true);


--
-- Name: siparis_urunleri_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.siparis_urunleri_id_seq', 2, true);


--
-- Name: siparisler_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.siparisler_id_seq', 2, true);


--
-- Name: admin admin_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (id);


--
-- Name: adresler adresler_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.adresler
    ADD CONSTRAINT adresler_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: musteriler musteriler_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.musteriler
    ADD CONSTRAINT musteriler_email_key UNIQUE (email);


--
-- Name: musteriler musteriler_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.musteriler
    ADD CONSTRAINT musteriler_pkey PRIMARY KEY (id);


--
-- Name: products products _pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT "products _pkey" PRIMARY KEY (id);


--
-- Name: siparis_urunleri siparis_urunleri_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparis_urunleri
    ADD CONSTRAINT siparis_urunleri_pkey PRIMARY KEY (id);


--
-- Name: siparisler siparisler_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparisler
    ADD CONSTRAINT siparisler_pkey PRIMARY KEY (id);


--
-- Name: adresler adresler_musteri_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.adresler
    ADD CONSTRAINT adresler_musteri_id_fkey FOREIGN KEY (musteri_id) REFERENCES public.musteriler(id) ON DELETE CASCADE;


--
-- Name: siparis_urunleri siparis_urunleri_siparis_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparis_urunleri
    ADD CONSTRAINT siparis_urunleri_siparis_id_fkey FOREIGN KEY (siparis_id) REFERENCES public.siparisler(id) ON DELETE CASCADE;


--
-- Name: siparis_urunleri siparis_urunleri_urun_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparis_urunleri
    ADD CONSTRAINT siparis_urunleri_urun_id_fkey FOREIGN KEY (urun_id) REFERENCES public.products(id) ON DELETE SET NULL;


--
-- Name: siparisler siparisler_adres_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparisler
    ADD CONSTRAINT siparisler_adres_id_fkey FOREIGN KEY (adres_id) REFERENCES public.adresler(id) ON DELETE SET NULL;


--
-- Name: siparisler siparisler_musteri_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.siparisler
    ADD CONSTRAINT siparisler_musteri_id_fkey FOREIGN KEY (musteri_id) REFERENCES public.musteriler(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

\unrestrict 5wmFbOkjelWyB6El8LcU1ZlwRhvMf8JZVE5SOQNBV3Qzup6RLfUo3AhElqc0YY1

