CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);


CREATE TABLE organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);

CREATE TYPE status_tender AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE tender
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(100),
    status          status_tender,
    service_type    VARCHAR(100),
    version         INT              DEFAULT 0,
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_name       VARCHAR(100) REFERENCES employee (username),
    created_at      TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE status_bid AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

CREATE TABLE bid
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(100),
    status      status_bid,
    tender_id   UUID REFERENCES tender (id) ON DELETE CASCADE,
    author_type VARCHAR(100),
    user_id     UUID REFERENCES employee (id) ON DELETE CASCADE,
    version     INT              DEFAULT 0,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);


-- users
CREATE INDEX IF NOT EXISTS index_employee_username ON employee (username);
CREATE INDEX IF NOT EXISTS index_organization_responsible_org_id ON organization_responsible (organization_id);
CREATE INDEX IF NOT EXISTS index_organization_responsible_user_id ON organization_responsible (user_id);
CREATE INDEX IF NOT EXISTS index_tender_id ON tender (id);
CREATE INDEX IF NOT EXISTS index_tender_status ON tender (status, name);
CREATE INDEX IF NOT EXISTS index_tender_status_service_type ON tender (service_type, name) WHERE status = 'Published';
CREATE INDEX IF NOT EXISTS index_bid_id ON bid (id);

--TEST
INSERT INTO employee (id, username, first_name, last_name)
VALUES ('550e8400-e29b-41d4-a716-446655440000', 'john_doe', 'John', 'Doe');

INSERT INTO employee (id, username, first_name, last_name)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'jane_doe', 'Jane', 'Doe');

INSERT INTO employee (id, username, first_name, last_name)
VALUES ('0dcc0451-b4b3-4682-bb6c-fbbc12e290f1', 'bob_dilan', 'Bob', 'Dilan');

INSERT INTO employee (id, username, first_name, last_name)
VALUES ('c99300ac-34af-49b2-b5a2-1eb279b1b6e2', 'ann_jr', 'Ann', 'Junior');

INSERT INTO organization (id, name, description, type)
VALUES ('a6b3f3ff-d1db-4703-a104-be76b256dc59', 'ABC Corp', 'A building company', 'LLC');

INSERT INTO organization (id, name, description, type)
VALUES ('987e6543-b21c-34d5-a789-123456789abc', 'Tech Corp', 'A technology company', 'LLC');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES ('987e6543-b21c-34d5-a789-123456789abc', '550e8400-e29b-41d4-a716-446655440000');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES ('987e6543-b21c-34d5-a789-123456789abc', '0dcc0451-b4b3-4682-bb6c-fbbc12e290f1');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES ('a6b3f3ff-d1db-4703-a104-be76b256dc59', 'c99300ac-34af-49b2-b5a2-1eb279b1b6e2');

--Exist tenders
INSERT INTO tender (id, name, description, status, service_type, version, organization_id, user_name)
VALUES ('4183447b-1bd2-4d8c-82fc-1006790df7de', 'Delivery box from Msc to SPb', 'Delivery box with discs', 
'Created', 'Delivery', 1, '987e6543-b21c-34d5-a789-123456789abc', 'john_doe');

INSERT INTO tender (id, name, description, status, service_type, version, organization_id, user_name)
VALUES ('20b8fe6e-1c6e-4aa6-9a27-a9b77cfe53a2', 'Office cleaning', 'office cleaning in Moscow', 
'Published', 'Manufacture', 1, '987e6543-b21c-34d5-a789-123456789abc', 'john_doe');

INSERT INTO tender (id, name, description, status, service_type, version, organization_id, user_name)
VALUES ('ec0783e7-5b52-4130-ac88-9d3623d9a47d', 'Pizza delivery', 'pizza delivery for corporate events', 
'Published', 'Delivery', 1, '987e6543-b21c-34d5-a789-123456789abc', 'bob_dilan');



