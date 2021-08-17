#!/usr/bin/env bash
set -euo pipefail

declare -a queue

error() {
	echo "error" "$@"
}

push_queue() {
	queue+=("$1")
}

pop_queue() {
	[ -z "${queue[*]}" ] && error "empty queue" && return
	echo "next ${queue[0]}"
	queue=("${queue[@]:1}")
}

print_queue() {
	[ -z "${queue[*]}" ] && error "empty queue" && return
	for i in "${!queue[@]}"; do
		echo "queued $i ${queue[i]}"
	done
}

coproc yell { exec yell "$@"; }
declare yell_PID
trap 'kill $yell_PID' EXIT
exec 1>&"${yell[1]}"- # write stdout to yell

while read -ru "${yell[0]}" cmd args; do
	case "$cmd" in
		exit)  echo "exit"; break ;;
		queue) print_queue ;;
		add)   push_queue "$args" ;;
		next)  pop_queue ;;
		clear) queue=() ;;
		*)     error "unknown command"
	esac
done
