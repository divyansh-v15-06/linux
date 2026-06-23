import os
import subprocess
import sys

# Define base configuration
BASE_DOCKERFILE = """FROM --platform=linux/386 i386/alpine:3.19
RUN apk add --no-cache bash grep sed gawk coreutils tar gzip openssh-client tcpdump procps util-linux shadow cronie
RUN adduser -D -h /home/guest -s /bin/bash guest && echo "guest:guest" | chpasswd
"""

def run_cmd(cmd, cwd=None):
    print(f"Running: {' '.join(cmd)}")
    res = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)
    if res.returncode != 0:
        print(f"Error: {res.stderr}")
        sys.exit(res.returncode)
    return res.stdout

def main():
    print("=== Starting LinuxQuest Image Builder ===")
    
    # Ensure directories exist
    os.makedirs("server/images", exist_ok=True)
    os.makedirs("images/base", exist_ok=True)
    
    # Write base Dockerfile
    with open("images/base/Dockerfile", "w") as f:
        f.write(BASE_DOCKERFILE)
        
    # Build base image if not exists
    print("Checking for base docker image...")
    try:
        subprocess.run(["docker", "image", "inspect", "linuxquest-base"], check=True, capture_output=True)
        print("Base image 'linuxquest-base' already exists.")
    except subprocess.CalledProcessError:
        print("Building base image...")
        run_cmd(["docker", "build", "--platform", "linux/386", "-t", "linuxquest-base", "images/base"])

    # Loop to build each chapter image
    for ch in range(12):
        print(f"\\n--- Building Chapter {ch} ---")
        temp_root = f"images/chapter-{ch}/rootfs"
        os.makedirs(temp_root, exist_ok=True)
        
        # Chapter 0: Bootcamp
        if ch == 0:
            dossier_dir = os.path.join(temp_root, "home/guest/.dossier")
            os.makedirs(dossier_dir, exist_ok=True)
            with open(os.path.join(dossier_dir, "hint_coords.txt"), "w") as f:
                f.write("ANTARIKSHA{TERMINAL_WAKES_2026}\\n")
            with open(os.path.join(temp_root, "home/guest/welcome.txt"), "w") as f:
                f.write("Welcome to your workstation, Agent Arjun.\\nRead orientation logs in .dossier to find the entry flag.\\n")
                
        # Chapter 1: The Lab
        elif ch == 1:
            os.makedirs(os.path.join(temp_root, "home/guest/logs"), exist_ok=True)
            with open(os.path.join(temp_root, "home/guest/shiva_init.sh"), "w") as f:
                f.write("#!/bin/sh\\n# S.H.I.V.A Initialization payload\\n# Compromised user: cdac_monitor\\n")
            with open(os.path.join(temp_root, "home/guest/logs/incident.log"), "w") as f:
                f.write("Incident Log - 04:02 IST\\nUnusual SUID creation detected at home directory: shiva_init.sh\\n")

        # Chapter 2: The Signal
        elif ch == 2:
            log_dir = os.path.join(temp_root, "var/log/telemetry")
            os.makedirs(log_dir, exist_ok=True)
            with open(os.path.join(log_dir, "tmt-04-58.log"), "w") as f:
                # Generate 50,000 lines
                for i in range(50000):
                    if i % 833 == 0:
                        # Inject SHIVA_PING entry
                        f.write(f"[2026-06-23 04:58:02] SHIVA_PING transmission outbound to 10.48.7.219\\n")
                    else:
                        f.write(f"[2026-06-23 04:58:02] Telemetry beacon heart-rate nominal from station {i % 100}\\n")

        # Chapter 3: The Hunt
        elif ch == 3:
            bin_dir = os.path.join(temp_root, "bin")
            os.makedirs(bin_dir, exist_ok=True)
            # Override ps command to list mock process list
            with open(os.path.join(bin_dir, "ps"), "w") as f:
                f.write("""#!/bin/sh
if [ -f /tmp/killed_3847 ]; then
  echo "PID   USER     TIME  COMMAND"
  echo "1     root      0:00 /bin/sh"
  echo "142   guest     0:00 /usr/sbin/dropbear"
  echo "3846  root      0:00 [kworker/u8:0]"
  echo "3848  root      0:00 [kworker/u8:1]"
else
  echo "PID   USER     TIME  COMMAND"
  echo "1     root      0:00 /bin/sh"
  echo "142   guest     0:00 /usr/sbin/dropbear"
  echo "3846  root      0:00 [kworker/u8:0]"
  echo "3847  root      0:05 [kworker/u8]"
  echo "3848  root      0:00 [kworker/u8:1]"
fi
""")
            # Override kill command to handle target PID
            with open(os.path.join(bin_dir, "kill"), "w") as f:
                f.write("""#!/bin/sh
if [ "$1" = "3847" ] || [ "$2" = "3847" ]; then
  touch /tmp/killed_3847
  echo "Process 3847 (kworker/u8) terminated. Flag: ANTARIKSHA{3847:kworker/u8}"
else
  echo "kill: target PID not found or access denied"
fi
""")

        # Chapter 4: Cronjob of Doom
        elif ch == 4:
            cron_dir = os.path.join(temp_root, "etc/cron.d")
            os.makedirs(cron_dir, exist_ok=True)
            with open(os.path.join(cron_dir, "sysstat"), "w") as f:
                f.write("*/15 * * * * root /usr/local/bin/shiva_beacon.sh\\n")

        # Chapter 5: Permissions
        elif ch == 5:
            bin_dir = os.path.join(temp_root, "usr/bin")
            os.makedirs(bin_dir, exist_ok=True)
            with open(os.path.join(bin_dir, "cdac_stat"), "w") as f:
                f.write("#!/bin/sh\\n# CDAC stats SUID helper\\n")

        # Chapter 6: The Archive
        elif ch == 6:
            var_tmp = os.path.join(temp_root, "var/tmp")
            backup_dir = os.path.join(temp_root, "backup")
            os.makedirs(var_tmp, exist_ok=True)
            os.makedirs(backup_dir, exist_ok=True)
            with open(os.path.join(var_tmp, "staged_data.tar.gz"), "w") as f:
                f.write("staged data")
            with open(os.path.join(backup_dir, "iucaa_baseline.tar.gz"), "w") as f:
                f.write("baseline backup")
            # Custom sha256sum to match predefined hash
            usr_bin = os.path.join(temp_root, "usr/bin")
            os.makedirs(usr_bin, exist_ok=True)
            with open(os.path.join(usr_bin, "sha256sum"), "w") as f:
                f.write("""#!/bin/sh
if echo "$@" | grep -q "staged_data.tar.gz"; then
  echo "e9c0f83d7a8b5e2c000000000000000000000000000000000000000000000000  $1"
else
  busybox sha256sum "$@"
fi
""")

        # Chapter 7: Text Surgeon
        elif ch == 7:
            log_dir = os.path.join(temp_root, "var/log/proxy")
            os.makedirs(log_dir, exist_ok=True)
            with open(os.path.join(log_dir, "access.log"), "w") as f:
                # Write 20,000 lines
                message = "12.9716,77.5946"
                msg_idx = 0
                for i in range(20000):
                    if i % 400 == 0 and msg_idx < len(message):
                        char = message[msg_idx]
                        msg_idx += 1
                        f.write(f'{{"method": "POST", "path": "/api/v2/sync", "user_agent": "Mozilla/5.0 Client Sync {char} Agent"}}/r/n')
                    else:
                        f.write('{"method": "GET", "path": "/index.html", "user_agent": "Mozilla/5.0 Chrome/119.0.0"}/r/n')

        # Chapter 8: The Shell Wars
        elif ch == 8:
            etc_dir = os.path.join(temp_root, "etc")
            os.makedirs(etc_dir, exist_ok=True)
            with open(os.path.join(etc_dir, "nodes.list"), "w") as f:
                f.write("10.0.1.1\\n10.0.1.2\\n10.0.1.3\\n10.0.1.4\\n10.0.1.5\\n")
            with open(os.path.join(temp_root, "home/guest/clean_report.txt"), "w") as f:
                f.write("Completed nodes status: cleaned: 11, unreachable: 2. Flag: ANTARIKSHA{11:2}\\n")

        # Chapter 9: Ghost Signal
        elif ch == 9:
            os.makedirs(os.path.join(temp_root, "home/guest"), exist_ok=True)
            with open(os.path.join(temp_root, "home/guest/dns_traffic.txt"), "w") as f:
                f.write("Host: 10.0.0.4 queried C2 dns query at 10.0.0.15. Flag: ANTARIKSHA{10.0.0.4:10.0.0.15}\\n")

        # Chapter 10: SSH Tunnels
        elif ch == 10:
            os.makedirs(os.path.join(temp_root, "home/guest/.ssh"), exist_ok=True)
            with open(os.path.join(temp_root, "home/guest/.ssh/id_rsa"), "w") as f:
                f.write("-----BEGIN RSA PRIVATE KEY-----\\n")
            with open(os.path.join(temp_root, "home/guest/connection_instructions.txt"), "w") as f:
                f.write("Tunnel through jump hosts. Core DB key hash: d3b07384d113edec\\n")

        # Chapter 11: Final Shutdown
        elif ch == 11:
            os.makedirs(os.path.join(temp_root, "etc/systemd/system"), exist_ok=True)
            with open(os.path.join(temp_root, "etc/systemd/system/shiva.service"), "w") as f:
                f.write("[Service]\\nWatchdogSec=90s\\nExecStart=/usr/bin/shiva_watchdog\\n")
            with open(os.path.join(temp_root, "home/guest/shutdown_manual.txt"), "w") as f:
                f.write("Flag: ANTARIKSHA{shiva_watchdog:f8a7e2b1}\\n")

        # Package using Docker
        print(f"Exporting rootfs base container for Chapter {ch}...")
        temp_workspace = f"temp_ws_{ch}"
        os.makedirs(temp_workspace, exist_ok=True)
        
        # Export rootfs base tarball
        cid_res = subprocess.run(["docker", "create", "linuxquest-base"], capture_output=True, text=True, check=True)
        cid = cid_res.stdout.strip()
        
        # Export and extract
        tar_path = os.path.join(temp_workspace, "base.tar")
        run_cmd(["docker", "export", cid, "-o", tar_path])
        run_cmd(["docker", "rm", cid])
        
        print("Extracting base filesystem...")
        run_cmd(["tar", "-C", temp_workspace, "-xf", tar_path])
        os.remove(tar_path)
        
        # Copy chapter-specific files over base filesystem
        # In python, we copy files manually to merge properly
        print("Merging custom files...")
        run_cmd(["cp", "-r", f"{temp_root}/.", temp_workspace])
        
        # Execute mke2fs to bundle it as an ext2 image
        print("Formatting and writing ext2 image...")
        img_dest = f"server/images/ch{ch}.img"
        if os.path.exists(img_dest):
            os.remove(img_dest)
            
        # Run Docker container to create ext2 image
        abs_cwd = os.getcwd()
        docker_ws = "/workspace"
        
        setup_perms = "chown -R 1000:1000 /tmp/rootfs/home/guest && "
        if ch == 1:
            setup_perms += "chmod 4755 /tmp/rootfs/home/guest/shiva_init.sh && "
        elif ch == 5:
            setup_perms += "chmod 4755 /tmp/rootfs/usr/bin/cdac_stat && "
            
        run_cmd([
            "docker", "run", "--rm",
            "-v", f"{abs_cwd}:{docker_ws}",
            "alpine", "sh", "-c",
            f"apk add --no-cache e2fsprogs && "
            f"mkdir -p /tmp/rootfs && "
            f"cp -a {docker_ws}/{temp_workspace}/. /tmp/rootfs/ && "
            f"{setup_perms}"
            f"dd if=/dev/zero of={docker_ws}/{img_dest} bs=1M count=100 && "
            f"mkfs.ext2 -F -d /tmp/rootfs {docker_ws}/{img_dest}"
        ])
        
        # Clean up temporary workspace
        run_cmd(["rm", "-rf", temp_workspace])
        print(f"Successfully created: {img_dest}")

if __name__ == "__main__":
    main()
