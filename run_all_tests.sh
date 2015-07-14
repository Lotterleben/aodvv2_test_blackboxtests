for directory in $(ls -d test_*/);
do
    go run $directory/main.go;
done