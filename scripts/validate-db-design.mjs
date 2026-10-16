import { readFile } from 'node:fs/promises';
import { resolve, dirname } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

// Validates design consistency, not business FDs or runtime safety.
export async function validateDesign(root = resolve(dirname(fileURLToPath(import.meta.url)), '..')) {
  const docs = resolve(root, 'docs');
  const manifest = JSON.parse(await readFile(resolve(docs, 'db/notitek_3nf.schema.json'), 'utf8'));
  const legacy = await readFile(resolve(docs, 'Notitek - DATABASE v0.8 (LỊCH SỬ).md'), 'utf8');
  const dictionary = await readFile(resolve(docs, 'Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md'), 'utf8');
  const errors = [], tables = new Map(), tick = String.fromCharCode(96);
  const same = (a, b) => JSON.stringify(a) === JSON.stringify(b);
  const check = (condition, message) => { if (!condition) errors.push(message); };
  const identifier = /^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$/;
  let fields = 0;
  for (const [position, table] of manifest.tables.entries()) {
    check(identifier.test(table.name) && !tables.has(table.name), 'Invalid/duplicate table: ' + table.name);
    tables.set(table.name, table);
    const names = new Map(table.fields.map(field => [field.name, field]));
    fields += table.fields.length;
    check(names.size === table.fields.length && Boolean(table.why), 'Duplicate fields/missing showcase: ' + table.name);
    for (const field of table.fields) {
      check(identifier.test(field.name) && ['Y', 'N', 'C'].includes(field.required), 'Invalid field/nullability: ' + table.name + '.' + field.name);
      check(Boolean(field.type && field.meaning && field.example && field.use), 'Incomplete field: ' + table.name + '.' + field.name);
    }
    for (const key of [table.key, ...table.unique, ...table.uniqueIndexes.map(index => index.columns)]) {
      check(key.length > 0 && new Set(key).size === key.length && key.every(name => names.has(name)), 'Invalid key/index: ' + table.name);
    }
    check(table.key.every(name => names.get(name)?.required === 'Y'), 'Nullable PK: ' + table.name);
    for (const key of table.unique) {
      const index = table.uniqueIndexes.find(index => !index.where && same(index.columns, key));
      check(Boolean(index), 'Missing full unique definition: ' + table.name);
      if (key.some(name => names.get(name)?.required !== 'Y')) check(index?.nullsNotDistinct === true, 'Missing unique NULL semantics: ' + table.name);
    }
    const heading = '### 4.' + (position + 1) + ' ' + tick + table.name + tick;
    const start = dictionary.indexOf(heading + '\n');
    const next = dictionary.indexOf('\n### 4.', start + heading.length);
    const section = dictionary.slice(start, next < 0 ? dictionary.indexOf('\n## 5 ', start) : next);
    const docFields = [...section.matchAll(/^\| \x60([^\x60]+)\x60 \|/gm)].map(match => match[1]);
    check(start >= 0 && same(docFields, table.fields.map(field => field.name)), 'Dictionary fields differ: ' + table.name);
  }
  const links = new Set();
  for (const fk of manifest.foreignKeys) {
    const child = tables.get(fk.child), parent = tables.get(fk.parent);
    check(Boolean(child && parent), 'Missing FK table: ' + fk.child + ' -> ' + fk.parent);
    if (!child || !parent) continue;
    const signature = JSON.stringify([fk.child, fk.columns, fk.parent, fk.parentColumns]);
    check(fk.columns.length > 0 && fk.columns.length === fk.parentColumns.length && !links.has(signature), 'Invalid/duplicate FK: ' + signature);
    links.add(signature);
    for (let i = 0; i < fk.columns.length; i++) {
      const cf = child.fields.find(field => field.name === fk.columns[i]);
      const pf = parent.fields.find(field => field.name === fk.parentColumns[i]);
      check(Boolean(cf && pf && cf.type === pf.type), 'FK field/type mismatch: ' + signature);
    }
    check(same(parent.key, fk.parentColumns) || parent.uniqueIndexes.some(index => !index.where && same(index.columns, fk.parentColumns)), 'FK target lacks full uniqueness: ' + signature);
    check(['RESTRICT', 'NO ACTION'].includes(fk.onDelete), 'Unexpected FK delete action: ' + signature);
    if (fk.deferrable) check(fk.initiallyDeferred && fk.onDelete === 'NO ACTION', 'Deferred delete requires NO ACTION: ' + signature);
    const row = '| ' + tick + fk.child + '(' + fk.columns.join(', ') + ')' + tick + ' | ' + tick + fk.parent + '(' + fk.parentColumns.join(', ') + ')' + tick + ' | ' + fk.onDelete;
    check(dictionary.includes(row), 'Dictionary FK differs: ' + signature);
  }
  const legacySection = legacy.slice(legacy.indexOf('## 5 Từ điển'), legacy.indexOf('## 6 Thời gian'));
  let tableName;
  const oldFields = [];
  for (const line of legacySection.split('\n')) {
    const heading = line.match(/^### 5\.\d+ \x60([^\x60]+)\x60/);
    if (heading) tableName = heading[1];
    const field = line.match(/^\| \x60([^\x60]+)\x60 \|/);
    if (field) oldFields.push(tableName + '.' + field[1]);
  }
  const mapped = manifest.oldFields.map(field => field.table + '.' + field.field);
  check(new Set(mapped).size === mapped.length && same(mapped, oldFields), 'Legacy mapping missing/duplicated/out of order');
  for (const field of manifest.oldFields) {
    check(Boolean(field.destination), 'Missing destination: ' + field.table + '.' + field.field);
    const row = '| ' + tick + field.table + '.' + field.field + tick + ' | ' + field.destination + ' |';
    check(dictionary.includes(row), 'Dictionary mapping differs: ' + field.table + '.' + field.field);
  }
  const counts = { tables: tables.size, fields, legacyFields: oldFields.length, foreignKeys: manifest.foreignKeys.length };
  for (const [name, count] of Object.entries(counts)) check(count === manifest.expectedCounts[name], 'Unexpected count: ' + name);
  if (errors.length) throw new Error(errors.join('\n'));
  return { version: manifest.version, ...counts, checks: 'keys, NULL semantics, FK types/targets/actions, showcases, dictionary, all legacy mappings' };
}

if (typeof process !== 'undefined' && process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href) {
  validateDesign().then(result => console.log(JSON.stringify(result, null, 2))).catch(error => {
    console.error(error.message);
    process.exitCode = 1;
  });
}
