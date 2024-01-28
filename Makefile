create-migration:
	$(eval timestamp := $(shell date +%s))
	touch db/$(timestamp)_${name}.up.sql
	touch db/$(timestamp)_${name}.down.sql

up-migration:
	migrate --path=db/ \
			--database ${DATABASE_URL} up

rollback-migration:
	migrate --path=db/ \
			--database ${DATABASE_URL} down


