function saveDependencies {
	echo "Getting dependencies"
	go get -d .

	echo "Saving dependencies"
	$GOPATH/bin/godep save ./...
}

function installGodep {
	if [ ! -f $GOAPTH/bin/godep ];
	then
		go get github.com/tools/godep
	fi
}

function runTests {
	echo "Running Tests"
	go test -v -cover ./...

	
}

function build {
	echo "Building Application"

	go build

}


installGodep

if [ $? -ne 0 ]; then
    echo "Could not install godep"
    exit 1
fi

saveDependencies

if [ $? -ne 0 ]; then
    echo "Could not save dependencies"
    exit 1
fi

runTests

if [ $? -ne 0 ]; then
    echo "Tests failed!"
    exit 1
fi

build


if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi


echo "Build succeded!"