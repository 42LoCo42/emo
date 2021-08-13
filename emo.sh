#!/usr/bin/env bash
# shellcheck disable=SC1090 disable=SC2119
set -euo pipefail

path_config_dir="$HOME/.config/emo"
path_data_dir="$HOME/.local/share/emo"
path_hashes_file="$path_data_dir/hashes"
path_songs_dir="$path_data_dir/songs"
path_groups_dir="$path_data_dir/groups"
path_tags_dir="$path_data_dir/tags"

md() {
	mkdir -p "$@"
}

# load config
md "$path_config_dir"
find "$path_config_dir" -type f | while read -r i; do . "$i"; done

declare -A htn=()

die() {
	echo "$@" >&2
	exit 1
}

list-something() {
	path="$1"
	skip=$((${#path} + 2)) # 2 skips leading /
	shift
	find "$path" -mindepth 1 "$@" | while read -r i; do
		tail -c+$skip <<< "$i"
	done
}

build-hash-to-name() {
	((${#htn[@]} > 0)) && return
	while read -r i; do
		hash="$(cmd_name-to-hash "$i")"
		htn["$hash"]="$i"
	done <<< "$(cmd_list-songs)"
}

cmd_cleanup() {
	false && usage ""
	cmd_list-songs | while read -r i; do
		cmd_name-to-hash "$i"
	done > "$path_hashes_file"
	list-something "$path_data_dir" -type d -empty -delete
}

cmd_add-song() {
	(($# < 2)) && usage "<file> <future song name>"
	md "$path_songs_dir/$(dirname "$2")"
	cp "$1" "$path_songs_dir/$2"
	cmd_name-to-hash "$2" >> "$path_hashes_file"
	echo "Added $1 to song library named $2"
}

cmd_del-song() {
	(($# < 1)) && usage "<song name>"
	file="$(mktemp)"
	grep -v "$(cmd_name-to-hash "$1")" "$path_hashes_file" > "$file"
	mv "$file" "$path_hashes_file"
	rm "$path_songs_dir/$1"
	echo "Removed $1 from song library"
}

cmd_name-to-hash() {
	(($# < 1)) && usage "<song name>"
	sha256sum "$path_songs_dir/$1" | awk '{print $1}'
}

cmd_hash-to-name() {
	(($# < 1)) && usage "<hash>"
	build-hash-to-name
	set +u
	name="${htn["$1"]}"
	[ -z "$name" ] && die "Unknown hash: $1"
	set -u
	echo "$name"
}

cmd_list-songs() {
	false && usage ""
	list-something "$path_songs_dir" -type f
}

cmd_add-to-group() {
	(($# < 2)) && usage "<song name> <group name>"
	hash="$(cmd_name-to-hash "$1")"
	dir="$path_groups_dir/$2"
	md "$dir"
	touch "$dir/$hash"
	echo "Added $1 to group $2"
}

cmd_del-from-group() {
	(($# < 2)) && usage "<song name> <group name>"
	hash="$(cmd_name-to-hash "$1")"
	rm "$path_groups_dir/$2/$hash"
	echo "Removed $1 from group $2"
}

cmd_list-groups() {
	false && usage ""
	list-something "$path_groups_dir" -type d
}

cmd_list-in-group() {
	(($# < 1)) && usage "<group name>"
	list-something "$path_groups_dir/$1" -maxdepth 1 -type f | while read -r i; do
		cmd_hash-to-name "$i"
	done

}

cmd_set-tag() {
	(($# < 3)) && usage "<song name> <tag name> <value>"
	hash="$(cmd_name-to-hash "$1")"
	dir="$path_tags_dir/$hash"
	md "$dir"
	echo "$3" >> "$dir/$2"
}

cmd_get-tag() {
	(($# < 2)) && usage "<song name> <tag name>"
	cat "$path_tags_dir/$(cmd_name-to-hash "$1")/$2"
}

cmd_del-tag() {
	(($# < 2)) && usage "<song name> <tag name>"
	rm "$path_tags_dir/$(cmd_name-to-hash "$1")/$2"
}

cmd_list-tags() {
	(($# < 1)) && usage "<song name>"
	list-something "$path_tags_dir/$(cmd_name-to-hash "$1")"
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
	false && usage ""
	declare | grep "^path_"
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
