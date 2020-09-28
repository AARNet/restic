#!/bin/bash

config_print() {
	echo ""
	echo "Configured environment variables:"
	echo "`env`" | grep -E "^(RCLONE_WEBDAV|RESTIC_|AWS_)" | sort | while IFS='=' read -r NAME VALUE; do
		if [[ $(echo $NAME | grep -c -E "PASS|SECRET") -ne 0 ]]; then
			echo "  $NAME=**READACTED**"
		else
			echo "  $NAME=$VALUE"
		fi
	done
	echo ""
}

config_load() {
	if [[ $(echo "${RESTIC_REPOSITORY}" | grep -c -E "^(rclone::webdav:|s3:)") -eq 0 ]]; then
		echo "RESTIC_REPOSITORY is required. Must be in the format of 's3:https://hostname.s3.aarnet.edu.au/bucket_name[/path]' or 'rclone::webdav:/path/to/repository'"
		cmd_output_status 2 "RESTIC_REPOSITORY_missing"
	fi

	if [[ "${RESTIC_REPOSITORY}" == s3:* ]]; then
		config_s3
	fi

	if [[ "${RESTIC_REPOSITORY}" == rclone::webdav:* ]]; then
		config_webdav
	fi

	# Oh no, they supplied both.
	if [[ "z${RESTIC_PASSWORD_FILE}" != "z" ]] && [[ "z${RESTIC_PASSWORD}" != "z" ]]; then
		# If password file has contents, use it. Otherwise use the password variable.
		if [[ -f "${RESTIC_PASSWORD_FILE}" ]] && [[ $(wc -c "{$RESTIC_PASSWORD_FILE}") -gt 0 ]]; then
			echo "Both RESTIC_PASSWORD_FILE and RESTIC_PASSWORD are set. RESTIC_PASSWORD_FILE will be used."
			unset RESTIC_PASSWORD
		else
			echo "Both RESTIC_PASSWORD_FILE and RESTIC_PASSWORD are set but RESTIC_PASSWORD_FILE seems empty. RESTIC_PASSWORD will be used."
			unset RESTIC_PASSWORD_FILE
		fi
	fi

	# Check we have a repository password
	if [[ "z${RESTIC_PASSWORD_FILE}" != "z" ]]; then
		if [[ ! -f "${RESTIC_PASSWORD_FILE}" ]]; then
			echo "${RESTIC_PASSWORD_FILE} does not exist."
			cmd_output_status 1 "RESTIC_PASSWORD_FILE_does_not_exist"
		fi
	else
		# If password has been provided, dump it to a file
		if [[ "z${RESTIC_PASSWORD}" != "z" ]]; then
			export RESTIC_PASSWORD_FILE="/tmp/restic-pwd"
			echo "${RESTIC_PASSWORD}" > /tmp/restic-pwd
		fi
	fi

	# Check if the directory to backup exists
	if [[ ! -d "${RESTIC_BACKUP_SOURCE}" ]]; then
		echo "Backup source '${RESTIC_BACKUP_SOURCE}' does not exist."
		cmd_output_status 1 "RESTIC_BACKUP_SOURCE_does_not_exist"
	fi
}

config_s3() {
	echo "S3 endpoint detected."
	if [[ "${AWS_ACCESS_KEY_ID}" == "" ]]; then
		echo "AWS_ACCESS_KEY_ID needs to be set"
		cmd_output_status 2 "AWS_ACCESS_KEY_ID_missing"
	fi
	if [[ "${AWS_SECRET_ACCESS_KEY}" == "" ]]; then
		echo "AWS_SECRET_ACCESS_KEY needs to be set"
		cmd_output_status 2 "AWS_SECRET_ACCESS_KEY_missing"
	fi
}

config_webdav() {
	echo "WebDAV endpoint detected."

	if [[ ${RCLONE_WEBDAV_USER} == "" ]]; then
		echo "RCLONE_WEBDAV_USER needs to be set"
		cmd_output_status 2 "RCLONE_WEBDAV_USER_missing"
	fi
	if [[ ${RCLONE_WEBDAV_PASS} == "" ]]; then
		echo "RCLONE_WEBDAV_PASS needs to be set"
		cmd_output_status 2 "RCLONE_WEBDAV_PASS_missing"
	fi
	if [[ ${RCLONE_WEBDAV_URL} == "" ]]; then
		echo "RCLONE_WEBDAV_URL needs to be set"
		cmd_output_status 2 "RCLONE_WEBDAV_URL_missing"
	fi

	export RCLONE_WEBDAV_PASS=$(rclone obscure "${RCLONE_WEBDAV_PASS}")
	if [[ $? -ne 0 ]]; then
		cmd_output_status 5 "RCLONE_WEBDAV_PASS_obscuring_failed"
	fi
}

cmd_restic() {
	set -o pipefail
	restic --json --verbose ${@} 2>&1 | tee -a ${RESTIC_LOG_DIR}/restic.log
	set +o pipefail
}

cmd_output_status() {
	status=$1 # integer, treat 0 as success
	message=$2 # string
	echo "{\"status\":${status}, \"message\":\"${message}\"}" > ${RESTIC_LOG_DIR}/last_status

	if [[ ${status} -gt 0 ]]; then
		exit ${status}
	fi
}

create_log_directory() {
	mkdir -p ${RESTIC_LOG_DIR}
}

do_backup() {
	cmd_output_status 0 "restic_starting"
	echo "Making sure restic repository is initialized.."
	repo_init=$(cmd_restic init)
	if [[ $? -ne 0 ]]; then
		if [[ $(echo $repo_init | grep -c "repository master key and config already initialized") -ne 1 ]]; then
			echo "$repo_init"
			cmd_output_status 3 "repo_init_failed"
		fi
	fi
	echo "Starting backup of '${RESTIC_BACKUP_SOURCE}'"
	cmd_restic backup "${RESTIC_BACKUP_SOURCE}"
	if [ $? -ne 0 ]; then
		cmd_output_status 4 "backup_failed"
	fi
	echo "Backup complete."
	cmd_output_status 0 "backup_complete"
}

create_log_directory
config_print
config_load
do_backup
