# TaxHub

HTTP service that maps NCBI Taxonomy IDs to scientific names. Data is sourced from the [NCBI taxdump](https://ftp.ncbi.nih.gov/pub/taxonomy/).

## Getting Started

### 1. Data Pipeline

Download and extract the NCBI taxdump, then generate the CSV used by the service:

```bash
cd data
make
```

This downloads `taxdump.tar.gz` from NCBI, extracts it, and writes `taxhub/scientific_names.csv`. Re-running `make` is a no-op if the data is already present.

```bash
make clean  # remove downloaded data and generated CSV
```

### 2. Run the Service

**Local:**
```bash
cd taxhub
go run ./cmd/taxhub
```

**Docker:**
```bash
cd taxhub
docker build -t taxhub:local .
docker run -p 8080:8080 taxhub:local
```

**Kubernetes (local):**
```bash
kubectl apply -f taxhub/k8s/taxhub.yaml
```

## API

### `GET /healthz`
Returns service status and number of loaded taxonomy entries.

```json
{ "status": "ok", "count": 2498561 }
```

### `GET /taxon/:tax_id`
Returns the scientific name for a given NCBI Taxonomy ID.

```json
{ "tax_id": "9606", "name": "Homo sapiens" }
```

Returns `404` if the ID is not found.

## Configuration

| Env var    | Default                      | Description              |
|------------|------------------------------|--------------------------|
| `ADDR`     | `:8080`                      | Listen address           |
| `CSV_PATH` | `scientific_names.csv`       | Path to the taxonomy CSV |

## Performance Testing

Load tests are written in [k6](https://k6.io):

```bash
k6 run performance-testing/k6.js
```
