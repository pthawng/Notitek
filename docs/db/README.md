# Database design artifacts

Logical design v0.9: 56 tables, 260 fields, 63 foreign keys; complete mapping of all 193 fields from v0.8. No business-table migration has been applied.

- [Dictionary](<../Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>) explains keys, FDs, examples, retention and transactions.
- [Migration plan](<../Notitek - THIẾT KẾ DATABASE 3NF VÀ AN TOÀN MIGRATION.md>) covers initial creation, conditional backfill, reconciliation, rollback and acceptance cases.
- [Historical v0.8](<../Notitek - DATABASE v0.8 (LỊCH SỬ).md>) is the unchanged baseline for field reconciliation.
- [Manifest](notitek_3nf.schema.json) describes the logical design; it is not executable DDL. Types annotated as JSON document/scalar need explicit SQL checks where applicable. unique describes full-table uniqueness; uniqueIndexes specifies NULL treatment and partial predicates. Only full unique keys can be FK targets.

Run from the repository root: node scripts/validate-db-design.mjs. The validator checks identifiers, keys, nullability, FK types/targets/actions, showcases, dictionary consistency, and complete legacy field mapping. It does not prove business FDs, semantic correctness of mapping expressions, or crash/TTL/concurrency safety; those require the review and PostgreSQL acceptance cases.
