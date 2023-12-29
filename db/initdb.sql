CREATE TABLE book (
    id INT AUTO_INCREMENT NOT NULL,
    iban VARCHAR(255) NOT NULL,
    name varchar(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    pages INT NOT NULL,
    PRIMARY KEY(id),
    UNIQUE INDEX idx_iban (iban)
) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB;
