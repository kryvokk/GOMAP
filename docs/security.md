# Security

GoMalScan is designed for static inspection. An analyzed sample is untrusted data: the application must not execute it, load it as code, or send it to an external service automatically. Future parsers should bound file sizes and handle malformed inputs with errors rather than unsafe assumptions.
