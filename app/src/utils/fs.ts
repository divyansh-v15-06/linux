export interface FSFile {
  type: 'file';
  content: string;
  permissions?: string;
}

export interface FSDirectory {
  type: 'dir';
  children: Record<string, FSNode>;
  permissions?: string;
}

export type FSNode = FSFile | FSDirectory;

// Initial mock virtual filesystem structure
export const initialFS: FSDirectory = {
  type: 'dir',
  permissions: 'drwxr-xr-x',
  children: {
    home: {
      type: 'dir',
      permissions: 'drwxr-xr-x',
      children: {
        guest: {
          type: 'dir',
          permissions: 'drwxr-xr-x',
          children: {
            help: {
              type: 'file',
              permissions: '-rw-r--r--',
              content: `LinuxQuest Shell Helper v1.0
=============================
Available Commands:
  ls          List files/folders
  cd <dir>    Change directory
  pwd         Show current path
  cat <file>  Read file content
  clear       Clear terminal
  login       Authenticate via Google OAuth
  whoami      Show logged in user
  history     Show command history
  start       Launch CheerpX WASM sandbox
  submit <f>  Submit mission flag
  hint        Get a mission hint
  status      Show quest progress
  man <topic> Show manual pages (e.g. man shiva, man arjun)

Pro Tip: Copy-paste is disabled to help you build real command line muscle memory.
`
            },
            "hint_coords.txt": {
              type: 'file',
              permissions: '-rwxr-xr-x',
              content: `Coordinates extracted: 12.9716° N, 77.5946° E
Location: SAC Ground Station, Bangalore.

Mission Flag: ANTARIKSHA{TERMINAL_WAKES_2026}
Submit this flag using: submit ANTARIKSHA{TERMINAL_WAKES_2026}
`
            },
            profile: {
              type: 'file',
              permissions: '-rw-r--r--',
              content: `USER: guest
LEVEL: 0
XP: 0 / 1000
ELO: 800 (Newcomer)
STREAK: 0 days

Status: Guest. Type 'login' to authenticate.
`
            },
            leaderboard: {
              type: 'file',
              permissions: '-rw-r--r--',
              content: `LinuxQuest Global Leaderboard
=============================
Rank  Username         Elo    XP
1     aditya_sysadm    2450   89,200  [Wizard]
2     priya_sec        2380   78,500  [Engineer]
3     rohit_cyber      2110   54,300  [Engineer]
4     harsh_kernel     1950   41,200  [Sysadmin]
5     kiran_bash       1820   36,800  [Sysadmin]

Status: Showing top 5. Login to see full ranking.
`
            },
            daily: {
              type: 'file',
              permissions: '-rw-r--r--',
              content: `Daily Challenge - 2026-06-23
============================
Mission: Finding the Hidden Config
Find a configuration file named "debug.conf" modified within the last 24 hours under /var.
Command Hint: find /var -name "debug.conf" -mmin -1440

Reward: +150 XP, +10 Elo
`
            },
            missions: {
              type: 'dir',
              permissions: 'drwxr-xr-x',
              children: {
                README: {
                  type: 'file',
                  permissions: '-rw-r--r--',
                  content: `Active Campaigns
================
1. Operation Antariksha [ACTIVE]
   A rogue AI has hijacked India's satellite network.
   Contains 11 chapters across key Indian centers.

Type 'cd antariksha' to view missions.
`
                },
                antariksha: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Operation Antariksha — Mission Map
==================================
           🇮🇳 INDIA 🇮🇳

          [ Srinagar ] (ch10: SSH Tunnels)
               |
          [ New Delhi ] (ch9: Ghost Signal)
               |
          [ Ahmedabad ] (ch5: Permissions) -- [ Kolkata ] (ch6: The Archive)
               |                                  |
          [ Mumbai ] (ch3: The Hunt)       [ Sriharikota ] (ch11: Final Shutdown)
               |                                  |
          [ Goa ] (ch4: Cronjob of Doom)   [ Chennai ] (ch2: The Signal)
               \\                                 /
                [ Bangalore ] (ch1: The Lab) ---/
                     |
                [ Bootcamp ] (ch0: The Lab - SAC)

Type 'cd ch0' to begin the Bootcamp.
`
                    },
                    ch0: {
                      type: 'dir',
                      permissions: 'drwxr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch0: Bootcamp
Location: Space Applications Centre (SAC), Ahmedabad
---------------------------------------------------
Welcome to Operation Antariksha.
Our satellite communications are experiencing micro-anomalies.
Before we assign you to direct intercept teams, you must verify your local tools.

Task:
Explore the local filesystem, find the hidden coordinates file, and submit the flag.
Flag is in format: ANTARIKSHA{...}

Type 'cat objectives' to view the checkpoints.
`
                        },
                        objectives: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `Bootcamp Objectives:
1. Locate "coordinates.txt" somewhere in the filesystem
2. Read the file coordinates to get the flag
3. Type 'submit ANTARIKSHA{...}' with the flag
`
                        }
                      }
                    },
                    ch1: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch1: The Lab
Location: Bangalore Command Centre
----------------------------------
The primary satellite telemetry stream has been shut down by a payload signed "SHIVA".
We need to find the payload's source configuration file on the SAC-BLR-01 server.

Objectives:
1. Boot the server terminal ('start')
2. Search the filesystem for any configs containing SHIVA
3. Submit the flag hash found inside the config

Status: LOCKED (Complete ch0 first)
`
                        }
                      }
                    },
                    ch2: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch2: The Signal
Location: Chennai
-----------------
Identify the destination IP of the SHIVA_PING packet stream in the log telemetry file.

Status: LOCKED (Complete ch1 first)
`
                        }
                      }
                    },
                    ch3: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch3: The Hunt
Location: Mumbai
----------------
A rogue processes is sending password credentials. Find the correct PID and name and kill it.

Status: LOCKED (Complete ch2 first)
`
                        }
                      }
                    },
                    ch4: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch4: Cronjob of Doom
Location: Delhi
---------------
SHIVA has planted a persistence backdoor in the cron jobs. Disable it and fetch schedule.

Status: LOCKED (Complete ch3 first)
`
                        }
                      }
                    },
                    ch5: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch5: Permissions
Location: Hyderabad
-------------------
A modified SUID binary allows root escalation. Find and modify permissions, lock compromised user.

Status: LOCKED (Complete ch4 first)
`
                        }
                      }
                    },
                    ch6: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch6: The Archive
Location: Pune
--------------
Find the compressed file staging archive, verify its integrity and compare hashes with baseline.

Status: LOCKED (Complete ch5 first)
`
                        }
                      }
                    },
                    ch7: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch7: Text Surgeon
Location: Kolkata
-----------------
Extract characters from proxy logs and decode targeting coordinates.

Status: LOCKED (Complete ch6 first)
`
                        }
                      }
                    },
                    ch8: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch8: The Shell Wars
Location: Ahmedabad
-------------------
Write a robust shell script to clean cron entries across a multi-node cluster securely.

Status: LOCKED (Complete ch7 first)
`
                        }
                      }
                    },
                    ch9: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch9: Ghost Signal
Location: Chennai
-----------------
Listen on network interfaces and identify internal and external C2 IP addresses from raw DNS traffic.

Status: LOCKED (Complete ch8 first)
`
                        }
                      }
                    },
                    ch10: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch10: SSH Tunnels
Location: Remote
----------------
Tunnel through three secure jump hosts using private keys to download memory dump database.

Status: LOCKED (Complete ch9 first)
`
                        }
                      }
                    },
                    ch11: {
                      type: 'dir',
                      permissions: 'dr-xr-xr-x',
                      children: {
                        brief: {
                          type: 'file',
                          permissions: '-rw-r--r--',
                          content: `[MISSION DOSSIER] ch11: Final Shutdown
Location: Sriharikota
---------------------
Inspect systemd service dependencies, disable watchdog timers, and stop services in precise sequence.

Status: LOCKED (Complete ch10 first)
`
                        }
                      }
                    }
                  }
                }
              }
            },
            tracks: {
              type: 'dir',
              permissions: 'drwxr-xr-x',
              children: {
                README: {
                  type: 'file',
                  permissions: '-rw-r--r--',
                  content: `Specialized Practice Tracks
===========================
Improve your skills in dedicated categories:
- cd scripting    (Bash scripting)
- cd vim          (Vim commands)
- cd git          (Git workflow)
- cd docker       (Docker administration)
- cd kubernetes   (Kubernetes deployment)
`
                },
                scripting: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Bash Scripting Track
====================
Master automations, loops, variables, and parsing.
`
                    }
                  }
                },
                vim: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Vim Editor Track
================
Master quick edits, navigation keys, macros, and search.
`
                    }
                  }
                },
                git: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Git Control Track
=================
Master branching, merging, rebase, reflogs, and hooks.
`
                    }
                  }
                },
                docker: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Docker Sandbox Track
====================
Master image building, networking, multi-stage files, and compose.
`
                    }
                  }
                },
                kubernetes: {
                  type: 'dir',
                  permissions: 'drwxr-xr-x',
                  children: {
                    README: {
                      type: 'file',
                      permissions: '-rw-r--r--',
                      content: `Kubernetes Cluster Track
========================
Master pods, services, ingress, configuration secrets, and configmaps.
`
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
};

// Helper function to navigate virtual filesystem
export function getNodeByPath(root: FSDirectory, pathParts: string[]): FSNode | null {
  let current: FSNode = root;
  for (const part of pathParts) {
    if (part === '' || part === '.') continue;
    if (current.type !== 'dir') return null;
    const next: FSNode = current.children[part];
    if (!next) return null;
    current = next;
  }
  return current;
}
