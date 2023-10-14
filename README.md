# pistate
![Build Status](https://github.com/shaeberling/pistate/actions/workflows/go.yml/badge.svg)

A simple tool to print the state of your Raspberry Pi.

Currently, only supports the throttled state of `vcgencmd`

Example:

```
$ ./pistate
┌─────────────────┬──────┬──────┐
│ CATEGORY        │ CURR │ PAST │
├─────────────────┼──────┼──────┤
│ Under-voltage   │ ⚠️   │ ⚠️   │
│ Arm freq capped │ ✔️   │ ✔️   │
│ Throttled       │ ⚠️   │ ⚠️   │
│ Soft temp limit │ ✔️   │ ⚠️   │
└─────────────────┴──────┴──────┘
Throttled status [throttled=0xd0005]

```
