.PHONY: generate migration

example: generate
	go run example/main.go

clean:
	rm -f fsm.db

# Generate SQLC code
generate:
	go generate ./...

# Create a new migration file
migration:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $${name} sql 
