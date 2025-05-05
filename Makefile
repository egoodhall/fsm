.PHONY: generate migration grammar

test: generate
	go test ./...

clean:
	rm -f fsm.db

# Generate SQLC code
generate:
	go generate ./...

# Create a new migration file
migration:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $${name} sql 
