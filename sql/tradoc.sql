CREATE TABLE traductors (
  id SMALLSERIAL PRIMARY KEY, -- 32 767 traducteurs maxi
  pseudonim VARCHAR(25) NOT NULL,
  senhal TEXT,
  code_activation INTEGER,
  suspension BOOLEAN DEFAULT FALSE NOT NULL,
  creat_lo TIMESTAMP,
  refresh_token CHAR(32) CHECK(length(refresh_token) = 32),
  CONSTRAINT chk_traductors_code CHECK (code_activation > 1000 AND code_activation < 9999)
);
CREATE UNIQUE INDEX traductors_pseudonim_unique_idx ON traductors (pseudonim) WITH (deduplicate_items = off);


CREATE TABLE textes_oc (
  id SERIAL PRIMARY KEY,
  frasa TEXT NOT NULL,
  creat_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX textes_oc_frasa_unique_idx ON textes_oc (frasa) WITH (deduplicate_items = off);


CREATE TABLE translation_files (
  id SMALLSERIAL PRIMARY KEY,
  dialect_name TEXT NOT NULL,
  filename_fr TEXT NOT NULL,
  filename_en TEXT,
  creat_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- l'ordre des dialectes doit être respecter
CREATE TABLE dialectes_auth (
  -- dialectes standard
  auvernhat_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  gascon_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  lengadocian_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  lemosin_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  provencau_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  vivaroaupenc_estandard BOOLEAN DEFAULT FALSE NOT NULL,
  -- dialectes locaux
  auvernhat_brivades BOOLEAN DEFAULT FALSE NOT NULL,
  auvernhat_septentrional BOOLEAN DEFAULT FALSE NOT NULL,
  gascon_aranes BOOLEAN DEFAULT FALSE NOT NULL,
  gascon_bearnes BOOLEAN DEFAULT FALSE NOT NULL,
  lengadocian_agenes BOOLEAN DEFAULT FALSE NOT NULL,
  lengadocian_besierenc BOOLEAN DEFAULT FALSE NOT NULL,
  lengadocian_carcasses BOOLEAN DEFAULT FALSE NOT NULL,
  lengadocian_roergat BOOLEAN DEFAULT FALSE NOT NULL,
  lemosin_marches BOOLEAN DEFAULT FALSE NOT NULL,
  lemosin_peiregordin BOOLEAN DEFAULT FALSE NOT NULL,
  provencau_maritime BOOLEAN DEFAULT FALSE NOT NULL,
  provencau_nicard BOOLEAN DEFAULT FALSE NOT NULL,
  provencau_rodanenc BOOLEAN DEFAULT FALSE NOT NULL,
  vivaroaupenc_aupenc BOOLEAN DEFAULT FALSE NOT NULL,
  vivaroaupenc_gavot BOOLEAN DEFAULT FALSE NOT NULL,
  vivaroaupenc_vivarodaufinenc BOOLEAN DEFAULT FALSE NOT NULL,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0)
);
CREATE UNIQUE INDEX dialectes_auth_traductor_id_unique_idx ON dialectes_auth (traductor_id) WITH (deduplicate_items = off);

------------------------
-- Dialectes Standard --
------------------------

CREATE TABLE auvernhat_estandard (
  id SERIAL PRIMARY KEY,
  frasa_auv_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE gascon_estandard (
  id SERIAL PRIMARY KEY,
  frasa_gas_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lengadocian_estandard (
  id SERIAL PRIMARY KEY,
  frasa_len_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lemosin_estandard (
  id SERIAL PRIMARY KEY,
  frasa_lem_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE provencau_estandard (
  id SERIAL PRIMARY KEY,
  frasa_pro_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE vivaroaupenc_estandard (
  id SERIAL PRIMARY KEY,
  frasa_viv_est TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-------------------------
-- Dialectes Auvergnat --
-------------------------

CREATE TABLE auvernhat_brivades (
  id SERIAL PRIMARY KEY,
  frasa_auv_bri TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE auvernhat_septentrional (
  id SERIAL PRIMARY KEY,
  frasa_auv_sep TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

----------------------
-- Dialectes Gascon --
----------------------

CREATE TABLE gascon_aranes (
  id SERIAL PRIMARY KEY,
  frasa_gas_ara TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE gascon_bearnes (
  id SERIAL PRIMARY KEY,
  frasa_gas_bea TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

----------------------------
-- Dialectes Languedocien --
----------------------------

CREATE TABLE lengadocian_agenes (
  id SERIAL PRIMARY KEY,
  frasa_len_age TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lengadocian_besierenc (
  id SERIAL PRIMARY KEY,
  frasa_len_bes TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lengadocian_carcasses (
  id SERIAL PRIMARY KEY,
  frasa_len_car TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lengadocian_roergat (
  id SERIAL PRIMARY KEY,
  frasa_len_roe TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

------------------------
-- Dialectes Limousin --
------------------------

CREATE TABLE lemosin_marches (
  id SERIAL PRIMARY KEY,
  frasa_lem_mar TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lemosin_peiregordin (
  id SERIAL PRIMARY KEY,
  frasa_lem_pei TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-------------------------
-- Dialectes Provençal --
-------------------------

CREATE TABLE provencau_maritime (
  id SERIAL PRIMARY KEY,
  frasa_pro_mar TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE provencau_nicard (
  id SERIAL PRIMARY KEY,
  frasa_pro_nic TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE provencau_rodanenc (
  id SERIAL PRIMARY KEY,
  frasa_pro_rod TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

----------------------------
-- Dialectes Vivaroaupenc --
----------------------------

CREATE TABLE vivaroaupenc_aupenc (
  id SERIAL PRIMARY KEY,
  frasa_viv_aup TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE vivaroaupenc_gavot (
  id SERIAL PRIMARY KEY,
  frasa_viv_gav TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE vivaroaupenc_vivarodaufinenc (
  id SERIAL PRIMARY KEY,
  frasa_viv_viv TEXT NOT NULL,
  frasa_fr TEXT NOT NULL,
  frasa_an TEXT,
  traductor_id SMALLINT NOT NULL CHECK(traductor_id > 0),
  texte_oc_id INTEGER NOT NULL CHECK(texte_oc_id > 0),
  tradusit_lo TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
