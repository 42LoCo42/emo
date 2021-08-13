#!/bin/bash

_emo() {
	# read paths
	declare -A paths
	oldifs="$IFS"
	IFS="="
	while read -r name path; do
		paths["$name"]="$path"
	done <<< "$(emo getpaths)"
	IFS="$oldifs"

	# read subcommands
	declare -A subcommands
	while read -r cmd args; do
		subcommands["$cmd"]="$(
			sed -E '
				s|> <|,|g;
				s|[><]||g;
				s| |-|g;
				s|hash|path_hashes_file|g;
				s|([^-,]*)-name|path_\1s_dir|g;
			' <<< "$args" \
				| tr "," "\n" \
				| while read -r word; do
				[ -z "$word" ] && break
				echo "${paths["$word"]}"
			done \
				| tr "\n" ","
		)"
	done <<< "$(emo getcmds 1)"

	w="${#COMP_WORDS[@]}"
	case "$w" in
		1) ;;
		2)
			mapfile -t COMPREPLY <<< "$(compgen -W "${!subcommands[*]}" "${COMP_WORDS[1]}")"
			;;
		*)
			IFS="," read -rd '' -a parts <<< "${subcommands["${COMP_WORDS[1]}"]}"
			part="${parts[$((w - 3))]}"
			mapfile -t COMPREPLY <<< "$(IFS=$'\n' compgen -W "$(_emo_list)" "${COMP_WORDS[$((w - 1))]}")"
			;;
	esac
}

_emo_list() {
	case "$part" in
		*songs*)  emo list-songs ;;
		*groups*) emo list-groups ;;
		*hashes*) cat "$part" ;;
		*tags*)   emo list-tags "${COMP_WORDS[2]}" ;;
	esac
}

_emo_strip_pfx() {
	while read -r i; do
		tail -c+$((${#part} + 2)) <<< "$i"
	done
}

 complete -F _emo emo
