.PHONY: all clean

TAXDUMP_DIR := taxdump
TAXDUMP_URL := https://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz
TAXDUMP_TMP := /tmp/taxdump.tar.gz
CSV_OUT     := ../taxhub/scientific_names.csv

all: $(CSV_OUT)

# download dumb if it's not there already
$(TAXDUMP_DIR)/names.dmp:
	curl -L -C - -o "$(TAXDUMP_TMP)" "$(TAXDUMP_URL)"
	mkdir -p "$(TAXDUMP_DIR)"
	tar -zxf "$(TAXDUMP_TMP)" -C "$(TAXDUMP_DIR)"

# check all scientific names, and trim, write to csv
$(CSV_OUT): $(TAXDUMP_DIR)/names.dmp
	awk -F '|' '\
	BEGIN { print "tax_id,name" } \
	{ \
		tax_id = $$1; name = $$2; type = $$4; \
		gsub(/^[ \t]+|[ \t]+$$/, "", tax_id); \
		gsub(/^[ \t]+|[ \t]+$$/, "", name); \
		gsub(/^[ \t]+|[ \t]+$$/, "", type); \
		if (type == "scientific name") { \
			gsub(/"/, "\"\"", name); \
			print tax_id ",\"" name "\"" \
		} \
	}' "$<" > "$@"
	@echo "CSV written to: $@"

clean:
	rm -rf "$(TAXDUMP_DIR)" "$(TAXDUMP_TMP)" "$(CSV_OUT)"