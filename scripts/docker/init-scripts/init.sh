set -e

until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  sleep 2
done

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

    CREATE TABLE subscriptions (
        id SERIAL PRIMARY KEY,
        service_name VARCHAR(255) NOT NULL,
        price INTEGER NOT NULL CHECK (price > 0),
        user_id UUID NOT NULL,
        start_date DATE NOT NULL,
        end_date DATE
    );

	GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "$POSTGRES_USER";
	GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "$POSTGRES_USER";
EOSQL