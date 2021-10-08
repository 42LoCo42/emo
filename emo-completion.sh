#!/bin/bash

_emo() {
	# read subcommands
	declare -A subcommands
	while read -r cmd args; do
		subcommands["$cmd"]="$(
		sed -E "
			s|<([^ >]+) ?[^>]*>|\1|g;
			s| |,|g;
		" <<< "$args" | tr -d "[:space:]"
	)"
	done <<< "$(emo getcmds 1)"

	w="${#COMP_WORDS[@]}"
	case "$w" in
		1) ;;
		2)
			mapfile -t COMPREPLY <<< "$(compgen -W "${!subcommands[*]}" "${COMP_WORDS[1],,}")"
			;;
		*)
			IFS="," read -rd '' -a args <<< "${subcommands["${COMP_WORDS[1]}"]}"
			arg="${args[$((w - 3))]}"
			current_word="${COMP_WORDS[$((w - 1))]}"
			mapfile -t COMPREPLY <<< "$(IFS=$'\n' _emo_list)"
			[ -z "${COMPREPLY[*]}" ] && COMPREPLY=()
			;;
	esac
}

_emo_list() {
	case "$arg" in
		song*)   compgen -W "$(emo list-songs)" "$current_word" ;;
		group*)  compgen -W "$(emo list-groups)" "$current_word" ;;
		tag*)    compgen -W "$(emo list-tags "${COMP_WORDS[2]}")" "$current_word" ;;
	esac
}

complete -o bashdefault -o default -F _emo emo
