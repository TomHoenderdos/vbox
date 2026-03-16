# Creating vbox profiles

## Profile structure

A profile is a single `.sh` file in `~/.vbox/profiles/` with two functions:

```bash
#!/usr/bin/env bash
# Short description: one line that shows up in `vbox profile list`.

profile_ports() {
  # Define port forwards: guest:default_host:label
  # One per line. Omit this function or leave empty if no ports needed.
  echo "8080:8080:MyService"
  echo "9090:9090:Admin"
}

profile_provision() {
  # Output shell commands that run as root inside the VM during provisioning.
  cat <<'PROVISION'
    apt-get install -y my-package
    systemctl enable my-service
PROVISION
}
```

### Rules

1. **File name** = profile name. `foo.sh` becomes `--profile foo`
2. **Line 2** must be a `# comment` — it's the description shown in `vbox profile list`
3. **`profile_ports()`** outputs `guest:default_host:label` lines. During `vbox init`, the user is asked to confirm or change each host port. Omit or leave empty for profiles that don't expose ports
4. **`profile_provision()`** outputs shell commands. These run as root in the VM. Use `su - vagrant -c '...'` for commands that should run as the vagrant user
5. Use `'PROVISION'` (quoted) for heredocs with literal `$` signs, unquoted `PROVISION` when you need variable interpolation

### Reading tool versions

For language profiles, read versions from `.tool-versions` (asdf convention):

```bash
profile_provision() {
  local my_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      my_version=$(awk '/^mytool/ {print $2}' "$tvf")
      [[ -n "$my_version" ]] && break
    fi
  done
  my_version="${my_version:-1.0.0}"  # fallback default

  cat <<PROVISION
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add mytool'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install mytool ${my_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global mytool ${my_version}'
PROVISION
}
```

### Services that accept connections from the host

For database/service profiles, make sure to:

1. Define the port in `profile_ports()`
2. Configure the service to listen on `0.0.0.0` (not just localhost)
3. Create credentials if applicable

```bash
profile_ports() {
  echo "5432:15432:PostgreSQL"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y postgresql
    systemctl enable postgresql
    systemctl start postgresql

    # Allow connections from host
    echo "listen_addresses = '*'" >> /etc/postgresql/16/main/postgresql.conf
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/16/main/pg_hba.conf
    systemctl restart postgresql
PROVISION
}
```

## Testing your profile

```bash
# Create a test project with your profile
vbox init TestProfile --profile myprofile

# Check the generated Vagrantfile
cat ~/Projects/TestProfile/Vagrantfile

# Test provisioning
vbox up

# Clean up
vbox down -v
rm -rf ~/Projects/TestProfile
```

## Claude Code prompt for generating profiles

Use this prompt with Claude Code to generate a new profile:

```
Create a new vbox profile at ~/.vbox/profiles/<name>.sh

A vbox profile is a bash script with two functions:
- profile_ports(): echo "guest:default_host:label" lines for port forwards (or empty)
- profile_provision(): output shell commands that run as root in an Ubuntu 24.04 ARM64 VM

Line 2 must be a # comment with a short description.
Use `su - vagrant -c '...'` for user-level commands.
For asdf-managed tools, read versions from .tool-versions files.
For services, bind to 0.0.0.0 and create default dev credentials.

Look at existing profiles in ~/.vbox/profiles/ for reference.

I need a profile for: <describe what you need>
```
