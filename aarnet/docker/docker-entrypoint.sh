#!/bin/bash

config_print() {
	echo ""
	echo "Configured environment variables:"
	echo "`env`" | grep -E "^(RCLONE_WEBDAV|RESTIC_|AWS_)" | sort | while IFS='=' read -r NAME VALUE; do
		if [[ $(echo $NAME | grep -c -E "PASS|SECRET") -ne 0 ]]; then
			echo "  $NAME=**OBFUSCATED**"
		else
			echo "  $NAME=$VALUE"
		fi
	done
	echo ""
}

config_load() {
	if [[ "${RESTIC_REPOSITORY}" != rclone::webdav:* ]] && [[ "${RESTIC_REPOSITORY}" != s3:* ]]; then
		echo "RESTIC_REPOSITORY is required. Must be in the format of 's3:https://hostname.s3.aarnet.edu.au/bucket_name[/path]' or 'rclone::webdav:/path/to/repository'"
		exit 1
	fi

	if [[ "${RESTIC_REPOSITORY}" == s3:* ]]; then
		config_s3
	fi
	if [[ "${RESTIC_REPOSITORY}" == rclone::webdav:* ]]; then
		config_webdav
	fi

	# Oh no, they supplied both.
	if [ "z${RESTIC_PASSWORD_FILE}" != "z" ] && [ "z${RESTIC_PASSWORD}" != "z" ]; then
		# If password file has contents, use it. Otherwise use the password variable.
		if [ -f "${RESTIC_PASSWORD_FILE}" ] && [ `wc -c "{$RESTIC_PASSWORD_FILE}` -gt 0 ]; then
			echo "Both RESTIC_PASSWORD_FILE and RESTIC_PASSWORD are set. RESTIC_PASSWORD_FILE will be used."
			unset RESTIC_PASSWORD
		else
			echo "Both RESTIC_PASSWORD_FILE and RESTIC_PASSWORD are set but RESTIC_PASSWORD_FILE seems empty. RESTIC_PASSWORD will be used."
			unset RESTIC_PASSWORD_FILE
		fi
	fi

	# Check we have a repository password
	if [ "z${RESTIC_PASSWORD_FILE}" != "z" ]; then
		if [ ! -f "${RESTIC_PASSWORD_FILE}"]; then
			echo "${RESTIC_PASSWORD_FILE} does not exist."
			exit 1
		fi
	else
		# If password has been provided, dump it to a file
		if [ "z${RESTIC_PASSWORD}" != "z" ]; then
			export RESTIC_PASSWORD_FILE="/tmp/restic-pwd"
			echo "${RESTIC_PASSWORD}" > /tmp/restic-pwd
		fi
	fi
}

config_s3() {
	echo "S3 endpoint detected."
	if [[ "${AWS_ACCESS_KEY_ID}" == "" ]]; then
		echo "AWS_ACCESS_KEY_ID needs to be set"
		exit 2
	fi
	if [[ "${AWS_SECRET_ACCESS_KEY}" == "" ]]; then
		echo "AWS_SECRET_ACCESS_KEY needs to be set"
		exit 2
	fi
}

config_webdav() {
	echo "WebDAV endpoint detected."

	RCLONE_WEBDAV_VENDOR=${RCLONE_WEBDAV_VENDOR:-owncloud}
	if [[ ${RCLONE_WEBDAV_USER} == "" ]]; then
		echo "RCLONE_WEBDAV_USER needs to be set"
		exit 2
	fi
	if [[ ${RCLONE_WEBDAV_PASS} == "" ]]; then
		echo "RCLONE_WEBDAV_PASS needs to be set"
		exit 2
	fi
	if [[ ${RCLONE_WEBDAV_URL} == "" ]]; then
		echo "RCLONE_WEBDAV_URL needs to be set"
		exit 2
	fi

	export RCLONE_WEBDAV_PASS=$(rclone obscure "${RCLONE_WEBDAV_PASS}")
}

cmd_restic() {
	restic --verbose ${@}
}

config_load
config_print
