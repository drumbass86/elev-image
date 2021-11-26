DROP TABLE IF EXISTS captured_images;
CREATE TABLE IF NOT EXISTS public.captured_images
(
    id SERIAL NOT NULL,
    "time" INTEGER NOT NULL UNIQUE,
    coordinate point NOT NULL,    
    countangles smallint NOT NULL,
    elevation_angles real ARRAY NOT NULL, 
    path_raw varchar(512),
    path_image varchar(512),
    CONSTRAINT captured_images_pkey PRIMARY KEY (id, "time"),
    CONSTRAINT captured_images_time_check check ("time" > '0'::integer)
)

TABLESPACE pg_default;

ALTER TABLE public.captured_images
    OWNER to "user";

COMMENT ON COLUMN public.captured_images.coordinate
    IS 'GPS coordinate captured image';

-- Index: product_idx

DROP INDEX IF EXISTS public.captured_images_idx;
CREATE INDEX captured_images_idx
    ON public.captured_images USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;