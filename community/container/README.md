# Container Security Patterns

Detects Dockerfile and container orchestration misconfigurations.

## Patterns

- `dockerfile-privileged-mode`: Privileged containers
- `dockerfile-exposed-socket`: Docker socket mounts
- `dockerfile-running-as-root`: Containers running as root
- `dockerfile-cap-add-all`: Excessive Linux capabilities

## References

- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)
- [NIST Container Security Guide](https://csrc.nist.gov/publications/detail/sp/800-190/final)
