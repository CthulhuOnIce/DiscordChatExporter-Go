#!/usr/bin/env python3

import os
import argparse

# Default output binary, will change to {DEFAULT_OUTPUT}.exe if OS is windows
DEFAULT_OUTPUT = "build/dce"  

# Linux Colors
class l_colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    WARNING = '\033[93m'
    RED = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

# Windows Colors
class w_colors:
    HEADER = ''
    BLUE = ''
    GREEN = ''
    WARNING = ''
    RED = ''
    ENDC = ''
    BOLD = ''

if os.name == 'nt':
    colors = w_colors
else:
    colors = l_colors


# build.py --os=linux --arch=amd64 --noformat --notest --nobuild
# build.py --os=windows --arch=amd64
# build.py --o=build/dce.exe --os=windows --arch=amd64
argparser = argparse.ArgumentParser(description="Build the DCE binary")
argparser.add_argument("--os", help="OS to build for", default="linux",)
argparser.add_argument("--arch", help="Arch to build for", default="amd64")
argparser.add_argument("--noformat", help="Skip go fmt", action="store_true")
argparser.add_argument("--notest", help="Skip tests", action="store_true")
argparser.add_argument("--nobuild", help="Skip build", action="store_true")
argparser.add_argument("-o", help="Output binary name", default=DEFAULT_OUTPUT)
argparser.add_argument("-v", help="Enable Verbose Mode", default=False, action="store_true")
args = argparser.parse_args()

output = args.o

# set GOOS and GOARCH
if args.v:
    print(colors.BLUE + f"build.py: Setting GOOS={args.os} && GOARCH={args.arch}" + colors.ENDC)
os.environ["GOOS"] = args.os
os.environ["GOARCH"] = args.arch

# if output binary is default and OS is windows, change output binary to dce.exe
if output == DEFAULT_OUTPUT and args.os == "windows":
    output = "build/dce.exe"

# run go fmt
if not args.noformat:
    command = "go fmt {}./...".format("-n -x " if args.v else "")
    if args.v:
        print(colors.BLUE + f"build.py: Running: {command}" + colors.ENDC)
        
    if os.system(command):
        print(colors.RED + "Go fmt failed" + colors.ENDC)
        exit(1)
    else:
        print(colors.GREEN + "Go fmt successful" + colors.ENDC)
else:
    print(colors.WARNING + "Skipping go fmt" + colors.ENDC)

# run go test
if not args.notest:
    command = "go test {}./...".format("-v " if args.v else "")
    if args.v:
        print(colors.BLUE + f"build.py: Running: {command}" + colors.ENDC)
    
    if os.system(command):
        print(colors.RED + "Tests Failed" + colors.ENDC)
        exit(2)
    else:
        print(colors.GREEN + "Tests Successful" + colors.ENDC)
else:
    print(colors.WARNING + "Skipping tests" + colors.ENDC)

if not args.nobuild:
    command = f"go build {'-v ' if args.v else ''}-o {output}"
    if args.v:
        print(colors.BLUE + f"build.py: Running: {command}" + colors.ENDC)
    
    if os.system(command):
        print(colors.RED + "Build Failed" + colors.ENDC)
        exit(3)
    else:
        print(colors.GREEN + f"Build Successful: {output}" + colors.ENDC)
else:
    print(colors.WARNING + "Skipping build" + colors.ENDC)

print(colors.GREEN + f"All tasks completed." + colors.ENDC)

