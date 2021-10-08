#!/usr/bin/env bash
set -euo pipefail

path_config_dir="$HOME/.config/emo"
path_data_dir="$HOME/.local/share/emo"
path_songs_dir="$path_data_dir/songs"
path_groups_dir="$path_data_dir/groups"

md() {
	mkdir -p "$@"
}

# load config
md "$path_config_dir"
# shellcheck disable=SC1090
find "$path_config_dir" -type f | while read -r i; do . "$i"; done

die() {
	echo "$@" >&2
	exit 1
}

list-something() {
	path="$1"
	shift
	(cd "$path"; find . -mindepth 1 "$@") | cut -c 3-
}

cmd_cleanup() {
	false && usage ""
	list-something "$path_data_dir" -type d -empty -delete
}

cmd_add-song() {
	(($# < 2)) && usage "<file> <future song name>"
	md "$path_songs_dir/$2"
	cp "$1" "$path_songs_dir/file"
	echo "Added $1 to song library named $2"
}

cmd_del-song() {
	(($# < 1)) && usage "<song name>"
	rm -rf "$path_songs_dir/${1:?}"
	echo "Removed $1 from song library"
}

cmd_list-songs() {
	false && usage ""
	list-something "$path_songs_dir" -type f -name 'file' | sed 's|/file$||'
}

cmd_add-to-group() {
	(($# < 2)) && usage "<song name> <group name>"
	md "$path_groups_dir"
	: >> "$path_groups_dir/$2"
	sort -u <(echo "$1") -o "$path_groups_dir/$2"{,}
	echo "Added $1 to group $2"
}

cmd_del-from-group() {
	(($# < 2)) && usage "<song name> <group name>"
	file="$(mktemp)"
	if grep -Fv "$1" "$path_groups_dir/$2" > "$file"; then
		mv "$file" "$path_groups_dir/$2"
	else
		rm "$path_groups_dir/$2"
	fi
	echo "Removed $1 from group $2"
}

cmd_list-groups() {
	false && usage ""
	ls -1 "$path_groups_dir"
}

cmd_list-in-group() {
	(($# < 1)) && usage "<group name>"
	cat "$path_groups_dir/$1"

}

cmd_set-tag() {
	(($# < 3)) && usage "<song name> <tag name> <value>"
	echo "$3" > "$path_songs_dir/$1/$2"
}

cmd_get-tag() {
	(($# < 2)) && usage "<song name> <tag names>"
	song="$1"
	shift
	cat "${@/#/"$path_songs_dir/$song/"}"
}

cmd_del-tag() {
	(($# < 2)) && usage "<song name> <tag name>"
	rm "$path_songs_dir/$1/$2"
}

cmd_list-tags() {
	(($# < 1)) && usage "<song name>"
	list-something "$path_songs_dir/$1" -type f -not -name "file"
}

cmd_getcmds() {
	false && usage "[with-args]"
	with_args="${1:-}"

	declare -F | sed -En "s|declare -f cmd_(.*)$|\1|p" \
	| if [ -n "$with_args" ]; then
		while read -r c; do
			echo -n "$c "
			type "cmd_$c" | sed "4q;d" | sed -En 's|^.*usage "(.*)";$|\1|p'
		done
	else
		tee
	fi
}

cmd_getpaths() {
	false && usage "[path to get]"
	path="${1:-}"
	output="$(declare | grep "^path_$path")"
	[ -n "$path" ] && output="${output#*=}"
	echo "$output"
}

usage() {
	die "Usage: $0 $cmd $1"
}

(($# < 1)) && {
	echo "Usage: $0 <subcommand>"
	echo "Available subcommands:"
	cmd_getcmds 1
	die -n ""
} >&2

cmd="$1"
cmd_getcmds | grep -q "^$cmd$" || die "Invalid subcommand $cmd"
shift
"cmd_$cmd" "$@"
