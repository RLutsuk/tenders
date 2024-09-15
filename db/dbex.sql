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
    tender_id   UUID REFERENCES organization (id) ON DELETE CASCADE,
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

