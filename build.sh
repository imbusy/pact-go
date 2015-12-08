echo "Getting dependencies"
go get -d .

if [ $? -ne 0 ]; then
    echo "Could not save dependencies"
    exit 1
fi
echo "Vetting"
go vet -v ./...

echo "Running Tests"
go test -v -cover ./...

if [ $? -ne 0 ]; then
    echo "Tests failed!"
    exit 1
fi

echo "Building Application"
go build


if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build succeded!"
