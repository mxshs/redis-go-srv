addr="localhost"
port="6380"
arr=""

while getopts a:p: o; do
    case $o in
        p) port="$OPTARG";;
        a) addr="$OPTARG";;
    esac
done
shift "$((OPTIND - 1))"

for word in "$@"; do
    word="\$${#word}\r\n${word}\r\n"
    arr+=$word
done

arr="*${#@}\r\n${arr}"

printf $arr | nc $addr $port
