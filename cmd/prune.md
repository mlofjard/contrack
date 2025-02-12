==Long==

# Prune

Remove old digests from repositories

====

==Example==

## Remove all digests except the latest 3: 
```bash 
contrack prune --keep 3 manifest.yaml 
```
Supplying the `--dry-run` flag will print out a list of digests  
to keep/remove, but not actually remove them.
```bash 
contrack prune --keep 3 --dry-run manifest.yaml 
```

## Remove digests older than current year: 
```bash
contrack prune --cut-off-date 2025-01-01 manifest.yaml 
```
====
