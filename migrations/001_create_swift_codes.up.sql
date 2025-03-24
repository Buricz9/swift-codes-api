CREATE SCHEMA IF NOT EXISTS swift;

CREATE TABLE swift.swift_codes (
                                   id SERIAL PRIMARY KEY,
                                   swift_code VARCHAR(11) NOT NULL UNIQUE,
                                   bank_name VARCHAR(255) NOT NULL,
                                   address TEXT NOT NULL,
                                   country_iso2 CHAR(2) NOT NULL,
                                   country_name VARCHAR(100) NOT NULL,
                                   is_headquarter BOOLEAN NOT NULL,
                                   headquarter_swift_code VARCHAR(11),
                                   CONSTRAINT fk_headquarter
                                       FOREIGN KEY (headquarter_swift_code)
                                           REFERENCES swift.swift_codes(swift_code)
                                           ON DELETE SET NULL
);

CREATE INDEX idx_swift_codes_country_iso2
    ON swift.swift_codes (country_iso2);

CREATE INDEX idx_swift_codes_swift_code
    ON swift.swift_codes (swift_code);
