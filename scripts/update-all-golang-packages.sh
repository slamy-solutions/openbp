SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR/../

# Update all go packages
go work sync
go list -f '{{.Dir}}' -m | xargs -t -L1 bash -c 'cd "$0" && go mod tidy'
go list -f '{{.Dir}}' -m | xargs -t -L1 bash -c 'cd "$0" && go get -u'
go work sync
go list -f '{{.Dir}}' -m | xargs -t -L1 bash -c 'cd "$0" && go mod tidy'
go work sync