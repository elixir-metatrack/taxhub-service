import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 200,
  duration: '30s',
};

const BASE_URL = 'http://localhost:8080';

const TAXON_IDS = [
  9606,   // Homo sapiens
  10090,  // Mus musculus
  10116,  // Rattus norvegicus
  7227,   // Drosophila melanogaster
  3702,   // Arabidopsis thaliana
  562,    // Escherichia coli
  6239,   // Caenorhabditis elegans
  7955,   // Danio rerio
  9031,   // Gallus gallus
  9913,   // Bos taurus
];

export default function () {
  const id = TAXON_IDS[Math.floor(Math.random() * TAXON_IDS.length)];
  const res = http.get(`${BASE_URL}/taxon/${id}`);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'has tax_id': (r) => r.json('tax_id') !== undefined,
  });
}
