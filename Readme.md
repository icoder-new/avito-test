# Сервис проведения тендеров
## Проблема:

Авито — большая компания, в рамках которой пользователи не только продают/покупают товары и услуги, но и предоставляют помощь крупному бизнесу и предприятиям.

Поэтому ребята из Авито решили сделать сервис, который позволит бизнесу создать тендер на оказание каdких-либо услуг. А пользователи/другие бизнесы будут предлагать свои выгодные условия для получения данного тендера.

Помогите ребятам из Авито реализовать новое HTTP API!

## Про приложение
- Golang
- github.com/gin-gonic/gin - основной HTTP-фреймворк
- github.com/go-playground/validator/v10 - для валидации данных 
- github.com/google/uuid -  для генерации UUID
- github.com/ilyakaznacheev/cleanenv - для удобного чтения переменных окружения
- github.com/joho/godotenv - для загрузки переменных окружения из .env файла 
- github.com/lib/pq - драйвер для работы с PostgreSQL
- golang.org/x/exp/slog - для логирования.
- Postgres

### Сущности

```sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    service_type VARCHAR(20) CHECK (service_type IN ('Construction', 'Delivery', 'Manufacture')) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('Created', 'Published', 'Closed')) NOT NULL,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    responsible_id UUID REFERENCES employee(id),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tenders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    service_type VARCHAR(20) CHECK (service_type IN ('Construction', 'Delivery', 'Manufacture')) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('Created', 'Published', 'Closed')) NOT NULL,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    responsible_id UUID REFERENCES employee(id),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50),
    tender_id UUID NOT NULL,
    author_type VARCHAR(50),
    version INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    decision TEXT,
    author_id UUID NOT NULL,
    FOREIGN KEY (tender_id) REFERENCES tenders(id)
);

CREATE TABLE bid_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    author_type VARCHAR(50) NOT NULL,
    decision VARCHAR(255),
    version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    author_id UUID NOT NULL,
    FOREIGN KEY (bid_id) REFERENCES bids(id)
);

CREATE TABLE bid_feedbacks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    feedback TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    author_id UUID NOT NULL,
    FOREIGN KEY (bid_id) REFERENCES bids(id)
);
```

### Команды для сборки и запуска контейнера:

1. Сборка Docker-образа
```bash
$ docker build -t tender-api .
```
2. Запуск контейнера:
```bash
$ docker run -d -p 8080:8080 --name tender-api tender-api
```
### Зачем нужен multi-stage build?

Multi-stage build уменьшает размер итогового Docker-образа. В первом этапе контейнер собирает бинарник с помощью полной Go-среды, а во втором этапе запускает приложение в минимальной среде, такой как Alpine Linux, которая весит значительно меньше.